package tags

import "strings"

func ParseSerialiseTag(tag string) (string, bool) {
	if tag == "-" {
		return "", true
	}

	if tag == "" {
		return "", false
	}

	parts := strings.Split(tag, ",")
	if parts[0] == "-" {
		return "", true
	}

	if parts[0] == "" {
		return "", false
	}

	return parts[0], false
}
