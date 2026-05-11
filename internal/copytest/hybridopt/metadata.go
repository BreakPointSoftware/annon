package hybridopt

import (
	"reflect"
	"sync"

	"github.com/BreakPointSoftware/annon/internal/detection"
)

type fieldPlan struct {
	index        int
	name         string
	typ          reflect.Type
	kind         reflect.Kind
	exported     bool
	sensitive    bool
	runtimeState bool
	reference    bool
}

type repairPlan struct {
	fields        []fieldPlan
	hasRepairWork bool
}

type planCache struct {
	plans sync.Map
}

var globalPlanCache planCache
var metadataDetector = detection.NewDetector(detection.DefaultRules(), false)

func (c *planCache) PlanFor(typ reflect.Type) *repairPlan {
	if cachedPlan, ok := c.plans.Load(typ); ok {
		return cachedPlan.(*repairPlan)
	}

	compiledPlan := compileRepairPlan(typ)
	c.plans.Store(typ, compiledPlan)
	return compiledPlan
}

func compileRepairPlan(typ reflect.Type) *repairPlan {
	compiledPlan := &repairPlan{}

	for fieldIndex := 0; fieldIndex < typ.NumField(); fieldIndex++ {
		structField := typ.Field(fieldIndex)
		kind := structField.Type.Kind()
		reference := isReferenceKind(kind)
		runtimeState := shouldZeroRuntimeState(structField.Type, kind)
		sensitive := metadataDetector.DetectField(structField.Name).Found()
		exported := structField.PkgPath == ""

		compiledPlan.fields = append(compiledPlan.fields, fieldPlan{
			index:        fieldIndex,
			name:         structField.Name,
			typ:          structField.Type,
			kind:         kind,
			exported:     exported,
			sensitive:    sensitive,
			runtimeState: runtimeState,
			reference:    reference,
		})

		if runtimeState || (exported && reference) {
			compiledPlan.hasRepairWork = true
		}

	}

	return compiledPlan
}
