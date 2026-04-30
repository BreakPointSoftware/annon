package main

import (
	"strings"
	"testing"
)

func TestBuildOutput(t *testing.T) {
	output := buildOutput()

	checks := []string{
		"Structured redaction with redact.Data",
		"Original",
		"Redacted",
		"greg@example.com",
		"g***@example.com",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected output to contain %q\n%s", check, output)
		}
	}
}
