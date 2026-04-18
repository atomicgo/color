package color

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestANSIColor_Sequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		color      ANSIColor
		background bool
		want       string
	}{
		{"black foreground", ANSIBlack, false, "\033[38;5;0m"},
		{"red foreground", ANSIRed, false, "\033[38;5;1m"},
		{"white foreground", ANSIWhite, false, "\033[38;5;7m"},
		{"bright black foreground", ANSIBrightBlack, false, "\033[38;5;8m"},
		{"bright white foreground", ANSIBrightWhite, false, "\033[38;5;15m"},
		{"black background", ANSIBlack, true, "\033[48;5;0m"},
		{"red background", ANSIRed, true, "\033[48;5;1m"},
		{"bright white background", ANSIBrightWhite, true, "\033[48;5;15m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.color.Sequence(tc.background); got != tc.want {
				t.Errorf("Sequence(%v) = %q, want %q", tc.background, got, tc.want)
			}
		})
	}
}

func TestANSIColor_ImplementsColor(t *testing.T) {
	t.Parallel()

	var _ Color = ANSIRed
}

func TestANSI256Color_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		color ANSI256Color
		want  string
	}{
		{"first color", ANSI256Color(0), "#000000"},
		{"red", ANSI256Color(9), "#ff0000"},
		{"bright red in 256", ANSI256Color(196), "#ff0000"},
		{"last grayscale", ANSI256Color(255), "#eeeeee"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.color.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestANSI256Color_Sequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		color      ANSI256Color
		background bool
		want       string
	}{
		{"red foreground", ANSI256Color(9), false, "\033[38;2;255;0;0m"},
		{"red background", ANSI256Color(9), true, "\033[48;2;255;0;0m"},
		{"black foreground", ANSI256Color(0), false, "\033[38;2;0;0;0m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.color.Sequence(tc.background); got != tc.want {
				t.Errorf("Sequence(%v) = %q, want %q", tc.background, got, tc.want)
			}
		})
	}
}

func TestANSI256Color_ImplementsColor(t *testing.T) {
	t.Parallel()

	var _ Color = ANSI256Color(0)
}

func TestANSI256Color_TableCoverage(t *testing.T) {
	t.Parallel()

	// Every entry should be a valid 7-char hex string. This guards against
	// accidental corruption of the lookup table.
	for i, hex := range ansi256Table {
		if len(hex) != 7 || hex[0] != '#' {
			t.Errorf("ansi256Table[%d] = %q, want 7-char \"#RRGGBB\"", i, hex)
		}
	}
}

func TestRGBColor_Hex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		color RGBColor
		want  string
	}{
		{"black", RGBColor{0, 0, 0}, "#000000"},
		{"white", RGBColor{255, 255, 255}, "#FFFFFF"},
		{"red", RGBColor{255, 0, 0}, "#FF0000"},
		{"green", RGBColor{0, 255, 0}, "#00FF00"},
		{"blue", RGBColor{0, 0, 255}, "#0000FF"},
		{"mid gray", RGBColor{128, 128, 128}, "#808080"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.color.Hex(); got != tc.want {
				t.Errorf("Hex() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRGBColor_Sequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		color      RGBColor
		background bool
		want       string
	}{
		{"red foreground", RGBColor{255, 0, 0}, false, "\033[38;2;255;0;0m"},
		{"red background", RGBColor{255, 0, 0}, true, "\033[48;2;255;0;0m"},
		{"black foreground", RGBColor{0, 0, 0}, false, "\033[38;2;0;0;0m"},
		{"white background", RGBColor{255, 255, 255}, true, "\033[48;2;255;255;255m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.color.Sequence(tc.background); got != tc.want {
				t.Errorf("Sequence(%v) = %q, want %q", tc.background, got, tc.want)
			}
		})
	}
}

func TestRGBColor_ImplementsColor(t *testing.T) {
	t.Parallel()

	var _ Color = RGBColor{}
}

func TestNewColorFromRGB(t *testing.T) {
	t.Parallel()

	c := NewColorFromRGB(10, 20, 30)

	rgb, ok := c.(RGBColor)
	if !ok {
		t.Fatalf("NewColorFromRGB returned %T, want RGBColor", c)
	}

	if rgb.R != 10 || rgb.G != 20 || rgb.B != 30 {
		t.Errorf("NewColorFromRGB = {%d, %d, %d}, want {10, 20, 30}", rgb.R, rgb.G, rgb.B)
	}
}

