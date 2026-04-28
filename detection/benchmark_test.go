package detection

import "testing"

func BenchmarkCompiledDetectorDetectField(b *testing.B) {
	detector := NewCompiledDetector(DefaultRules(), PatternValueDetector{}, false)
	fields := []string{
		"email",
		"phoneNumber",
		"customerName",
		"vehicleRegistration",
		"username",
		"plainField",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = detector.DetectField(fields[i%len(fields)])
	}
}

func BenchmarkCompiledDetectorDetectValue(b *testing.B) {
	detector := NewCompiledDetector(DefaultRules(), PatternValueDetector{}, true)
	values := []string{
		"greg@example.com",
		"07700 900123",
		"TN9 1XA",
		"AB12 CDE",
		"plain text",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = detector.Detect("note", values[i%len(values)])
	}
}
