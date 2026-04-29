package detection

import "testing"

func BenchmarkDetectorDetectField(b *testing.B) {
	detector := NewDetector(DefaultRules(), false)
	fields := []string{"email", "phoneNumber", "customerName", "vehicleRegistration", "username", "plainField"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = detector.DetectField(fields[i%len(fields)])
	}
}
