package main

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/examples/common"
	"github.com/BreakPointSoftware/annon/redact"
)

type Step struct {
	Title       string
	Description string
	InputTitle  string
	InputText   string
	OutputTitle string
	OutputText  string
}

func buildSteps() []Step {
	originalCustomer := common.DemoCustomer()
	redactedCustomer := redact.Data(originalCustomer)
	jsonOutput := string(redact.JSON(originalCustomer))
	fallbackJSON := string(redact.JSONBytes(common.MalformedJSON()))
	stringOutput := fmt.Sprintf(
		"String:   %s\nEmail:    %s\nPhone:    %s\nPostcode: %s",
		redact.String("greg@example.com"),
		redact.Email("greg@example.com"),
		redact.Phone("07700 900123"),
		redact.Postcode("TN9 1XA"),
	)

	return []Step{
		{
			Title:       "Step 1: Structured redaction",
			Description: "Start with a nested customer object and show how redact.Data preserves shape while masking sensitive fields.",
			InputTitle:  "Original object",
			InputText:   common.PrettyJSON(originalCustomer),
			OutputTitle: "redact.Data output",
			OutputText:  common.PrettyJSON(redactedCustomer),
		},
		{
			Title:       "Step 2: JSON export",
			Description: "Show the same customer rendered through redact.JSON for export or logging use cases.",
			InputTitle:  "Original object",
			InputText:   common.PrettyJSON(originalCustomer),
			OutputTitle: "redact.JSON output",
			OutputText:  jsonOutput,
		},
		{
			Title:       "Step 3: Malformed JSON fallback",
			Description: "Invalid raw bytes never panic and always return a valid fallback payload.",
			InputTitle:  "Malformed JSON input",
			InputText:   string(common.MalformedJSON()),
			OutputTitle: "redact.JSONBytes output",
			OutputText:  fallbackJSON,
		},
		{
			Title:       "Step 4: Direct string helpers",
			Description: "Use direct value helpers when you already know the data type you are redacting.",
			InputTitle:  "Helper inputs",
			InputText:   "greg@example.com\n07700 900123\nTN9 1XA",
			OutputTitle: "Helper outputs",
			OutputText:  stringOutput,
		},
	}
}
