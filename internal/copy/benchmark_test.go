package copy

import (
	"testing"

	"github.com/BreakPointSoftware/annon/internal/decision"
	"github.com/BreakPointSoftware/annon/internal/detection"
	"github.com/BreakPointSoftware/annon/internal/redactcore"
)

func benchmarkCopier() *Copier {
	config := decision.Config{UseTags: true, UseFieldDetection: true, UseValueDetection: true, Detector: detection.NewDetector(detection.DefaultRules(), true), Preservation: redactcore.DefaultConfig()}
	return New(config, decision.New(config), nil)
}

func BenchmarkCopierCopy(b *testing.B) {
	copier := benchmarkCopier()
	input := typedCustomer{Email: "greg@example.com", Secret: "secret"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = copier.Copy(input)
	}
}
