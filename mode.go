package color

type Mode int

const (
	TrueColor Mode = iota
	ANSI256
	ANSI

	Disabled
)

func (m Mode) String() string {
	switch m {
	case TrueColor:
		return "TrueColor"
	case ANSI256:
		return "ANSI256"
	case ANSI:
		return "ANSI"
	case Disabled:
		return "Disabled"
	default:
		return "Unknown"
	}
}
