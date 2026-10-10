package color

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// kind tells which palette a Color comes from.
type kind uint8

const (
	kindNone kind = iota
	kindBasic
	kind256
	kindRGB
)

// Color is a single terminal color: one of the 16 basic ANSI colors, an
// entry of the xterm 256-color palette, or a 24-bit RGB value.
//
// The zero value is None, which renders nothing. Colors are comparable and
// safe to copy.
type Color struct {
	kind    kind
	r, g, b uint8 // RGB value, or the palette index in r for basic and 256
}

// None is the absent color. Styles skip it, so a Style with a None
// background leaves the terminal's background alone.
var None = Color{}

// The 16 basic ANSI colors. Terminals render them from their own palette,
// so they follow the user's terminal theme.
var (
	Black         = Basic(0)
	Red           = Basic(1)
	Green         = Basic(2)
	Yellow        = Basic(3)
	Blue          = Basic(4)
	Magenta       = Basic(5)
	Cyan          = Basic(6)
	White         = Basic(7)
	BrightBlack   = Basic(8)
	BrightRed     = Basic(9)
	BrightGreen   = Basic(10)
	BrightYellow  = Basic(11)
	BrightBlue    = Basic(12)
	BrightMagenta = Basic(13)
	BrightCyan    = Basic(14)
	BrightWhite   = Basic(15)
)

var basicNames = [16]string{
	"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
	"brightblack", "brightred", "brightgreen", "brightyellow",
	"brightblue", "brightmagenta", "brightcyan", "brightwhite",
}

// ErrInvalidColor is returned by Parse and Hex for input that names no color.
var ErrInvalidColor = errors.New("color: invalid color")

// Basic returns basic ANSI color n (0-15). Values above 15 wrap around.
func Basic(n uint8) Color {
	return Color{kind: kindBasic, r: n % 16}
}

// Index returns entry n of the xterm 256-color palette.
func Index(n uint8) Color {
	return Color{kind: kind256, r: n}
}

// RGB returns a 24-bit color. Terminals without true color support get the
// closest color they can show.
func RGB(r, g, b uint8) Color {
	return Color{kind: kindRGB, r: r, g: g, b: b}
}

// Hex parses "#rrggbb", "rrggbb", "#rgb" or "rgb" into an RGB color.
func Hex(s string) (Color, error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}

	if len(h) != 6 {
		return None, fmt.Errorf("%w: %q", ErrInvalidColor, s)
	}

	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return None, fmt.Errorf("%w: %q", ErrInvalidColor, s)
	}

	return RGB(uint8(v>>16), uint8(v>>8), uint8(v)), nil //nolint:gosec // v has 24 bits, each shift keeps one byte.
}

// MustHex is like Hex but panics on invalid input. Use it for constants in
// your own code, never for user input.
func MustHex(s string) Color {
	c, err := Hex(s)
	if err != nil {
		panic(err)
	}

	return c
}

// Parse reads a color the way a user would write it in a config file or an
// environment variable:
//
//   - a hex value: "#b7bdf8", "b7bdf8", "#fff"
//   - a basic color name: "red", "brightblue", "bright-blue", "bright_blue"
//   - a 256-color palette index: "0" to "255"
//   - "none" or "" for None
func Parse(s string) (Color, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "", "none", "default":
		return None, nil
	}

	if isDigits(s) {
		n, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return None, fmt.Errorf("%w: %q is not a palette index (0-255)", ErrInvalidColor, s)
		}

		return Index(uint8(n)), nil
	}

	name := strings.NewReplacer("-", "", "_", "", " ", "").Replace(s)
	if name == "gray" || name == "grey" {
		name = "brightblack"
	}

	for i, n := range basicNames {
		if name == n {
			return Basic(uint8(i)), nil
		}
	}

	return Hex(s)
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return s != ""
}

// IsNone reports whether c is the absent color.
func (c Color) IsNone() bool {
	return c.kind == kindNone
}

// RGB returns the 24-bit value of c. Basic colors report the xterm default
// palette, because the real value depends on the user's terminal theme.
// None reports black.
func (c Color) RGB() (r, g, b uint8) {
	switch c.kind {
	case kindRGB:
		return c.r, c.g, c.b
	case kindBasic:
		p := basicPalette[c.r]
		return uint8(p[0]), uint8(p[1]), uint8(p[2]) //nolint:gosec // palette values are bytes.
	case kind256:
		return index256ToRGB(c.r)
	default:
		return 0, 0, 0
	}
}

// Hex returns c as "#rrggbb". See Color.RGB for basic and palette colors.
func (c Color) Hex() string {
	r, g, b := c.RGB()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// String describes c in the format Parse accepts.
func (c Color) String() string {
	switch c.kind {
	case kindBasic:
		return basicNames[c.r]
	case kind256:
		return strconv.Itoa(int(c.r))
	case kindRGB:
		return c.Hex()
	default:
		return "none"
	}
}

// Mix blends c toward other. t=0 returns c, t=1 returns other, and values in
// between interpolate in RGB space. Use it for gradients and for tints of a
// theme color. The result is always an RGB color.
func (c Color) Mix(other Color, t float64) Color {
	t = min(max(t, 0), 1)
	r1, g1, b1 := c.RGB()
	r2, g2, b2 := other.RGB()

	lerp := func(a, b uint8) uint8 {
		return uint8(float64(a) + (float64(b)-float64(a))*t + 0.5)
	}

	return RGB(lerp(r1, r2), lerp(g1, g2), lerp(b1, b2))
}

// Luminance returns the relative luminance of c from 0 (black) to 1 (white),
// as defined by WCAG.
func (c Color) Luminance() float64 {
	r, g, b := c.RGB()

	return 0.2126*linear(r) + 0.7152*linear(g) + 0.0722*linear(b)
}

// IsDark reports whether c is a dark color. Use it to pick readable text on
// top of c, or to tell a dark terminal background from a light one.
func (c Color) IsDark() bool {
	return c.Luminance() < 0.2
}

func linear(v uint8) float64 {
	f := float64(v) / 255
	if f <= 0.04045 {
		return f / 12.92
	}

	return math.Pow((f+0.055)/1.055, 2.4)
}

// Fg returns a Style with c as its foreground.
func (c Color) Fg() Style {
	return Style{Foreground: c}
}

// Bg returns a Style with c as its background.
func (c Color) Bg() Style {
	return Style{Background: c}
}

// Sprint colors the text with c as foreground for standard output.
func (c Color) Sprint(a ...any) string {
	return c.Fg().Sprint(a...)
}

// Sprintf formats and colors the text with c as foreground for standard
// output.
func (c Color) Sprintf(format string, a ...any) string {
	return c.Fg().Sprintf(format, a...)
}
