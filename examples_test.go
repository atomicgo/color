package template_test

import (
	"fmt"

	"atomicgo.dev/template"
)

func Example_demo() {
	fmt.Println(template.HelloWorld())
	// Output: Hello, World!
}

func ExampleHelloWorld() {
	fmt.Println(template.HelloWorld())
	// Output: Hello, World!
}
