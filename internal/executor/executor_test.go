package executor

import (
	"testing"

	"github.com/deker104/cli/internal/env"
)

func TestEcho(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())
	if exec.Execute([]string{"echo", "hello"}) != 0 {
		t.Errorf("echo failed")
	}
}
