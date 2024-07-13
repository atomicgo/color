package color

import (
	"fmt"
	"strings"
)

type Style struct {
	Foreground Color
	Background Color

	Modifiers []Modifier
}

func NewStyle(foregroundColor, backgroundColor Color, modifiers ...Modifier) Style {
	if foregroundColor == nil {
		foregroundColor = NoColor
	}

	if backgroundColor == nil {
		backgroundColor = NoColor
	}

	return Style{
		Foreground: foregroundColor,
		Background: backgroundColor,
		Modifiers:  modifiers,
	}
}

func (s Style) Sequence() string {
	var sb strings.Builder

	if s.Foreground != nil {
		sb.WriteString(s.Foreground.Sequence(false))
	}

	if s.Background != nil {
		sb.WriteString(s.Background.Sequence(true))
	}

	for _, m := range s.Modifiers {
		sb.WriteString(m.Sequence())
	}

	return sb.String()
}

func (s Style) Sprint(a ...any) string {
	var sb strings.Builder

	sb.WriteString(s.Sequence())
	sb.WriteString(fmt.Sprint(a...))
	sb.WriteString(Reset.Sequence())

	return sb.String()
}
