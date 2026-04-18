package color

import (
	"fmt"
	"strconv"
	"strings"
)

// RGBColor represents a color in the RGB color space.
type RGBColor struct {
	R, G, B uint8
}

// Hex returns the hex representation of the color.
func (c RGBColor) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// NewColorFromRGB creates a new Color from RGB values.
func NewColorFromRGB(r, g, b uint8) Color {
	return RGBColor{r, g, b}
}

// NewColorFromHex creates a new Color from a hex string.
// Accepts both the 6-digit form ("#RRGGBB"/"RRGGBB") and the 3-digit
// shorthand ("#RGB"/"RGB"). If the hex string is invalid, NoColor is returned.
func NewColorFromHex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")

	// Expand 3-digit shorthand (e.g. "f0a" -> "ff00aa").
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	if len(hex) != 6 {
		return NoColor
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return NoColor
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return NoColor
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return NoColor
	}

	return NewColorFromRGB(uint8(r), uint8(g), uint8(b))
}

// Sequence returns the ANSI escape sequence for the color.
func (c RGBColor) Sequence(background bool) string {
	if background {
		return fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}
