package walker

import "strings"

type Path []string

func (p Path) Append(part string) Path {
	clone := append(Path{}, p...)
	return append(clone, part)
}

func (p Path) String() string {
	return strings.Join(p, ".")
}
