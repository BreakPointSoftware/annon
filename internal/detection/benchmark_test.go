package detection

import "testing"

func BenchmarkCompiledDetectorDetectField(b *testing.B) {
	detector := NewCompiledDetector(DefaultRules(), PatternValueDetector{}, false)
	fields := []string{"email", "phoneNumber", "customerName", "vehicleRegistration", "username", "plainField"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = detector.DetectField(fields[i%len(fields)])
	}
}
