package main

import (
	"strings"
	"testing"
)

func TestBuildSteps(t *testing.T) {
	steps := buildSteps()
	if len(steps) == 0 {
		t.Fatal("expected presentation steps")
	}

	if steps[0].Title != "Step 1: Structured redaction" {
		t.Fatalf("unexpected first step title: %q", steps[0].Title)
	}
}

func TestRenderStep(t *testing.T) {
	step := buildSteps()[0]
	rendered := renderStep(1, 4, step)

	checks := []string{
		"Step 1 of 4",
		"Original object",
		"redact.Data output",
		"greg@example.com",
		"g***@example.com",
	}

	for _, check := range checks {
		if !strings.Contains(rendered, check) {
			t.Fatalf("expected rendered step to contain %q\n%s", check, rendered)
		}
	}
}
