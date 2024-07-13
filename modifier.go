package color

import "fmt"

// Modifier type for text modifiers
type Modifier int

const modifierSequence = "\033[%dm"

// Modifiers
const (
	Reset Modifier = iota
	Bold
	Italic
	Underline
)

func (m Modifier) Sequence() string {
	return fmt.Sprintf(modifierSequence, m)
}
