package color

import (
	"bytes"
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want Color
	}{
		{"#b7bdf8", RGB(0xb7, 0xbd, 0xf8)},
		{"B7BDF8", RGB(0xb7, 0xbd, 0xf8)},
		{"#fff", RGB(255, 255, 255)},
		{"red", Red},
		{"Bright-Blue", BrightBlue},
		{"bright_white", BrightWhite},
		{"gray", BrightBlack},
		{"208", Index(208)},
		{"0", Index(0)},
		{"", None},
		{"none", None},
	}

	for _, tt := range tests {
		got, err := Parse(tt.in)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.in, err)
			continue
		}

		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}

	for _, bad := range []string{"#12", "purple-ish", "256", "#ggg"} {
		if _, err := Parse(bad); !errors.Is(err, ErrInvalidColor) {
			t.Errorf("Parse(%q) error = %v, want ErrInvalidColor", bad, err)
		}
	}
}

func TestColorString(t *testing.T) {
	t.Parallel()

	for _, c := range []Color{Red, BrightCyan, Index(42), RGB(1, 2, 3), None} {
		back, err := Parse(c.String())
		if err != nil || back != c {
			t.Errorf("Parse(%q) = %v, %v; want %v", c.String(), back, err, c)
		}
	}
}

func TestSequenceDownsamples(t *testing.T) {
	t.Parallel()

	lavender := RGB(0xb7, 0xbd, 0xf8)

	tests := []struct {
		name  string
		style Style
		level Level
		want  string
	}{
		{"truecolor", lavender.Fg(), LevelTrueColor, "\x1b[38;2;183;189;248m"},
		{"256", lavender.Fg(), Level256, "\x1b[38;5;147m"},
		{"basic", lavender.Fg(), LevelBasic, "\x1b[37m"},
		{"none", lavender.Fg(), LevelNone, ""},
		{"basic color stays basic", Red.Fg(), LevelTrueColor, "\x1b[31m"},
		{"bright background", BrightBlue.Bg(), LevelBasic, "\x1b[104m"},
		{"index at truecolor", Index(208).Fg(), LevelTrueColor, "\x1b[38;5;208m"},
		{"index at basic", Index(196).Fg(), LevelBasic, "\x1b[91m"},
		{"attrs and colors", New(Bold, Underline).Fg(Green).Bg(Black), LevelBasic, "\x1b[1;4;32;40m"},
		{"zero", Style{}, LevelTrueColor, ""},
	}

	for _, tt := range tests {
		if got := tt.style.Sequence(tt.level); got != tt.want {
			t.Errorf("%s: Sequence = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestRenderNests(t *testing.T) {
	t.Parallel()

	inner := Red.Fg().Render(LevelBasic, "inner")
	got := New(Bold).Render(LevelBasic, "a "+inner+" b")
	want := "\x1b[1ma \x1b[31minner\x1b[0m\x1b[1m b\x1b[0m"

	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}
}

func TestRenderPerLine(t *testing.T) {
	t.Parallel()

	got := Red.Bg().Render(LevelBasic, "a\n\nb")
	want := "\x1b[41ma\x1b[0m\n\n\x1b[41mb\x1b[0m"

	if got != want {
		t.Errorf("Render = %q, want %q", got, want)
	}

	if got := Red.Fg().Render(LevelNone, "a"); got != "a" {
		t.Errorf("Render at LevelNone = %q, want plain text", got)
	}
}

func TestStrip(t *testing.T) {
	t.Parallel()

	styled := New(Bold).Fg(Red).Render(LevelTrueColor, "hi") + Hyperlink(LevelBasic, "https://pterm.sh", "link") + "\x1b[2K"
	if got := Strip(styled); got != "hilink" {
		t.Errorf("Strip = %q, want %q", got, "hilink")
	}
}

func TestDetect(t *testing.T) {
	t.Parallel()

	env := func(vars map[string]string) func(string) string {
		return func(k string) string { return vars[k] }
	}

	tests := []struct {
		name string
		vars map[string]string
		tty  bool
		want Level
	}{
		{"pipe", map[string]string{"COLORTERM": "truecolor"}, false, LevelNone},
		{"pipe forced", map[string]string{"FORCE_COLOR": "1"}, false, LevelBasic},
		{"pipe forced 3", map[string]string{"FORCE_COLOR": "3"}, false, LevelTrueColor},
		{"force off wins", map[string]string{"FORCE_COLOR": "0", "COLORTERM": "truecolor"}, true, LevelNone},
		{"clicolor force", map[string]string{"CLICOLOR_FORCE": "1", "TERM": "xterm-256color"}, false, Level256},
		{"no color", map[string]string{"NO_COLOR": "1", "COLORTERM": "truecolor"}, true, LevelNone},
		{"truecolor", map[string]string{"COLORTERM": "truecolor"}, true, LevelTrueColor},
		{"256", map[string]string{"TERM": "xterm-256color"}, true, Level256},
		{"dumb", map[string]string{"TERM": "dumb"}, true, LevelNone},
		{"kitty", map[string]string{"TERM": "xterm-kitty"}, true, LevelTrueColor},
		{"windows terminal", map[string]string{"WT_SESSION": "x"}, true, LevelTrueColor},
		{"unknown terminal", map[string]string{}, true, LevelBasic},
	}

	for _, tt := range tests {
		if got := detect(env(tt.vars), tt.tty, LevelNone, false); got != tt.want {
			t.Errorf("%s: detect = %v, want %v", tt.name, got, tt.want)
		}
	}

	// A legacy Windows console without ANSI support gets plain text.
	if got := detect(env(nil), true, LevelNone, true); got != LevelNone {
		t.Errorf("legacy console: detect = %v, want none", got)
	}
}

func TestDetectBuffer(t *testing.T) {
	t.Setenv("FORCE_COLOR", "")
	t.Setenv("CLICOLOR_FORCE", "")

	if got := Detect(&bytes.Buffer{}); got != LevelNone {
		t.Errorf("Detect(buffer) = %v, want none", got)
	}

	SetLevel(Level256)

	defer ResetLevel()

	if got := Detect(&bytes.Buffer{}); got != Level256 {
		t.Errorf("Detect with SetLevel = %v, want 256", got)
	}
}

func TestOSC11(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want Color
		ok   bool
	}{
		{"\x1b]11;rgb:2424/2727/3a3a\x1b\\", RGB(0x24, 0x27, 0x3a), true},
		{"\x1b]11;rgb:ff/ff/ff\x07", RGB(255, 255, 255), true},
		{"\x1b]11;rgba:0000/0000/0000/ffff\x1b\\", RGB(0, 0, 0), true},
		{"\x1b]11;rgb:f/0/0\x07", RGB(255, 0, 0), true},
		{"garbage", None, false},
		{"\x1b]11;rgb:zz/00/00\x07", None, false},
	}

	for _, tt := range tests {
		got, ok := parseOSC11(tt.in)
		if ok != tt.ok || got != tt.want {
			t.Errorf("parseOSC11(%q) = %v, %v; want %v, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}

	if c, ok := parseColorFgBg("15;0"); !ok || c != Black || !c.IsDark() {
		t.Errorf("COLORFGBG 15;0 = %v, %v; want dark black", c, ok)
	}

	if c, ok := parseColorFgBg("0;default;15"); !ok || c.IsDark() {
		t.Errorf("COLORFGBG 0;default;15 = %v, %v; want light", c, ok)
	}
}

func TestMixAndLuminance(t *testing.T) {
	t.Parallel()

	black, white := RGB(0, 0, 0), RGB(255, 255, 255)
	if got := black.Mix(white, 0.5); got != RGB(128, 128, 128) {
		t.Errorf("Mix = %v", got)
	}

	if !RGB(0x24, 0x27, 0x3a).IsDark() || RGB(0xef, 0xf1, 0xf5).IsDark() {
		t.Error("IsDark misjudges Macchiato base or Latte base")
	}

	if r, g, b := Index(196).RGB(); r != 255 || g != 0 || b != 0 {
		t.Errorf("Index(196).RGB = %d,%d,%d", r, g, b)
	}

	if r, g, b := Index(232).RGB(); r != 8 || g != 8 || b != 8 {
		t.Errorf("Index(232).RGB = %d,%d,%d", r, g, b)
	}
}

func TestDownsampleRoundTrip(t *testing.T) {
	t.Parallel()

	for i := range 256 {
		r, g, b := index256ToRGB(uint8(i))
		if got := rgbTo256(r, g, b); i >= 16 && index256ToRGB2(got) != [3]uint8{r, g, b} {
			t.Errorf("rgbTo256(index %d) = %d", i, got)
		}
	}
}

func index256ToRGB2(n uint8) [3]uint8 {
	r, g, b := index256ToRGB(n)
	return [3]uint8{r, g, b}
}
