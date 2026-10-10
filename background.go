package color

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Background returns the background color of the terminal, for example to
// choose between a dark and a light theme with Color.IsDark.
//
// It reads COLORFGBG first. When that is unset it asks the terminal with an
// OSC 11 query and waits at most timeout for the answer. The query needs a
// controlling terminal on Unix; on Windows only COLORFGBG is read. The bool
// is false when the background could not be found out.
//
// The query reads from the terminal, so call it before reading user input,
// not while another goroutine does. Keys pressed during the query are lost.
func Background(timeout time.Duration) (Color, bool) {
	if c, ok := parseColorFgBg(os.Getenv("COLORFGBG")); ok {
		return c, true
	}

	return queryBackground(timeout)
}

// parseColorFgBg reads the background from COLORFGBG, which rxvt, Konsole
// and others set to "fg;bg" or "fg;default;bg" with basic color indexes.
func parseColorFgBg(v string) (Color, bool) {
	if v == "" {
		return None, false
	}

	parts := strings.Split(v, ";")

	n, err := strconv.ParseUint(parts[len(parts)-1], 10, 8)
	if err != nil || n > 15 {
		return None, false
	}

	return Basic(uint8(n)), true
}

// parseOSC11 extracts the color from a terminal's answer to an OSC 11 query,
// such as "\x1b]11;rgb:1e1e/1e1e/2e2e\x1b\\". Each channel has 1 to 4 hex
// digits.
func parseOSC11(resp string) (Color, bool) {
	_, after, ok := strings.Cut(resp, "]11;")
	if !ok {
		return None, false
	}

	end := strings.IndexAny(after, "\x07\x1b")
	if end < 0 {
		return None, false
	}

	spec := after[:end]
	if !strings.HasPrefix(spec, "rgb:") && !strings.HasPrefix(spec, "rgba:") {
		return None, false
	}

	_, spec, _ = strings.Cut(spec, ":")

	channels := strings.Split(spec, "/")
	if len(channels) < 3 {
		return None, false
	}

	var rgb [3]uint8

	for i := range 3 {
		ch := channels[i]
		if ch == "" || len(ch) > 4 {
			return None, false
		}

		v, err := strconv.ParseUint(ch, 16, 16)
		if err != nil {
			return None, false
		}

		maxValue := uint64(1)<<(4*len(ch)) - 1
		rgb[i] = uint8(v * 255 / maxValue) //nolint:gosec // scaled into 0-255.
	}

	return RGB(rgb[0], rgb[1], rgb[2]), true
}
