package walk

import "strings"

type parsedTag struct {
	empty        bool
	skip         bool
	auto         bool
	remove       bool
	strategyName string
}

func parseTag(raw string) parsedTag {
	tag := strings.TrimSpace(raw)
	switch tag {
	case "", "-":
		return parsedTag{empty: true}
	case "false":
		return parsedTag{skip: true}
	case "true", "auto":
		return parsedTag{auto: true}
	case "remove":
		return parsedTag{remove: true, strategyName: "remove"}
	default:
		return parsedTag{strategyName: tag}
	}
}
