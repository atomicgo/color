package color

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Attr is a set of text attributes. Combine them with |, for example
// Bold|Underline.
type Attr uint16

// Text attributes. Terminals that do not support one ignore it.
const (
	Bold Attr = 1 << iota
	Dim
	Italic
	Underline
	Blink
	Reverse
	Hidden
	Strikethrough
	Overline
)

// attrCodes lists the SGR code of each attribute in bit order.
var attrCodes = [...]string{"1", "2", "3", "4", "5", "7", "8", "9", "53"}

const (
	esc = "\x1b"

	// Reset is the escape sequence that turns off every color and attribute.
	Reset = esc + "[0m"
)

// Style is a foreground color, a background color and text attributes.
// The zero value renders text unchanged.
//
// Styles are values: every method returns a new Style and leaves the
// receiver alone, so a Style can be shared between goroutines.
type Style struct {
	Foreground Color
	Background Color
	Attrs      Attr
}

// New returns a Style with the given attributes and no colors.
func New(attrs ...Attr) Style {
	var s Style
	for _, a := range attrs {
		s.Attrs |= a
	}

	return s
}

// Fg returns a copy of s with c as foreground.
func (s Style) Fg(c Color) Style {
	s.Foreground = c
	return s
}

// Bg returns a copy of s with c as background.
func (s Style) Bg(c Color) Style {
	s.Background = c
	return s
}

// Add returns a copy of s with the given attributes turned on.
func (s Style) Add(attrs ...Attr) Style {
	for _, a := range attrs {
		s.Attrs |= a
	}

	return s
}

// Bold returns a copy of s with Bold turned on.
func (s Style) Bold() Style { return s.Add(Bold) }

// Dim returns a copy of s with Dim turned on.
func (s Style) Dim() Style { return s.Add(Dim) }

// Italic returns a copy of s with Italic turned on.
func (s Style) Italic() Style { return s.Add(Italic) }

// Underline returns a copy of s with Underline turned on.
func (s Style) Underline() Style { return s.Add(Underline) }

// Strikethrough returns a copy of s with Strikethrough turned on.
func (s Style) Strikethrough() Style { return s.Add(Strikethrough) }

// Reverse returns a copy of s with Reverse turned on.
func (s Style) Reverse() Style { return s.Add(Reverse) }

// IsZero reports whether s has no colors and no attributes.
func (s Style) IsZero() bool {
	return s == Style{}
}

// Sequence returns the escape sequence that switches the terminal to s at
// level l. It returns "" when l is LevelNone or s is zero.
func (s Style) Sequence(l Level) string {
	if l == LevelNone {
		return ""
	}

	codes := make([]string, 0, 4)

	for i, code := range attrCodes {
		if s.Attrs&(1<<i) != 0 {
			codes = append(codes, code)
		}
	}

	if c := l.fgCode(s.Foreground); c != "" {
		codes = append(codes, c)
	}

	if c := l.bgCode(s.Background); c != "" {
		codes = append(codes, c)
	}

	if len(codes) == 0 {
		return ""
	}

	return esc + "[" + strings.Join(codes, ";") + "m"
}

// Render returns text styled with s at level l.
//
// Styled text can be nested: a reset inside text switches back to s instead
// of plain text, so an inner Render does not end the outer style. Each line
// is styled on its own, so a background color does not bleed into the
// margin when the text is printed next to something else.
func (s Style) Render(l Level, text string) string {
	seq := s.Sequence(l)
	if seq == "" || text == "" {
		return text
	}

	text = strings.ReplaceAll(text, Reset, Reset+seq)

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = seq + line + Reset
		}
	}

	return strings.Join(lines, "\n")
}

// Sprint styles the text for standard output, at the level Detect(os.Stdout)
// reports.
func (s Style) Sprint(a ...any) string {
	return s.Render(Detect(os.Stdout), fmt.Sprint(a...))
}

// Sprintf formats and styles the text for standard output.
func (s Style) Sprintf(format string, a ...any) string {
	return s.Render(Detect(os.Stdout), fmt.Sprintf(format, a...))
}

// Fprint styles the text at the level Detect(w) reports and writes it to w.
func (s Style) Fprint(w io.Writer, a ...any) (int, error) {
	return io.WriteString(w, s.Render(Detect(w), fmt.Sprint(a...)))
}

// Fprintln is like Fprint and adds a newline.
func (s Style) Fprintln(w io.Writer, a ...any) (int, error) {
	return io.WriteString(w, s.Render(Detect(w), fmt.Sprint(a...))+"\n")
}

// Fprintf formats, styles and writes the text to w.
func (s Style) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return io.WriteString(w, s.Render(Detect(w), fmt.Sprintf(format, a...)))
}

// Print writes the styled text to standard output.
func (s Style) Print(a ...any) {
	_, _ = s.Fprint(os.Stdout, a...)
}

// Println writes the styled text and a newline to standard output.
func (s Style) Println(a ...any) {
	_, _ = s.Fprintln(os.Stdout, a...)
}

// Printf formats and writes the styled text to standard output.
func (s Style) Printf(format string, a ...any) {
	_, _ = s.Fprintf(os.Stdout, format, a...)
}

// ansiRegex matches CSI sequences (colors, cursor movement) and OSC
// sequences (hyperlinks, titles) terminated by BEL or ST.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;:?<=>]*[ -/]*[@-~]|\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)`)

// Strip removes every escape sequence from s: colors, attributes, cursor
// movement and hyperlinks. The text of a hyperlink stays.
func Strip(s string) string {
	if !strings.Contains(s, esc) {
		return s
	}

	return ansiRegex.ReplaceAllString(s, "")
}

// Hyperlink returns text as a clickable link to url (OSC 8). At LevelNone it
// returns text unchanged. Terminals without link support show the text and
// ignore the link.
func Hyperlink(l Level, url, text string) string {
	if l == LevelNone || url == "" {
		return text
	}

	return esc + "]8;;" + url + esc + "\\" + text + esc + "]8;;" + esc + "\\"
}

// SGR returns the escape sequence for raw SGR parameters, for the rare
// attribute this package has no name for, such as SGR(21) for a double
// underline.
func SGR(params ...int) string {
	codes := make([]string, len(params))
	for i, p := range params {
		codes[i] = strconv.Itoa(p)
	}

	return esc + "[" + strings.Join(codes, ";") + "m"
}
