package walk

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func benchmarkWalker() *Walker {
	return New(Config{UseTags: true, UseFieldDetection: true, UseValueDetection: true, Detector: detection.NewDetector(detection.DefaultRules(), true), Preservation: redactcore.DefaultConfig()}, nil)
}

func BenchmarkWalkerCopy(b *testing.B) {
	w := benchmarkWalker(); input := typedCustomer{Email: "greg@example.com", Secret: "secret"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ { _, _ = w.Copy(input) }
}
