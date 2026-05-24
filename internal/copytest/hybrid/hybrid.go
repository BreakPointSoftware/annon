package hybrid

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/BreakPointSoftware/annon/internal/detection"
)

type FlagReason string
type FlagAction string

const (
	SensitiveFieldName          FlagReason = "SensitiveFieldName"
	ExportedReferenceDeepCopied FlagReason = "ExportedReferenceDeepCopied"
	UnexportedReferenceShared   FlagReason = "UnexportedReferenceShared"
	RuntimeStateZeroed          FlagReason = "RuntimeStateZeroed"
	UnsupportedKind             FlagReason = "UnsupportedKind"
	RecursiveReferenceReused    FlagReason = "RecursiveReferenceReused"
)

const (
	ActionDeepCopied  FlagAction = "deep-copied"
	ActionSkipped     FlagAction = "skipped"
	ActionZeroed      FlagAction = "zeroed"
	ActionShared      FlagAction = "shared"
	ActionUnsupported FlagAction = "unsupported"
	ActionReused      FlagAction = "reused"
)

type FieldFlag struct {
	Path   string
	Type   reflect.Type
	Kind   reflect.Kind
	Reason FlagReason
	Action FlagAction
}

type HybridCopyResult[T any] struct {
	Copy  T
	Flags []FieldFlag
}

type visitKey struct {
	typ reflect.Type
	ptr uintptr
}

type hybridWalker struct {
	visited  map[visitKey]reflect.Value
	flags    []FieldFlag
	detector *detection.Detector
}

func Copy[T any](input T) (HybridCopyResult[T], error) {
	copyWalker := &hybridWalker{visited: map[visitKey]reflect.Value{}, detector: detection.NewDetector(detection.DefaultRules(), false)}
	inputValue := reflect.ValueOf(input)
	copiedValue, err := copyWalker.copyValue(inputValue, "")
	if err != nil {
		return HybridCopyResult[T]{}, err
	}
	if !copiedValue.IsValid() {
		var zero T
		return HybridCopyResult[T]{Copy: zero, Flags: copyWalker.flags}, nil
	}
	return HybridCopyResult[T]{Copy: copiedValue.Interface().(T), Flags: copyWalker.flags}, nil
}

func (w *hybridWalker) copyValue(inputValue reflect.Value, path string) (reflect.Value, error) {
	if !inputValue.IsValid() {
		return inputValue, nil
	}

	switch inputValue.Kind() {
	case reflect.Interface:
		return w.copyInterface(inputValue, path)
	case reflect.Pointer:
		return w.copyPointer(inputValue, path)
	case reflect.Struct:
		return w.copyStruct(inputValue, path)
	case reflect.Map:
		return w.copyMap(inputValue, path)
	case reflect.Slice:
		return w.copySlice(inputValue, path)
	case reflect.Array:
		return w.copyArray(inputValue, path)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		w.flags = append(w.flags, FieldFlag{Path: path, Type: inputValue.Type(), Kind: inputValue.Kind(), Reason: UnsupportedKind, Action: ActionZeroed})
		return reflect.Zero(inputValue.Type()), nil
	default:
		return inputValue, nil
	}
}

func (w *hybridWalker) copyInterface(inputValue reflect.Value, path string) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}
	innerCopy, err := w.copyValue(inputValue.Elem(), path)
	if err != nil {
		return reflect.Value{}, err
	}
	outputValue := reflect.New(inputValue.Type()).Elem()
	outputValue.Set(innerCopy)
	return outputValue, nil
}

func (w *hybridWalker) copyPointer(inputValue reflect.Value, path string) (reflect.Value, error) {
	if inputValue.IsNil() {
		return reflect.Zero(inputValue.Type()), nil
	}
	key := visitKey{typ: inputValue.Type(), ptr: inputValue.Pointer()}
	if reusedValue, ok := w.visited[key]; ok {
		w.flags = append(w.flags, FieldFlag{Path: path, Type: inputValue.Type(), Kind: inputValue.Kind(), Reason: RecursiveReferenceReused, Action: ActionReused})
		return reusedValue, nil
	}
	outputPointer := reflect.New(inputValue.Type().Elem())
	w.visited[key] = outputPointer
	childCopy, err := w.copyValue(inputValue.Elem(), path)
	if err != nil {
		return reflect.Value{}, err
	}
	outputPointer.Elem().Set(childCopy)
	return outputPointer, nil
}

func (w *hybridWalker) copyStruct(structValue reflect.Value, path string) (reflect.Value, error) {
	outputStruct := reflect.New(structValue.Type()).Elem()
	outputStruct.Set(structValue)

	for fieldIndex := 0; fieldIndex < structValue.NumField(); fieldIndex++ {
		structField := structValue.Type().Field(fieldIndex)
		fieldValue := structValue.Field(fieldIndex)
		fieldPath := joinPath(path, structField.Name)

		if match := w.detector.DetectField(structField.Name); match.Found() {
			w.flags = append(w.flags, FieldFlag{Path: fieldPath, Type: structField.Type, Kind: fieldValue.Kind(), Reason: SensitiveFieldName, Action: ActionSkipped})
		}

		if shouldZeroRuntimeState(structField.Type, fieldValue.Kind()) {
			if outputStruct.Field(fieldIndex).CanSet() {
				outputStruct.Field(fieldIndex).Set(reflect.Zero(structField.Type))
				w.flags = append(w.flags, FieldFlag{Path: fieldPath, Type: structField.Type, Kind: fieldValue.Kind(), Reason: RuntimeStateZeroed, Action: ActionZeroed})
			} else {
				w.flags = append(w.flags, FieldFlag{Path: fieldPath, Type: structField.Type, Kind: fieldValue.Kind(), Reason: RuntimeStateZeroed, Action: ActionShared})
			}
			continue
		}

		if structField.PkgPath != "" {
			if isReferenceKind(fieldValue.Kind()) {
				w.flags = append(w.flags, FieldFlag{Path: fieldPath, Type: structField.Type, Kind: fieldValue.Kind(), Reason: UnexportedReferenceShared, Action: ActionShared})
			}
			continue
		}

		if !isReferenceKind(fieldValue.Kind()) {
			continue
		}

		copiedValue, err := w.copyValue(fieldValue, fieldPath)
		if err != nil {
			return reflect.Value{}, err
		}
		outputStruct.Field(fieldIndex).Set(copiedValue)
		w.flags = append(w.flags, FieldFlag{Path: fieldPath, Type: structField.Type, Kind: fieldValue.Kind(), Reason: ExportedReferenceDeepCopied, Action: ActionDeepCopied})
	}

	return outputStruct, nil
}

