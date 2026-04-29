package walk

import "testing"

func TestOutputFromNeutral(t *testing.T) {
	w := New(testWalkerConfig(), nil)
	result, err := w.OutputFromNeutral(map[string]any{"email": "greg@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	outputMap := result.(map[string]any)
	if outputMap["email"] != "g***@example.com" {
		t.Fatalf("unexpected result: %#v", outputMap)
	}
}