func TestNewColorFromHex_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		hex  string
		want RGBColor
	}{
		{"with hash", "#ff0000", RGBColor{255, 0, 0}},
		{"without hash", "00ff00", RGBColor{0, 255, 0}},
		{"mixed case", "#AaBbCc", RGBColor{0xAA, 0xBB, 0xCC}},
		{"black", "#000000", RGBColor{0, 0, 0}},
		{"white", "#ffffff", RGBColor{255, 255, 255}},
		{"shorthand 3-digit", "#f0a", RGBColor{0xFF, 0x00, 0xAA}},
		{"shorthand without hash", "abc", RGBColor{0xAA, 0xBB, 0xCC}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := NewColorFromHex(tc.hex)

			rgb, ok := c.(RGBColor)
			if !ok {
				t.Fatalf("NewColorFromHex(%q) returned %T, want RGBColor", tc.hex, c)
			}

			if rgb != tc.want {
				t.Errorf("NewColorFromHex(%q) = %+v, want %+v", tc.hex, rgb, tc.want)
			}
		})
	}
}

func TestNewColorFromHex_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		hex  string
	}{
		{"empty string", ""},
		{"only hash", "#"},
		{"too short", "#12"},
		{"4 chars", "#1234"},
		{"5 chars", "#12345"},
		{"too long", "#1234567"},
		{"non-hex chars", "#gggggg"},
		{"non-hex in green", "#00zz00"},
		{"non-hex in blue", "#0000zz"},
		{"non-hex 3-digit", "#xyz"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := NewColorFromHex(tc.hex); got != NoColor {
				t.Errorf("NewColorFromHex(%q) = %v, want NoColor", tc.hex, got)
			}
		})
	}
}

func TestNewColorFromHex_RoundTrip(t *testing.T) {
	t.Parallel()

	original := RGBColor{R: 0x12, G: 0x34, B: 0x56}

	c := NewColorFromHex(original.Hex())

	rgb, ok := c.(RGBColor)
	if !ok {
		t.Fatalf("round-trip returned %T, want RGBColor", c)
	}

	if rgb != original {
		t.Errorf("round-trip = %+v, want %+v", rgb, original)
	}
}

func TestNoColor_Sequence(t *testing.T) {
	t.Parallel()

	if got := NoColor.Sequence(false); got != "" {
		t.Errorf("NoColor.Sequence(false) = %q, want empty", got)
	}

	if got := NoColor.Sequence(true); got != "" {
		t.Errorf("NoColor.Sequence(true) = %q, want empty", got)
	}
}

