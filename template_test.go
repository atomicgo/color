package template

import "testing"

func TestHelloWorld(t *testing.T) {
	t.Parallel()

	if HelloWorld() != "Hello, World!" {
		t.Fatal("Not equal")
	}
}
