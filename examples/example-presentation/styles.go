package main

import (
	"fmt"
	"strings"

	"github.com/BreakPointSoftware/annon/examples/common"
)

func renderStep(stepNumber, totalSteps int, step Step) string {
	return common.Section(step.Title) + "\n" +
		fmt.Sprintf("Step %d of %d\n\n", stepNumber, totalSteps) +
		step.Description + "\n\n" +
		renderPanel(step.InputTitle, step.InputText) + "\n\n" +
		renderPanel(step.OutputTitle, step.OutputText)
}

func renderPanel(title, body string) string {
	return title + "\n" + strings.Repeat("-", len(title)) + "\n" + body
}
