package anonymise

import (
	"errors"
	"reflect"
	"testing"

	"github.com/BreakPointSoftware/annon/detection"
)

func TestNewInitialisesWalker(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if a.walker == nil || a.cache == nil || a.config.Strategies["email"] == nil {
		t.Fatal("anonymiser was not initialised correctly")
	}
}

func TestNewReturnsOptionError(t *testing.T) {
	want := errors.New("boom")
	_, err := New(func(*Config) error { return want })
	if !errors.Is(err, want) {
		t.Fatalf("expected option error, got %v", err)
	}
}

func TestAdditionalFieldRulesAreApplied(t *testing.T) {
	type custom struct {
		Alias string `json:"customerAlias"`
	}
	a, err := New(WithFieldRules(detection.StrongRule(detection.Redact, "customerAlias")))
	if err != nil {
		t.Fatal(err)
	}
	result, err := a.JSON(custom{Alias: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != `{"customerAlias":"[REDACTED]"}` {
		t.Fatalf("unexpected custom field-rule output: %s", result)
	}
}

func TestCustomDetectorAndFieldRulesConflict(t *testing.T) {
	_, err := New(
		WithDetector(detection.NewCompiledDetector(nil, detection.PatternValueDetector{}, false)),
		WithFieldRules(detection.StrongRule(detection.Redact, "customerAlias")),
	)
	if err == nil {
		t.Fatal("expected conflicting detector/field-rule configuration to fail")
	}
}

func TestReusableAnonymiserReusesTypeCache(t *testing.T) {
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	type sample struct {
		Email string `json:"email"`
	}
	if _, err := a.JSON(sample{Email: "greg@example.com"}); err != nil {
		t.Fatal(err)
	}
	first := a.cache.StructFields(reflect.TypeOf(sample{}))
	if _, err := a.JSON(sample{Email: "one@example.com"}); err != nil {
		t.Fatal(err)
	}
	second := a.cache.StructFields(reflect.TypeOf(sample{}))
	if &first[0] != &second[0] {
		t.Fatal("expected reusable anonymiser to retain cached field metadata")
	}
}
