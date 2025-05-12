package executor

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/deker104/cli/internal/env"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

func TestEcho(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	output := captureOutput(func() {
		exec.Execute([][]string{{"echo", "hello"}})
	})

	expected := "hello\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestPwd(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	output := captureOutput(func() {
		exec.Execute([][]string{{"pwd"}})
	})

	if output == "" {
		t.Errorf("Expected non-empty output for pwd")
	}
}

func TestExit(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	code := exec.Execute([][]string{{"exit"}})
	if code != ExitCode {
		t.Errorf("Expected exit code %d, got %d", ExitCode, code)
	}
}

func TestUnknownCommand(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	output := captureOutput(func() {
		exec.Execute([][]string{{"unknown_command"}})
	})

	if output == "" {
		t.Errorf("Expected error output for unknown command")
	}
}

func TestPipeEchoToWc(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	output := captureOutput(func() {
		exec.Execute([][]string{
			{"echo", "one two three"},
			{"wc"},
		})
	})

	if !(strings.Contains(output, "1") && strings.Contains(output, "3")) {
		t.Errorf("Expected line and word count, got: %q", output)
	}
}

func TestPipeEchoToGrep(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	output := captureOutput(func() {
		exec.Execute([][]string{
			{"echo", "this is a test line"},
			{"grep", "test"},
		})
	})

	if !strings.Contains(output, "test line") {
		t.Errorf("Expected grep to find 'test', got: %q", output)
	}
}

func TestExternalCommandFailureCode(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	// Преднамеренно ошибка
	code := exec.Execute([][]string{{"nonexistent_command"}})
	if code == 0 {
		t.Errorf("Expected non-zero exit code for invalid command, got %d", code)
	}
}
