package parser

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	input := "echo Hello World"
	expected := []string{"echo", "Hello", "World"}

	result := ParseCommand(input)

	if len(result) != len(expected) {
		t.Errorf("Ожидалось %v, получено %v", expected, result)
	}
}
