package env

import "testing"

func TestEnvManager(t *testing.T) {
	env := NewEnvManager()

	env.Set("FOO", "bar")
	if env.Get("FOO") != "bar" {
		t.Errorf("Expected 'bar', got '%s'", env.Get("FOO"))
	}
}

func TestEnvManagerUnsetVariable(t *testing.T) {
	env := NewEnvManager()
	value := env.Get("UNSET_VAR")
	if value != "" {
		t.Errorf("Expected empty string, got '%s'", value)
	}
}

func TestEnvManagerOverrideVariable(t *testing.T) {
	env := NewEnvManager()
	env.Set("FOO", "bar")
	env.Set("FOO", "baz")

	if env.Get("FOO") != "baz" {
		t.Errorf("Expected 'baz', got '%s'", env.Get("FOO"))
	}
}

func TestEnvManagerSystemVariable(t *testing.T) {
	env := NewEnvManager()
	value := env.Get("PATH")
	if value == "" {
		t.Errorf("Expected non-empty PATH variable, got empty string")
	}
}

func TestEnvManagerMultipleVariables(t *testing.T) {
	env := NewEnvManager()
	env.Set("A", "1")
	env.Set("B", "2")

	if env.Get("A") != "1" || env.Get("B") != "2" {
		t.Errorf("Expected A=1 and B=2, got A=%s B=%s", env.Get("A"), env.Get("B"))
	}
}
