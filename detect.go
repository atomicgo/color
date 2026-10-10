package color

import (
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/term"
)

// fder is implemented by *os.File and by writers that wrap a file descriptor.
type fder interface {
	Fd() uintptr
}

var (
	// override holds the level set by SetLevel, or nil.
	override atomic.Pointer[Level]

	// ttyCache remembers which file descriptors are terminals, because the
	// check costs a syscall and Detect runs for every styled print.
	ttyCache sync.Map // map[uintptr]bool
)

// SetLevel makes Detect return l for every writer, ignoring the environment
// and the terminal. Use it for a --color flag or in tests. ResetLevel undoes
// it. Safe for concurrent use.
func SetLevel(l Level) {
	override.Store(&l)
}

// ResetLevel undoes SetLevel, so Detect inspects the environment and the
// writer again.
func ResetLevel() {
	override.Store(nil)
}

// Detect reports how many colors w can show. Pass the writer you are about
// to write to: standard output can be a terminal while standard error is
// piped into a file, or the other way around.
//
// The checks, in order:
//
//  1. A level set with SetLevel wins.
//  2. FORCE_COLOR forces colors on, even for pipes and files: "1" or "true"
//     for at least 16 colors, "2" for 256, "3" for true color. "0", "false"
//     and "none" force colors off. CLICOLOR_FORCE=1 also forces colors on.
//  3. NO_COLOR (any value, see https://no-color.org) and CLICOLOR=0 turn
//     colors off.
//  4. A writer that is not a terminal (a pipe, a file, a bytes.Buffer) gets
//     LevelNone.
//  5. TERM=dumb gets LevelNone. COLORTERM, TERM, TERM_PROGRAM and a few
//     terminal-specific variables decide between 16, 256 and true color.
//     On Windows the first call turns on ANSI processing in the console, and
//     consoles that cannot do it get LevelNone.
func Detect(w io.Writer) Level {
	if l := override.Load(); l != nil {
		return *l
	}

	console, consoleKnown := consoleSupport()

	return detect(os.Getenv, isTerminal(w), console, consoleKnown)
}

// DetectEnv is Detect for a terminal you know is there: it applies every
// check except the terminal check of step 4.
func DetectEnv() Level {
	if l := override.Load(); l != nil {
		return *l
	}

	console, consoleKnown := consoleSupport()

	return detect(os.Getenv, true, console, consoleKnown)
}

// IsTerminal reports whether w writes to a terminal.
func IsTerminal(w io.Writer) bool {
	return isTerminal(w)
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(fder)
	if !ok || f == nil {
		return false
	}

	fd := f.Fd()
	if v, ok := ttyCache.Load(fd); ok {
		return v.(bool) //nolint:forcetypeassert // ttyCache only stores bools.
	}

	tty := term.IsTerminal(int(fd))
	ttyCache.Store(fd, tty)

	return tty
}

// detect is the testable core of Detect: every input is explicit.
func detect(getenv func(string) string, tty bool, console Level, consoleKnown bool) Level {
	if force := getenv("FORCE_COLOR"); force != "" {
		switch strings.ToLower(force) {
		case "0", "false", "none":
			return LevelNone
		case "2":
			return Level256
		case "3":
			return LevelTrueColor
		default:
			return max(envLevel(getenv, console, consoleKnown), LevelBasic)
		}
	}

	if force := getenv("CLICOLOR_FORCE"); force != "" && force != "0" {
		return max(envLevel(getenv, console, consoleKnown), LevelBasic)
	}

	if getenv("NO_COLOR") != "" || getenv("CLICOLOR") == "0" {
		return LevelNone
	}

	if !tty {
		return LevelNone
	}

	return envLevel(getenv, console, consoleKnown)
}

// envLevel reads the color level from variables that describe the terminal
// without forcing anything, plus the platform console verdict.
func envLevel(getenv func(string) string, console Level, consoleKnown bool) Level {
	colorTerm := strings.ToLower(getenv("COLORTERM"))
	term := strings.ToLower(getenv("TERM"))

	if strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit") {
		return LevelTrueColor
	}

	if strings.Contains(term, "truecolor") || strings.Contains(term, "24bit") || strings.Contains(term, "direct") {
		return LevelTrueColor
	}

	if term == "dumb" {
		return LevelNone
	}

	// These terminals show true color but do not always set COLORTERM.
	if getenv("WT_SESSION") != "" || getenv("ConEmuANSI") == "ON" || getenv("TERMINAL_EMULATOR") == "JetBrains-JediTerm" {
		return LevelTrueColor
	}

	switch getenv("TERM_PROGRAM") {
	case "iTerm.app", "Hyper", "vscode", "ghostty", "WezTerm", "rio", "Tabby":
		return LevelTrueColor
	case "Apple_Terminal":
		// Older Terminal.app versions garble true color, and 256 colors
		// look right on every version.
		return Level256
	}

	if strings.Contains(term, "kitty") || strings.Contains(term, "alacritty") {
		return LevelTrueColor
	}

	if strings.Contains(term, "256color") {
		return Level256
	}

	// ANSICON translates ANSI sequences for legacy Windows consoles.
	if getenv("ANSICON") != "" {
		return LevelBasic
	}

	// On Windows the console decides: true color when it processes ANSI
	// sequences, nothing when it cannot (cmd.exe before Windows 10).
	if consoleKnown {
		return console
	}

	return LevelBasic
}
