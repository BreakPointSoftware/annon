package output

type Format string

const (
	JSON Format = "json"
	YAML Format = "yaml"
)

func (f Format) String() string {
	return string(f)
}
