package walker

import (
	"testing"

	"github.com/your-org/annon/detection"
	"github.com/your-org/annon/preservation"
	strategypkg "github.com/your-org/annon/strategy"
)

type typedCustomer struct {
	ID     string
	Email  string `json:"email"`
	Secret string `anonymise:"remove"`
	Note   string `anonymise:"false"`
}

type typedNested struct {
	Contact *typedCustomer
	Items   []typedCustomer
	Labels  map[string]string
	Codes   [2]string
}

func testWalkerConfig() Config {
	strategies := map[string]strategypkg.Strategy{}
	for _, s := range strategypkg.DefaultStrategies() {
		strategies[s.Name()] = s
	}
	return Config{
		UseTags:           true,
		UseFieldDetection: true,
		UseValueDetection: false,
		Detector:          detection.NewCompiledDetector(detection.DefaultRules(), detection.PatternValueDetector{}, false),
		Strategies:        strategies,
		Preservation:      preservation.Default(),
	}
}

func TestCopyStruct(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	input := typedCustomer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "greg@example.com"}
	resultAny, err := w.Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	result := resultAny.(typedCustomer)
	if result.Email != "g***@example.com" {
		t.Fatalf("unexpected masked email: %q", result.Email)
	}
	if result.Secret != "" {
		t.Fatalf("expected removed secret to be zero, got %q", result.Secret)
	}
	if result.Note != "greg@example.com" {
		t.Fatalf("expected false tag to preserve note, got %q", result.Note)
	}
	if input.Secret != "secret" {
		t.Fatalf("input mutated: %+v", input)
	}
}

func TestFieldDetectionCanBeDisabled(t *testing.T) {
	cfg := testWalkerConfig()
	cfg.UseFieldDetection = false
	w := New(cfg, nil)
	input := typedCustomer{Email: "greg@example.com"}
	resultAny, err := w.Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	result := resultAny.(typedCustomer)
	if result.Email != "greg@example.com" {
		t.Fatalf("expected email to remain unchanged, got %q", result.Email)
	}
}

func TestValueDetectionCanBeEnabledIndependently(t *testing.T) {
	cfg := testWalkerConfig()
	cfg.UseFieldDetection = false
	cfg.UseValueDetection = true
	cfg.Detector = detection.NewCompiledDetector(nil, detection.PatternValueDetector{}, true)
	w := New(cfg, nil)
	input := map[string]string{"note": "greg@example.com"}
	resultAny, err := w.Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	result := resultAny.(map[string]string)
	if result["note"] != "g***@example.com" {
		t.Fatalf("expected value detection to mask email, got %#v", result)
	}
}

func TestInvalidTagReturnsError(t *testing.T) {
	type invalid struct {
		Email string `anonymise:"notARealStrategy"`
	}
	w := New(testWalkerConfig(), nil)
	_, err := w.Copy(invalid{Email: "greg@example.com"})
	if err == nil {
		t.Fatal("expected invalid tag error")
	}
}

func TestCopyHandlesNilAndNestedValues(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	input := typedNested{
		Contact: &typedCustomer{ID: "1", Email: "greg@example.com"},
		Items:   []typedCustomer{{Email: "one@example.com"}, {Email: "two@example.com"}},
		Labels:  map[string]string{"email": "greg@example.com", "plain": "ok"},
		Codes:   [2]string{"AB12 CDE", "TN9 1XA"},
	}
	resultAny, err := w.Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	result := resultAny.(typedNested)
	if result.Contact == nil || result.Contact.Email != "g***@example.com" {
		t.Fatalf("unexpected nested pointer result: %+v", result.Contact)
	}
	if result.Items[0].Email != "o**@example.com" || result.Items[1].Email != "t**@example.com" {
		t.Fatalf("unexpected slice result: %+v", result.Items)
	}
	if result.Labels["email"] != "g***@example.com" || result.Labels["plain"] != "ok" {
		t.Fatalf("unexpected map result: %+v", result.Labels)
	}
	if result.Codes[0] != "AB12 CDE" || result.Codes[1] != "TN9 1XA" {
		t.Fatalf("unexpected array result: %+v", result.Codes)
	}
	if input.Contact.Email != "greg@example.com" || input.Items[0].Email != "one@example.com" {
		t.Fatalf("input mutated: %+v", input)
	}
}

func TestCopyHandlesNilPointersMapsAndSlices(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	input := typedNested{}
	resultAny, err := w.Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	result := resultAny.(typedNested)
	if result.Contact != nil || result.Items != nil || result.Labels != nil {
		t.Fatalf("expected nil values to remain nil: %+v", result)
	}
}
