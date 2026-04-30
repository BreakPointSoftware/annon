package testdata

import (
	"context"
	"sync"
	"net/http"
)

func DemoValueOnly() ValueOnly {
	v := ValueOnly{Name: "value-only", Age: 42, Enabled: true}
	v.Hash[0] = 1
	v.Hash[31] = 8
	return v
}

func DemoExportedRefs() ExportedRefs {
	return ExportedRefs{
		Child: &Child{Name: "child"},
		Names: []string{"one", "two"},
		Meta:  map[string]int{"a": 1, "b": 2},
	}
}

func DemoCycle() *Cat {
	cat := &Cat{Name: "milo"}
	cat.Sibling = cat
	return cat
}

func DemoSiblingCycle() *Cat {
	left := &Cat{Name: "left"}
	right := &Cat{Name: "right"}
	left.Sibling = right
	right.Sibling = left
	return left
}

func DemoSharedChild() SharedChild {
	child := &Child{Name: "shared"}
	return SharedChild{Left: child, Right: child}
}

func DemoRuntimeState() RuntimeState {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	state := RuntimeState{
		Once:      sync.Once{},
		WaitGroup: waitGroup,
		Cond:      sync.NewCond(&sync.Mutex{}),
		Chan:      make(chan string, 1),
		Func:      func() string { return "runtime" },
		Ctx:       context.Background(),
		Client:    &http.Client{},
		Name:      "runtime-state",
	}
	state.Map.Store("key", "value")
	state.Atomic.Store("atomic")
	return state
}

func DemoInterfaceFields() InterfaceFields {
	shared := &Child{Name: "interface-child"}
	return InterfaceFields{
		Pointer: shared,
		Slice:   []*Child{shared},
		Map:     map[string]*Child{"primary": shared},
	}
}

func DemoNestedCollections() NestedCollections {
	shared := &Child{Name: "nested"}
	return NestedCollections{
		Groups: map[string][]*Child{
			"one": []*Child{shared},
			"two": []*Child{shared},
		},
		Items: []map[string]*Child{{"first": shared}, {"second": shared}},
	}
}

func DemoNilRefs() NilRefs {
	return NilRefs{}
}

func DemoLargeValue() LargeValue {
	v := LargeValue{Name: "large", Count: 99}
	v.A[0], v.B[1], v.C[2], v.D[3] = 1, 2, 3, 4
	v.E[4], v.F[5], v.G[6], v.H[7] = 5, 6, 7, 8
	return v
}

func DemoManyPointers() ManyPointers {
	return ManyPointers{
		One:   &Child{Name: "one"},
		Two:   &Child{Name: "two"},
		Three: &Child{Name: "three"},
		Four:  &Child{Name: "four"},
		Five:  &Child{Name: "five"},
		Six:   &Child{Name: "six"},
		Seven: &Child{Name: "seven"},
		Eight: &Child{Name: "eight"},
	}
}

func DemoTree() *Tree {
	return &Tree{
		Name: "root",
		Children: []*Tree{
			{Name: "left", Children: []*Tree{{Name: "left-leaf"}}},
			{Name: "right", Children: []*Tree{{Name: "right-leaf"}}},
		},
	}
}

func DemoDomainObject() DomainObject {
	return DomainObject{
		Customer: DemoExportedRefs(),
		Shared:   DemoSharedChild(),
		Nested:   DemoNestedCollections(),
		Private:  NewWithPrivateRef(),
		Runtime:  DemoRuntimeState(),
	}
}
