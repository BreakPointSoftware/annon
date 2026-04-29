package tags

import "testing"

func TestParseSerialiseTag(t *testing.T) {
	testCases := []struct {
		name      string
		tag       string
		expected  string
		ignoreTag bool
	}{
		{name: "empty tag", tag: "", expected: "", ignoreTag: false},
		{name: "ignore tag", tag: "-", expected: "", ignoreTag: true},
		{name: "name only", tag: "email", expected: "email", ignoreTag: false},
		{name: "name with options", tag: "email,omitempty", expected: "email", ignoreTag: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			name, ignoreTag := ParseSerialiseTag(testCase.tag)
			if name != testCase.expected || ignoreTag != testCase.ignoreTag {
				t.Fatalf("ParseSerialiseTag(%q) = (%q, %v)", testCase.tag, name, ignoreTag)
			}
		})
	}
}
