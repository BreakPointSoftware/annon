package main

import (
	"fmt"

	"github.com/BreakPointSoftware/annon/examples/common"
)

func main() {
	steps := buildSteps()

	for index, step := range steps {
		fmt.Println(renderStep(index+1, len(steps), step))
		if index < len(steps)-1 {
			common.WaitForEnter()
			fmt.Println()
		}
	}
}
