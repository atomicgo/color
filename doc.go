/*
Package color colors terminal output. It handles ANSI, 256-color and true
color output, finds out what the terminal can show, and falls back to fewer
colors or plain text when it has to.

It is only about color and text attributes (bold, italic, underline...). It
has no predefined styles such as "error" or "success"; build those in your
own code from a Style.

# Colors

A Color is one of the 16 basic ANSI colors (Red, BrightBlue...), an entry of
the 256-color palette (Index), or a 24-bit value (RGB, Hex, Parse):

	color.Red
	color.Index(208)
	color.RGB(183, 189, 248)
	c, err := color.Hex("#b7bdf8")
	c, err := color.Parse("brightblue") // also "#fff", "208", "none"

# Styles

A Style combines a foreground, a background and attributes. Styles are
values: every method returns a copy, so they are safe to share between
goroutines.

	title := color.New(color.Bold).Fg(color.MustHex("#b7bdf8"))
	fmt.Println(title.Sprint("Hello"))
	fmt.Println(color.Green.Sprint("done"))

Sprint and Print style for standard output, Fprint for any writer. They pick
the color level with Detect, so output piped into a file or another program
stays plain. To control the level yourself, use Render:

	s := title.Render(color.LevelTrueColor, "Hello")

# Detection

Detect reports the level a writer can show: LevelNone, LevelBasic (16
colors), Level256 or LevelTrueColor. It honors FORCE_COLOR, NO_COLOR,
CLICOLOR and CLICOLOR_FORCE, checks whether the writer is a terminal, and
reads COLORTERM, TERM and TERM_PROGRAM. On Windows it turns on ANSI
processing in the console. SetLevel overrides it, for example for a --color
flag.

Colors above the detected level are converted to the closest color the
terminal can show. Basic colors are never converted, so they follow the
user's terminal palette.

Background asks the terminal for its background color, so a program can pick
a light or a dark palette.
*/
package color
