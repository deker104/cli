package substitutor

import (
	"os"
	"reflect"
	"testing"
)

func TestSubstituteEnvVariable(t *testing.T) {
	os.Setenv("USER", "testuser")
	input := []string{"echo", "$USER"}
	expected := []string{"echo", "testuser"}

	result := Substitute(input, 0)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Substitute(%v) = %v; want %v", input, result, expected)
	}
}

func TestSubstituteExitCode(t *testing.T) {
	input := []string{"echo", "$?"}
	expected := []string{"echo", "42"}

	result := Substitute(input, 42)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Substitute(%v) = %v; want %v", input, result, expected)
	}
}

func TestSubstituteUnknownVariable(t *testing.T) {
	input := []string{"echo", "$UNKNOWN"}
	expected := []string{"echo", ""}

	result := Substitute(input, 0)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Substitute(%v) = %v; want %v", input, result, expected)
	}
}

func TestSubstituteNoChange(t *testing.T) {
	input := []string{"echo", "hello", "world"}
	expected := []string{"echo", "hello", "world"}

	result := Substitute(input, 0)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Substitute(%v) = %v; want %v", input, result, expected)
	}
}

func TestSubstituteMixedVariables(t *testing.T) {
	os.Setenv("HOME", "/home/test")
	input := []string{"echo", "$HOME", "$?"}
	expected := []string{"echo", "/home/test", "0"}

	result := Substitute(input, 0)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Substitute(%v) = %v; want %v", input, result, expected)
	}
}
