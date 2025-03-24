package parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"echo hello", []string{"echo", "hello"}},
		{"echo 'hello world'", []string{"echo", "hello world"}},
		{"echo \"hello world\"", []string{"echo", "hello world"}},
		{"ls -la", []string{"ls", "-la"}},
		{"echo 'a b' c", []string{"echo", "a b", "c"}},
	}

	for _, test := range tests {
		result := Parse(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Parse(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestParseEmptyInput(t *testing.T) {
	result := Parse("")
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestParseMultipleSpaces(t *testing.T) {
	result := Parse("   echo   hello    world   ")
	expected := []string{"echo", "hello", "world"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseEscapedQuotes(t *testing.T) {
	result := Parse(`echo "hello \"world\"!"`)
	expected := []string{"echo", `hello "world"!`}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseUnmatchedQuotes(t *testing.T) {
	result := Parse(`echo "hello world`)
	expected := []string{"echo", `"hello world`}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
