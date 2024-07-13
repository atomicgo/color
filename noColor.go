package color

type noColor struct{}

func (noColor) Sequence(_ bool) string {
	return ""
}

var NoColor Color = noColor{}
