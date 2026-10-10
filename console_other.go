//go:build !windows

package color

// consoleSupport reports the color support the platform console guarantees.
// Outside of Windows the console needs no setup and has no say, so the
// environment decides.
func consoleSupport() (Level, bool) {
	return LevelNone, false
}
