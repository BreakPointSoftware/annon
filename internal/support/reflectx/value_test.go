package reflectx

import (
	"reflect"
	"testing"
)

func TestInterface(t *testing.T) {
	if got := Interface(reflect.Value{}); got != nil {
		t.Fatalf("expected nil for invalid reflect value, got %#v", got)
	}

	if got := Interface(reflect.ValueOf("value")); got != "value" {
		t.Fatalf("expected reflected string value, got %#v", got)
	}
}
