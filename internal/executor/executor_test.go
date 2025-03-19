package executor

import "testing"

func TestExecuteCommand(t *testing.T) {
	ExecuteCommand("echo", []string{"Hello"})
}
