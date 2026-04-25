package color

type noColor struct{}

func (noColor) Sequence(_ bool) string {
	return ""
}

// NoColor disables color output for a style component.
var NoColor Color = noColor{}
