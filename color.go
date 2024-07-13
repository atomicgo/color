package color

import (
	"io"
	"os"
)

// Writer is the writer to write colorized output to.
var Writer io.Writer = os.Stdout

// Color is an interface for colors.
type Color interface {
	Sequence(background bool) string
}
