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
