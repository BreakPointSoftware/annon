package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func PrettyJSON(value any) string {
	encodedBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", value)
	}

	return string(encodedBytes)
}

func Section(title string) string {
	divider := strings.Repeat("=", 60)
	return divider + "\n" + title + "\n" + divider
}

func WaitForEnter() {
	fmt.Print("\nPress Enter to continue...")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
