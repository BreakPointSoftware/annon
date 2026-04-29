package copy

import (
	"reflect"
	"testing"
)

type cacheExample struct {
	Email string `json:"email" yaml:"email_address" anonymise:"email"`
	Skip  string `json:"-"`
	plain string
}

func TestTypeCacheReusesCompiledMetadata(t *testing.T) {
	cache := NewTypeCache()
	typeOf := reflect.TypeOf(cacheExample{})
	first := cache.StructFields(typeOf)
	second := cache.StructFields(typeOf)
	if &first[0] != &second[0] {
		t.Fatal("expected reuse")
	}
}
