package anonymise

import "testing"

type benchmarkCustomer struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Phone   string `json:"phoneNumber"`
	Name    string `json:"customerName"`
	Vehicle struct {
		Registration string `json:"reg"`
		Postcode     string `json:"postcode"`
	} `json:"vehicle"`
}

func benchmarkAnonymiserInput() benchmarkCustomer {
	input := benchmarkCustomer{
		ID:    "123",
		Email: "greg@example.com",
		Phone: "07700 900123",
		Name:  "Greg Bryant",
	}
	input.Vehicle.Registration = "AB12 CDE"
	input.Vehicle.Postcode = "TN9 1XA"
	return input
}

func BenchmarkAnonymiserJSON(b *testing.B) {
	a, err := New(WithValueDetection(true))
	if err != nil {
		b.Fatal(err)
	}
	input := benchmarkAnonymiserInput()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := a.JSON(input); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnonymiserFromJSON(b *testing.B) {
	a, err := New(WithValueDetection(true))
	if err != nil {
		b.Fatal(err)
	}
	input := []byte(`{"email":"greg@example.com","phoneNumber":"07700 900123","vehicle":{"reg":"AB12 CDE","postcode":"TN9 1XA"}}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := a.FromJSON(input); err != nil {
			b.Fatal(err)
		}
	}
}
