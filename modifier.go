package color

import "fmt"

// Modifier type for text modifiers.
type Modifier int

const modifierSequence = "\033[%dm"

// Modifiers
const (
	Reset Modifier = iota
	Bold
	Italic
	Underline
)

// Sequence returns the ANSI escape sequence for the modifier.
func (m Modifier) Sequence() string {
	return fmt.Sprintf(modifierSequence, m)
}
