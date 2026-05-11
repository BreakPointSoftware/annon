package hybridopt

import (
	"reflect"
	"testing"

	"github.com/BreakPointSoftware/annon/internal/copytest/testdata"
)

func TestCompileRepairPlanForValueOnlyStruct(t *testing.T) {
	plan := compileRepairPlan(reflect.TypeOf(testdata.DemoValueOnly()))
	if plan.hasRepairWork {
		t.Fatalf("expected no repair work for value-only struct")
	}
	if !plan.hasFieldFlags {
		t.Fatalf("expected field flags due to sensitive field-name detection")
	}
}

func TestCompileRepairPlanForExportedRefs(t *testing.T) {
	plan := compileRepairPlan(reflect.TypeOf(testdata.DemoExportedRefs()))
	if !plan.hasRepairWork {
		t.Fatalf("expected repair work for exported references")
	}
	if len(plan.fields) == 0 {
		t.Fatal("expected planned fields")
	}
}

func TestCompileRepairPlanForRuntimeState(t *testing.T) {
	plan := compileRepairPlan(reflect.TypeOf(testdata.DemoRuntimeState()))
	if !plan.hasRepairWork {
		t.Fatalf("expected runtime state to require repair")
	}
	var foundRuntimeState bool
	for _, field := range plan.fields {
		if field.runtimeState {
			foundRuntimeState = true
			break
		}
	}
	if !foundRuntimeState {
		t.Fatal("expected at least one runtime-state field in plan")
	}
}
