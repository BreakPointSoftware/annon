package walker

import "testing"

func TestParseTag(t *testing.T) {
	tests := map[string]parsedTag{
		"":                     {empty: true},
		"false":                {skip: true},
		"true":                 {auto: true},
		"auto":                 {auto: true},
		"email":                {strategyName: "email"},
		"vehicleRegistration":  {strategyName: "vehicleRegistration"},
		"redact":               {strategyName: "redact"},
		"remove":               {remove: true, strategyName: "remove"},
	}
	for input, want := range tests {
		if got := parseTag(input); got != want {
			t.Fatalf("parseTag(%q) = %+v, want %+v", input, got, want)
		}
	}
}
