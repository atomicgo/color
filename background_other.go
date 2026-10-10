//go:build !unix

package color

import "time"

// queryBackground has no terminal query on this platform.
func queryBackground(time.Duration) (Color, bool) {
	return None, false
}
