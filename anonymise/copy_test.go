package anonymise

import "testing"

type customer struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Secret string `json:"secret" anonymise:"remove"`
	Note   string `json:"note" anonymise:"false"`
}

func TestCopy(t *testing.T) {
	input := customer{ID: "123", Email: "greg@example.com", Secret: "secret", Note: "greg@example.com"}
	a, err := New()
	if err != nil {
		t.Fatal(err)
	}
	resultAny, err := a.Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	result := resultAny.(customer)
	if result.Email != "g***@example.com" || result.Secret != "" || result.Note != "greg@example.com" {
		t.Fatalf("unexpected copy result: %+v", result)
	}
	if input.Secret != "secret" {
		t.Fatalf("input mutated: %+v", input)
	}
}

func TestPackageCopyPreservesType(t *testing.T) {
	input := customer{ID: "123", Email: "greg@example.com"}
	result, err := Copy(input)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := any(result).(customer); !ok {
		t.Fatalf("unexpected result type: %T", result)
	}
}
