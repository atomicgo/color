package color_test

import (
	"fmt"

	"atomicgo.dev/color"
)

func Example_demo() {
	// Examples render at a fixed level so the output is predictable.
	// Real programs call Sprint or Print and let Detect pick the level.
	l := color.LevelBasic

	fmt.Printf("%q\n", color.Green.Fg().Render(l, "done"))
	fmt.Printf("%q\n", color.New(color.Bold).Fg(color.Red).Render(l, "failed"))
	fmt.Printf("%q\n", color.Blue.Bg().Render(l, "info"))
	// Output:
	// "\x1b[32mdone\x1b[0m"
	// "\x1b[1;31mfailed\x1b[0m"
	// "\x1b[44minfo\x1b[0m"
}

func ExampleParse() {
	for _, s := range []string{"#b7bdf8", "brightblue", "208", "none"} {
		c, err := color.Parse(s)
		fmt.Println(c, err)
	}
	// Output:
	// #b7bdf8 <nil>
	// brightblue <nil>
	// 208 <nil>
	// none <nil>
}

func ExampleStyle_Render() {
	lavender := color.MustHex("#b7bdf8")
	title := color.New(color.Bold).Fg(lavender)

	// The same style at each level. Terminals with fewer colors get the
	// closest color they can show.
	for _, l := range []color.Level{color.LevelTrueColor, color.Level256, color.LevelBasic, color.LevelNone} {
		fmt.Printf("%-9s %q\n", l, title.Render(l, "Hi"))
	}
	// Output:
	// truecolor "\x1b[1;38;2;183;189;248mHi\x1b[0m"
	// 256       "\x1b[1;38;5;147mHi\x1b[0m"
	// 16        "\x1b[1;37mHi\x1b[0m"
	// none      "Hi"
}

func ExampleColor_Mix() {
	from, to := color.MustHex("#b7bdf8"), color.MustHex("#8bd5ca")

	// A three-step gradient from lavender to teal.
	for i := range 3 {
		fmt.Println(from.Mix(to, float64(i)/2))
	}
	// Output:
	// #b7bdf8
	// #a1c9e1
	// #8bd5ca
}

func ExampleStrip() {
	s := color.Red.Fg().Render(color.LevelTrueColor, "error")
	fmt.Println(color.Strip(s))
	// Output: error
}
