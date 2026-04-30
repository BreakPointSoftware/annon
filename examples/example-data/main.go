package main

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/examples/common"
	"github.com/BreakPointSoftware/annon/redact"
)

func buildOutput() string {
	originalCustomer := common.DemoCustomer()
	redactedCustomer := redact.Data(originalCustomer)

	return common.Section("Structured redaction with redact.Data") + "\n\n" +
		"Original\n--------\n" + common.PrettyJSON(originalCustomer) + "\n\n" +
		"Redacted\n--------\n" + common.PrettyJSON(redactedCustomer) + "\n"
}

func main() {
	fmt.Print(buildOutput())
}