func TestModifier_Sequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		modifier Modifier
		want     string
	}{
		{"reset", Reset, "\033[0m"},
		{"bold", Bold, "\033[1m"},
		{"faint", Faint, "\033[2m"},
		{"italic", Italic, "\033[3m"},
		{"underline", Underline, "\033[4m"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.modifier.Sequence(); got != tc.want {
				t.Errorf("Sequence() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMode_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mode Mode
		want string
	}{
		{"TrueColor", TrueColor, "TrueColor"},
		{"ANSI256", ANSI256, "ANSI256"},
		{"ANSI", ANSI, "ANSI"},
		{"Disabled", Disabled, "Disabled"},
		{"unknown", Mode(99), "Unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.mode.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNewStyle_NilColorsBecomeNoColor(t *testing.T) {
	t.Parallel()

	s := NewStyle(nil, nil)

	if s.Foreground != NoColor {
		t.Errorf("Foreground = %v, want NoColor", s.Foreground)
	}

	if s.Background != NoColor {
		t.Errorf("Background = %v, want NoColor", s.Background)
	}

	if len(s.Modifiers) != 0 {
		t.Errorf("Modifiers = %v, want empty", s.Modifiers)
	}
}

func TestNewStyle_WithColorsAndModifiers(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, ANSIBlue, Bold, Underline)

	if s.Foreground != ANSIRed {
		t.Errorf("Foreground = %v, want ANSIRed", s.Foreground)
	}

	if s.Background != ANSIBlue {
		t.Errorf("Background = %v, want ANSIBlue", s.Background)
	}

	if len(s.Modifiers) != 2 || s.Modifiers[0] != Bold || s.Modifiers[1] != Underline {
		t.Errorf("Modifiers = %v, want [Bold, Underline]", s.Modifiers)
	}
}

func TestStyle_AddModifier_Deduplicates(t *testing.T) {
	t.Parallel()

	s := NewStyle(nil, nil)
	s.AddModifier(Bold)
	s.AddModifier(Bold)
	s.AddModifier(Italic)
	s.AddModifier(Italic)

	if len(s.Modifiers) != 2 {
		t.Errorf("Modifiers = %v, want [Bold, Italic]", s.Modifiers)
	}
}

func TestStyle_WithModifier(t *testing.T) {
	t.Parallel()

	base := NewStyle(nil, nil)
	extended := base.WithModifier(Bold)

	if len(base.Modifiers) != 0 {
		t.Errorf("WithModifier mutated receiver: %v", base.Modifiers)
	}

	if len(extended.Modifiers) != 1 || extended.Modifiers[0] != Bold {
		t.Errorf("extended.Modifiers = %v, want [Bold]", extended.Modifiers)
	}
}

func TestStyle_Sequence(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, ANSIBlue, Bold)

	want := ANSIRed.Sequence(false) + ANSIBlue.Sequence(true) + Bold.Sequence()
	if got := s.Sequence(); got != want {
		t.Errorf("Sequence() = %q, want %q", got, want)
	}
}

func TestStyle_Sequence_NoColors(t *testing.T) {
	t.Parallel()

	s := NewStyle(nil, nil, Bold)

	// NoColor contributes no sequence, so only Bold is emitted.
	if got, want := s.Sequence(), Bold.Sequence(); got != want {
		t.Errorf("Sequence() = %q, want %q", got, want)
	}
}

func TestStyle_Sprint(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, nil)

	got := s.Sprint("hello")

	if !strings.HasPrefix(got, s.Sequence()) {
		t.Errorf("Sprint result %q missing style prefix", got)
	}

	if !strings.HasSuffix(got, Reset.Sequence()) {
		t.Errorf("Sprint result %q missing reset suffix", got)
	}

	if !strings.Contains(got, "hello") {
		t.Errorf("Sprint result %q missing payload", got)
	}
}

func TestStyle_Sprintf(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, nil)

	got := s.Sprintf("hello %s %d", "world", 42)

	if !strings.Contains(got, "hello world 42") {
		t.Errorf("Sprintf result %q missing formatted payload", got)
	}
}

func TestStyle_Fprint(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, nil)

	var buf bytes.Buffer

	n, err := s.Fprint(&buf, "hello")
	if err != nil {
		t.Fatalf("Fprint returned error: %v", err)
	}

	if n != buf.Len() {
		t.Errorf("Fprint returned n=%d, buffer has %d", n, buf.Len())
	}

	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("Fprint wrote %q, missing payload", buf.String())
	}
}

func TestStyle_Fprintf(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, nil)

	var buf bytes.Buffer

	if _, err := s.Fprintf(&buf, "hello %s", "world"); err != nil {
		t.Fatalf("Fprintf returned error: %v", err)
	}

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Fprintf wrote %q, missing payload", buf.String())
	}
}

func TestStyle_Fprintln(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, nil)

	var buf bytes.Buffer

	if _, err := s.Fprintln(&buf, "hello"); err != nil {
		t.Fatalf("Fprintln returned error: %v", err)
	}

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("Fprintln result %q missing trailing newline", buf.String())
	}
}

func TestStyle_Fprintfln(t *testing.T) {
	t.Parallel()

	s := NewStyle(ANSIRed, nil)

	var buf bytes.Buffer

	if _, err := s.Fprintfln(&buf, "hello %s", "world"); err != nil {
		t.Fatalf("Fprintfln returned error: %v", err)
	}

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Fprintfln wrote %q, missing payload", buf.String())
	}

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("Fprintfln result %q missing trailing newline", buf.String())
	}
}

func TestStyle_Print_WritesToWriter(t *testing.T) {
	// Not parallel: mutates the package-level Writer.
	var buf bytes.Buffer

	oldWriter := Writer
	Writer = &buf

	t.Cleanup(func() { Writer = oldWriter })

	NewStyle(ANSIRed, nil).Print("hello")

	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("Print wrote %q, missing payload", buf.String())
	}
}

