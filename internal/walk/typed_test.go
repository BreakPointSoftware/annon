package walk

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

type typedCustomer struct {
	Email string `json:"email"`
	Secret string `anonymise:"remove"`
}

func testWalkerConfig() Config {
	return Config{UseTags: true, UseFieldDetection: true, UseValueDetection: false, Detector: detection.NewCompiledDetector(detection.DefaultRules(), detection.PatternValueDetector{}, false), Preservation: redactcore.DefaultConfig()}
}

func TestCopyStruct(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	resultAny, err := w.Copy(typedCustomer{Email: "greg@example.com", Secret: "secret"})
	if err != nil { t.Fatal(err) }
	result := resultAny.(typedCustomer)
	if result.Email != "g***@example.com" || result.Secret != "" { t.Fatalf("unexpected copy result: %+v", result) }
}
