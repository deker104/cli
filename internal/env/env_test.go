package env

import "testing"

func TestEnvSetAndGet(t *testing.T) {
	Set("VAR", "3228")
	if Get("VAR") != "3228" {
		t.Errorf("Ожидалось '3228', получено '%s'", Get("VAR"))
	}
}