func (w *hybridWalker) copyMap(mapValue reflect.Value, path string) (reflect.Value, error) {
	if mapValue.IsNil() {
		return reflect.Zero(mapValue.Type()), nil
	}
	key := visitKey{typ: mapValue.Type(), ptr: mapValue.Pointer()}
	if reusedValue, ok := w.visited[key]; ok {
		w.flags = append(w.flags, FieldFlag{Path: path, Type: mapValue.Type(), Kind: mapValue.Kind(), Reason: RecursiveReferenceReused, Action: ActionReused})
		return reusedValue, nil
	}
	outputMap := reflect.MakeMapWithSize(mapValue.Type(), mapValue.Len())
	w.visited[key] = outputMap
	iter := mapValue.MapRange()
	for iter.Next() {
		copiedValue, err := w.copyValue(iter.Value(), fmt.Sprintf("%s[%v]", path, iter.Key().Interface()))
		if err != nil {
			return reflect.Value{}, err
		}
		outputMap.SetMapIndex(iter.Key(), copiedValue)
	}
	return outputMap, nil
}

func (w *hybridWalker) copySlice(sliceValue reflect.Value, path string) (reflect.Value, error) {
	if sliceValue.IsNil() {
		return reflect.Zero(sliceValue.Type()), nil
	}
	key := visitKey{typ: sliceValue.Type(), ptr: sliceValue.Pointer()}
	if sliceValue.Pointer() != 0 {
		if reusedValue, ok := w.visited[key]; ok {
			w.flags = append(w.flags, FieldFlag{Path: path, Type: sliceValue.Type(), Kind: sliceValue.Kind(), Reason: RecursiveReferenceReused, Action: ActionReused})
			return reusedValue, nil
		}
	}
	outputSlice := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), sliceValue.Len())
	if sliceValue.Pointer() != 0 {
		w.visited[key] = outputSlice
	}
	for index := 0; index < sliceValue.Len(); index++ {
		copiedValue, err := w.copyValue(sliceValue.Index(index), fmt.Sprintf("%s[%d]", path, index))
		if err != nil {
			return reflect.Value{}, err
		}
		outputSlice.Index(index).Set(copiedValue)
	}
	return outputSlice, nil
}

func (w *hybridWalker) copyArray(arrayValue reflect.Value, path string) (reflect.Value, error) {
	outputArray := reflect.New(arrayValue.Type()).Elem()
	for index := 0; index < arrayValue.Len(); index++ {
		copiedValue, err := w.copyValue(arrayValue.Index(index), fmt.Sprintf("%s[%d]", path, index))
		if err != nil {
			return reflect.Value{}, err
		}
		outputArray.Index(index).Set(copiedValue)
	}
	return outputArray, nil
}

func isReferenceKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface:
		return true
	default:
		return false
	}
}

func joinPath(parent, field string) string {
	if parent == "" {
		return field
	}
	return parent + "." + field
}

var (
	mutexType       = reflect.TypeOf(sync.Mutex{})
	rwMutexType     = reflect.TypeOf(sync.RWMutex{})
	onceType        = reflect.TypeOf(sync.Once{})
	waitGroupType   = reflect.TypeOf(sync.WaitGroup{})
	syncMapType     = reflect.TypeOf(sync.Map{})
	atomicValueType = reflect.TypeOf(atomic.Value{})
	contextType     = reflect.TypeOf((*context.Context)(nil)).Elem()
	filePtrType     = reflect.TypeOf((*os.File)(nil))
	connType        = reflect.TypeOf((*net.Conn)(nil)).Elem()
	dbPtrType       = reflect.TypeOf((*sql.DB)(nil))
	txPtrType       = reflect.TypeOf((*sql.Tx)(nil))
	clientPtrType   = reflect.TypeOf((*http.Client)(nil))
	condPtrType     = reflect.TypeOf((*sync.Cond)(nil))
)

func shouldZeroRuntimeState(typ reflect.Type, kind reflect.Kind) bool {
	if kind == reflect.Chan || kind == reflect.Func {
		return true
	}
	if typ == mutexType || typ == rwMutexType || typ == onceType || typ == waitGroupType || typ == syncMapType || typ == atomicValueType {
		return true
	}
	if typ == condPtrType || typ == filePtrType || typ == dbPtrType || typ == txPtrType || typ == clientPtrType {
		return true
	}
	if typ.Implements(contextType) {
		return true
	}
	if typ.Kind() == reflect.Interface && typ.Implements(connType) {
		return true
	}
	return false
}