func TestStyle_Printf_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer

	oldWriter := Writer
	Writer = &buf

	t.Cleanup(func() { Writer = oldWriter })

	NewStyle(ANSIRed, nil).Printf("hello %s", "world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Printf wrote %q, missing payload", buf.String())
	}
}

func TestStyle_Println_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer

	oldWriter := Writer
	Writer = &buf

	t.Cleanup(func() { Writer = oldWriter })

	NewStyle(ANSIRed, nil).Println("hello")

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("Println output %q missing trailing newline", buf.String())
	}
}

func TestStyle_Printfln_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer

	oldWriter := Writer
	Writer = &buf

	t.Cleanup(func() { Writer = oldWriter })

	NewStyle(ANSIRed, nil).Printfln("hello %s", "world")

	out := buf.String()
	if !strings.Contains(out, "hello world") || !strings.HasSuffix(out, "\n") {
		t.Errorf("Printfln wrote %q, want payload + newline", out)
	}
}

func TestPresetStyleFuncs(t *testing.T) {
	t.Parallel()

	// Each preset wraps the payload with its own escape sequence and a Reset.
	tests := []struct {
		name string
		fn   func(a ...any) string
	}{
		{"Black", Black},
		{"BrightBlack", BrightBlack},
		{"Red", Red},
		{"BrightRed", BrightRed},
		{"Green", Green},
		{"BrightGreen", BrightGreen},
		{"Yellow", Yellow},
		{"BrightYellow", BrightYellow},
		{"Blue", Blue},
		{"BrightBlue", BrightBlue},
		{"Magenta", Magenta},
		{"BrightMagenta", BrightMagenta},
		{"Cyan", Cyan},
		{"BrightCyan", BrightCyan},
		{"White", White},
		{"BrightWhite", BrightWhite},
		{"Success", Success},
		{"Info", Info},
		{"Warning", Warning},
		{"Error", Error},
		{"Fatal", Fatal},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := tc.fn("hi")
			if !strings.Contains(out, "hi") {
				t.Errorf("%s(%q) = %q, missing payload", tc.name, "hi", out)
			}

			if !strings.HasSuffix(out, Reset.Sequence()) {
				t.Errorf("%s output %q missing Reset suffix", tc.name, out)
			}
		})
	}
}

func TestWriter_DefaultIsStdout(t *testing.T) {
	t.Parallel()

	// The package exposes Writer as a mutable var; it should default to a
	// non-nil sink so zero-configuration use of Print* does not panic.
	if Writer == nil {
		t.Fatal("Writer is nil on package init")
	}
}

// Sanity check that the ANSI256 -> RGB -> Sequence pipeline produces a
// truecolor escape that matches a direct construction from the same hex.
func TestANSI256_MatchesHex(t *testing.T) {
	t.Parallel()

	for i, hex := range ansi256Table {
		fromTable := ANSI256Color(i).Sequence(false)
		fromHex := NewColorFromHex(hex).Sequence(false)

		if fromTable != fromHex {
			t.Fatalf("index %d: ANSI256 sequence %q != hex sequence %q", i, fromTable, fromHex)
		}
	}
}

// Double-check the escape character encoding used by all color sequences.
func TestSequencesUseEscChar(t *testing.T) {
	t.Parallel()

	seqs := []string{
		ANSIRed.Sequence(false),
		ANSI256Color(9).Sequence(false),
		RGBColor{1, 2, 3}.Sequence(false),
		Bold.Sequence(),
	}

	for _, s := range seqs {
		if !strings.HasPrefix(s, "\x1b[") {
			t.Errorf("%q does not start with ESC[", s)
		}

		if !strings.HasSuffix(s, "m") {
			t.Errorf("%q does not end with 'm'", s)
		}
	}
}

// Ensure RGBColor.Hex is compatible with fmt.Sprintf's %X verb format.
func TestRGBColor_Hex_MatchesSprintf(t *testing.T) {
	t.Parallel()

	c := RGBColor{R: 1, G: 2, B: 3}

	want := fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	if got := c.Hex(); got != want {
		t.Errorf("Hex() = %q, want %q", got, want)
	}
}
