package parser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected [][]string
	}{
		{"echo hello", [][]string{{"echo", "hello"}}},
		{"echo 'hello world'", [][]string{{"echo", "hello world"}}},
		{"echo \"hello world\"", [][]string{{"echo", "hello world"}}},
		{"ls -la", [][]string{{"ls", "-la"}}},
		{"echo 'a b' c", [][]string{{"echo", "a b", "c"}}},
	}

	for _, test := range tests {
		result := Parse(test.input)
		fmt.Printf("Expected: %#v, Got: %#v\n", test.expected, result) // Debug output
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Parse(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestParsePipes(t *testing.T) {
	result := Parse("ls | grep go")
	expected := [][]string{{"ls"}, {"grep", "go"}}

	fmt.Printf("Expected: %#v, Got: %#v\n", expected, result) // Debugging output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseMultiplePipes(t *testing.T) {
	result := Parse("cat file.txt | grep foo | wc -l")
	expected := [][]string{{"cat", "file.txt"}, {"grep", "foo"}, {"wc", "-l"}}

	fmt.Printf("Expected: %#v, Got: %#v\n", expected, result) // Debugging output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
