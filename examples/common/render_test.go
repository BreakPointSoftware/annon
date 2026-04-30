package common

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPrettyJSON(t *testing.T) {
	rendered := PrettyJSON(DemoCustomer())

	var decoded any
	if err := json.Unmarshal([]byte(rendered), &decoded); err != nil {
		t.Fatalf("PrettyJSON returned invalid json: %v\n%s", err, rendered)
	}

	if !strings.Contains(rendered, "\n  ") {
		t.Fatalf("expected indented output, got %q", rendered)
	}
}

func TestSection(t *testing.T) {
	rendered := Section("Demo")
	if !strings.Contains(rendered, "Demo") || !strings.Contains(rendered, strings.Repeat("=", 60)) {
		t.Fatalf("unexpected section render: %q", rendered)
	}
}
