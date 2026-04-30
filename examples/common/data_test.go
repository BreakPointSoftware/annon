package common

import "testing"

func TestDemoCustomer(t *testing.T) {
	customer := DemoCustomer()

	if customer.Email == "" || customer.Phone == "" || customer.Postcode == "" {
		t.Fatalf("expected populated demo customer, got %+v", customer)
	}

	if customer.Notes == "" || customer.Secret == "" {
		t.Fatalf("expected note and secret demo fields, got %+v", customer)
	}

	if customer.Vehicle.Registration == "" || customer.Contact.Email == "" {
		t.Fatalf("expected nested demo values, got %+v", customer)
	}
}

func TestMalformedJSON(t *testing.T) {
	if got := string(MalformedJSON()); got != `{"email":` {
		t.Fatalf("unexpected malformed json fixture: %q", got)
	}
}
