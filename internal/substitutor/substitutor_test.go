package substitutor

import (
	"os"
	"testing"

	"github.com/deker104/cli/internal/parser"
)

func makeTokens(strs []string, subst bool) []parser.Token {
	var tokens []parser.Token
	for _, s := range strs {
		tokens = append(tokens, parser.Token{Value: s, SubstituteEnv: subst})
	}
	return tokens
}

func TestSubstituteBasic(t *testing.T) {
	input := makeTokens([]string{"echo", "hello"}, true)
	expected := []string{"echo", "hello"}

	result := Substitute(input, 0)
	compare(t, expected, result)
}

func TestSubstituteEnvVar(t *testing.T) {
	os.Setenv("TEST_VAR", "abc123")
	input := makeTokens([]string{"echo", "$TEST_VAR"}, true)
	expected := []string{"echo", "abc123"}

	result := Substitute(input, 0)
	compare(t, expected, result)
}

func TestSubstituteExitCode(t *testing.T) {
	input := makeTokens([]string{"echo", "$?"}, true)
	expected := []string{"echo", "42"}

	result := Substitute(input, 42)
	compare(t, expected, result)
}

func TestNoSubstitutionForLiterals(t *testing.T) {
	input := []parser.Token{
		{Value: "echo", SubstituteEnv: true},
		{Value: "$HOME", SubstituteEnv: false}, // strong quote
	}
	expected := []string{"echo", "$HOME"}

	result := Substitute(input, 0)
	compare(t, expected, result)
}

func TestMixedTokens(t *testing.T) {
	os.Setenv("MIXED", "YES")
	input := []parser.Token{
		{Value: "echo", SubstituteEnv: true},
		{Value: "$MIXED", SubstituteEnv: true},
		{Value: "$IGNORED", SubstituteEnv: false},
	}
	expected := []string{"echo", "YES", "$IGNORED"}

	result := Substitute(input, 0)
	compare(t, expected, result)
}

func compare(t *testing.T, expected, got []string) {
	if len(expected) != len(got) {
		t.Fatalf("Expected len %d, got %d\nExpected: %v\nGot: %v", len(expected), len(got), expected, got)
	}
	for i := range expected {
		if expected[i] != got[i] {
			t.Errorf("Mismatch at index %d: expected %q, got %q", i, expected[i], got[i])
		}
	}
}
