package env

import "testing"

func TestEnvManager(t *testing.T) {
	env := NewEnvManager()

	env.Set("FOO", "bar")
	if env.Get("FOO") != "bar" {
		t.Errorf("Expected 'bar', got '%s'", env.Get("FOO"))
	}
}
