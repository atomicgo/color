package color

import (
	"io"
	"os"
)

var Writer io.Writer = os.Stdout

type Color interface {
	Sequence(background bool) string
}
