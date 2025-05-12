package parser

import (
	"reflect"
	"testing"
)

func makeTokens(args ...Token) []Token {
	return args
}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected [][]Token
	}{
		{"echo hello", [][]Token{
			makeTokens(Token{"echo", true}, Token{"hello", true}),
		}},
		{"echo 'hello world'", [][]Token{
			makeTokens(Token{"echo", true}, Token{"hello world", false}),
		}},
		{`echo "hello world"`, [][]Token{
			makeTokens(Token{"echo", true}, Token{"hello world", true}),
		}},
		{"ls -la", [][]Token{
			makeTokens(Token{"ls", true}, Token{"-la", true}),
		}},
		{"echo 'a b' c", [][]Token{
			makeTokens(Token{"echo", true}, Token{"a b", false}, Token{"c", true}),
		}},
	}

	for _, tt := range tests {
		got := Parse(tt.input)
		if !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("Parse(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestParsePipes(t *testing.T) {
	input := "ls | grep go"
	expected := [][]Token{
		makeTokens(Token{"ls", true}),
		makeTokens(Token{"grep", true}, Token{"go", true}),
	}

	got := Parse(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestParseMultiplePipes(t *testing.T) {
	input := "cat file.txt | grep foo | wc -l"
	expected := [][]Token{
		makeTokens(Token{"cat", true}, Token{"file.txt", true}),
		makeTokens(Token{"grep", true}, Token{"foo", true}),
		makeTokens(Token{"wc", true}, Token{"-l", true}),
	}

	got := Parse(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestStrongQuotes(t *testing.T) {
	input := "echo '$HOME'"
	expected := [][]Token{
		makeTokens(Token{"echo", true}, Token{"$HOME", false}),
	}

	got := Parse(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestWeakQuotes(t *testing.T) {
	input := `echo "$HOME"`
	expected := [][]Token{
		makeTokens(Token{"echo", true}, Token{"$HOME", true}),
	}

	got := Parse(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
