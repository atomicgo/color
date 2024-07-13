package color

import "fmt"

const (
	ANSIBlack ANSIColor = iota
	ANSIRed
	ANSIGreen
	ANSIYellow
	ANSIBlue
	ANSIMagenta
	ANSICyan
	ANSIWhite
	ANSIBrightBlack
	ANSIBrightRed
	ANSIBrightGreen
	ANSIBrightYellow
	ANSIBrightBlue
	ANSIBrightMagenta
	ANSIBrightCyan
	ANSIBrightWhite
)

// ANSIColor represents an ANSI color code.
type ANSIColor int

// Sequence represents the ANSI escape sequence for the color.
func (c ANSIColor) Sequence(background bool) string {
	if background {
		return fmt.Sprintf("\033[48;5;%dm", c)
	}

	return fmt.Sprintf("\033[38;5;%dm", c)
}
