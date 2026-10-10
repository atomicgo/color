package color

import "strconv"

// Level is how many colors a terminal can show.
type Level uint8

const (
	// LevelNone shows no colors and no text attributes. Rendering at this
	// level returns the text unchanged.
	LevelNone Level = iota

	// LevelBasic shows the 16 basic ANSI colors.
	LevelBasic

	// Level256 shows the xterm 256-color palette.
	Level256

	// LevelTrueColor shows 24-bit RGB colors.
	LevelTrueColor
)

// String returns "none", "16", "256" or "truecolor".
func (l Level) String() string {
	switch l {
	case LevelNone:
		return "none"
	case LevelBasic:
		return "16"
	case Level256:
		return "256"
	case LevelTrueColor:
		return "truecolor"
	default:
		return "unknown"
	}
}

// fgCode returns the SGR parameters that set c as foreground at level l, or
// "" when l cannot show colors or c is None.
func (l Level) fgCode(c Color) string {
	return l.code(c, false)
}

// bgCode is fgCode for the background.
func (l Level) bgCode(c Color) string {
	return l.code(c, true)
}

func (l Level) code(c Color, background bool) string {
	if l == LevelNone || c.kind == kindNone {
		return ""
	}

	// A basic color keeps its basic code at every level, so it follows the
	// user's terminal palette. Richer colors fall back to what l can show.
	switch {
	case c.kind == kindBasic:
		return basicCode(c.r, background)
	case c.kind == kind256 && l >= Level256:
		return layer(background) + ";5;" + strconv.Itoa(int(c.r))
	case c.kind == kindRGB && l == LevelTrueColor:
		return layer(background) + ";2;" + strconv.Itoa(int(c.r)) + ";" + strconv.Itoa(int(c.g)) + ";" + strconv.Itoa(int(c.b))
	case l == Level256:
		r, g, b := c.RGB()
		return layer(background) + ";5;" + strconv.Itoa(int(rgbTo256(r, g, b)))
	default:
		r, g, b := c.RGB()
		return basicCode(rgbToBasic(r, g, b), background)
	}
}

func layer(background bool) string {
	if background {
		return "48"
	}

	return "38"
}

// basicCode returns the SGR code of basic color n: 30-37 and 90-97 for
// foregrounds, 40-47 and 100-107 for backgrounds.
func basicCode(n uint8, background bool) string {
	base := 30
	if n >= 8 {
		base = 90
		n -= 8
	}

	if background {
		base += 10
	}

	return strconv.Itoa(base + int(n))
}
