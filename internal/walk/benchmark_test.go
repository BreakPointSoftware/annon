package walk

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func benchmarkWalker() *Walker {
	walkerConfig := decision.Config{UseTags: true, UseFieldDetection: true, UseValueDetection: true, Detector: detection.NewDetector(detection.DefaultRules(), true), Preservation: redactcore.DefaultConfig()}
	return New(walkerConfig, decision.New(walkerConfig), nil)
}

func BenchmarkWalkerCopy(b *testing.B) {
	walker := benchmarkWalker()
	input := typedCustomer{Email: "greg@example.com", Secret: "secret"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = walker.Copy(input)
	}
}
