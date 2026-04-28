package walker

import (
	"reflect"
	"testing"
)

type cacheExample struct {
	Email string `json:"email" yaml:"email_address" anonymise:"email"`
	Skip  string `json:"-"`
	plain string
}

func TestTypeCache(t *testing.T) {
	cache := NewTypeCache()
	fields := cache.StructFields(reflect.TypeOf(cacheExample{}))
	if len(fields) != 2 {
		t.Fatalf("unexpected field count: %d", len(fields))
	}
	if fields[0].JSONName != "email" || fields[0].YAMLName != "email_address" || fields[0].AnonymiseTag != "email" {
		t.Fatalf("unexpected field metadata: %+v", fields[0])
	}
	if fields[1].IgnoreJSON != true {
		t.Fatalf("expected json ignore: %+v", fields[1])
	}
}

func TestTypeCacheReusesCompiledMetadata(t *testing.T) {
	cache := NewTypeCache()
	typeOf := reflect.TypeOf(cacheExample{})
	first := cache.StructFields(typeOf)
	second := cache.StructFields(typeOf)
	if len(first) == 0 || len(second) == 0 {
		t.Fatal("expected cached fields")
	}
	if &first[0] != &second[0] {
		t.Fatal("expected repeated lookups to reuse cached metadata slice")
	}
}
