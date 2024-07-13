package color_test

import (
	"fmt"

	"atomicgo.dev/color"
)

func Example_demo() {
	// Simple coloring
	fmt.Println("Hello, " + color.Green("World") + "!")

	fmt.Println() // blank line

	// Theme colors - can be customized in init() function if needed
	fmt.Println(color.Success("Success message"))
	fmt.Println(color.Info("Info message"))
	fmt.Println(color.Warning("Warning message"))
	fmt.Println(color.Error("Error message"))
	fmt.Println(color.Fatal("Fatal message"))

	fmt.Println() // blank line

	// Supports ANSI colors
	ansiRed := color.NewStyle(color.ANSIRed, nil).Sprint
	fmt.Println(ansiRed("This is printed red using an ANSI color code"))

	// Supports ANSI256 colors
	ansi256Red := color.NewStyle(color.ANSI256Color(196), nil).Sprint
	fmt.Println(ansi256Red("This is printed red using an ANSI256 color code"))

	// Supports RGB colors
	redRGB := color.NewStyle(color.NewColorFromRGB(255, 0, 0), nil).Sprint
	fmt.Println(redRGB("This is printed red using a RGB color code"))

	// Supports hex colors
	redHex := color.NewStyle(color.NewColorFromHex("#ff0000"), nil).Sprint
	fmt.Println(redHex("This is printed red using a hex color code"))
}
