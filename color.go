package color

import (
	"fmt"
)

type Color interface {
	Sequence(background bool) string
}

func Sprint(a ...any) string {
	return fmt.Sprint(a...)
}
