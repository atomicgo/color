package color

import (
	"fmt"
	"io"
	"strings"
)

// Style represents a text style with a foreground and background color and modifiers.
type Style struct {
	Foreground Color
	Background Color

	Modifiers []Modifier
}

// NewStyle creates a new Style with the given foreground and background colors and modifiers.
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

// Sequence returns the ANSI escape sequence for the style.
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

// Sprint formats using the default formats for its operands and returns the resulting string.
func (s Style) Sprint(a ...any) string {
	var sb strings.Builder

	sb.WriteString(s.Sequence())
	sb.WriteString(fmt.Sprint(a...))
	sb.WriteString(Reset.Sequence())

	return sb.String()
}

// Sprintf formats according to a format specifier and returns the resulting string.
func (s Style) Sprintf(format string, a ...any) string {
	return s.Sprint(fmt.Sprintf(format, a...))
}

// Fprint formats using the default formats for its operands and writes to w.
func (s Style) Fprint(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprint(w, s.Sprint(a...)) //nolint:wrapcheck // stdlib error
}

// Fprintf formats according to a format specifier and writes to w.
func (s Style) Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
	return fmt.Fprint(w, s.Sprintf(format, a...)) //nolint:wrapcheck // stdlib error
}

// Fprintln formats using the default formats for its operands and writes to w.
func (s Style) Fprintln(w io.Writer, a ...any) (n int, err error) {
	return fmt.Fprintln(w, s.Sprint(a...)) //nolint:wrapcheck // stdlib error
}

// Fprintfln formats according to a format specifier and writes to w.
func (s Style) Fprintfln(w io.Writer, format string, a ...any) (n int, err error) {
	return fmt.Fprintln(w, s.Sprintf(format, a...)) //nolint:wrapcheck // stdlib error
}

// Print formats using the default formats for its operands and writes to standard output.
func (s Style) Print(a ...any) {
	_, _ = s.Fprint(Writer, a...)
}

// Printf formats according to a format specifier and writes to standard output.
func (s Style) Printf(format string, a ...any) {
	_, _ = s.Fprintf(Writer, format, a...)
}

// Println formats using the default formats for its operands and writes to standard output.
func (s Style) Println(a ...any) {
	_, _ = s.Fprintln(Writer, a...)
}

// Printfln formats according to a format specifier and writes to standard output.
func (s Style) Printfln(format string, a ...any) {
	_, _ = s.Fprintfln(Writer, format, a...)
}
