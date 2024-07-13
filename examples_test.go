package color_test

import (
	"fmt"

	"atomicgo.dev/color"
)

func Example_demo() {
	// Simple coloring
	fmt.Println("Hello, " + color.Red("World") + "!")

	// Theme colors - can be customized in init() function if needed
	fmt.Println(color.Success("Success message"))
	fmt.Println(color.Info("Info message"))
	fmt.Println(color.Warning("Warning message"))
	fmt.Println(color.Error("Error message"))
	fmt.Println(color.Fatal("Fatal message"))
}
