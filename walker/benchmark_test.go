package walker

import (
	"testing"

	"github.com/your-org/annon/detection"
	"github.com/your-org/annon/preservation"
	strategypkg "github.com/your-org/annon/strategy"
)

type benchmarkCustomer struct {
	ID      string            `json:"id"`
	Email   string            `json:"email"`
	Phone   string            `json:"phoneNumber"`
	Name    string            `json:"customerName"`
	Vehicle benchmarkVehicle  `json:"vehicle"`
	Labels  map[string]string `json:"labels"`
	Items   []benchmarkItem   `json:"items"`
}

type benchmarkVehicle struct {
	Registration string `json:"reg"`
	Postcode     string `json:"postcode"`
}

type benchmarkItem struct {
	Email string `json:"email"`
	Note  string `json:"note"`
}

func benchmarkWalker() *Walker {
	strategies := map[string]strategypkg.Strategy{}
	for _, s := range strategypkg.DefaultStrategies() {
		strategies[s.Name()] = s
	}
	return New(Config{
		UseTags:           true,
		UseFieldDetection: true,
		UseValueDetection: true,
		Detector:          detection.NewCompiledDetector(detection.DefaultRules(), detection.PatternValueDetector{}, true),
		Strategies:        strategies,
		Preservation:      preservation.Default(),
	}, nil)
}

func benchmarkInput() benchmarkCustomer {
	return benchmarkCustomer{
		ID:    "123",
		Email: "greg@example.com",
		Phone: "07700 900123",
		Name:  "Greg Bryant",
		Vehicle: benchmarkVehicle{
			Registration: "AB12 CDE",
			Postcode:     "TN9 1XA",
		},
		Labels: map[string]string{
			"email":    "greg@example.com",
			"reference": "ABC123",
		},
		Items: []benchmarkItem{{Email: "one@example.com", Note: "plain"}, {Email: "two@example.com", Note: "07700 900123"}},
	}
}

func BenchmarkWalkerCopy(b *testing.B) {
	w := benchmarkWalker()
	input := benchmarkInput()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := w.Copy(input); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWalkerBlobFromValue(b *testing.B) {
	w := benchmarkWalker()
	input := benchmarkInput()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := w.BlobFromValue(input, "json"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWalkerBlobFromNeutral(b *testing.B) {
	w := benchmarkWalker()
	input := map[string]any{
		"email":       "greg@example.com",
		"phoneNumber": "07700 900123",
		"vehicle": map[string]any{
			"reg":      "AB12 CDE",
			"postcode": "TN9 1XA",
		},
		"items": []any{
			map[string]any{"email": "one@example.com"},
			map[string]any{"email": "two@example.com", "note": "plain"},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := w.BlobFromNeutral(input); err != nil {
			b.Fatal(err)
		}
	}
}
