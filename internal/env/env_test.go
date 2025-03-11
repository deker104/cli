package env

import "testing"

func TestEnvSetAndGet(t *testing.T) {
	Set("VAR", "42")
	if Get("VAR") != "42" {
		t.Errorf("Ожидалось '42', получено '%s'", Get("VAR"))
	}
}