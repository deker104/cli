package executor

import (
	"bytes"
	"io"
	"os"
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

func TestGrep(t *testing.T) {
	exec := NewExecutor(env.NewEnvManager())

	testFile := "testfile.txt"
	content := "hello world\nthis is a test\ngrep is cool\ntest again\nanother test line\ntasty line\n"
	os.WriteFile(testFile, []byte(content), 0644)
	defer os.Remove(testFile)

	output := captureOutput(func() {
		exec.Execute([][]string{{"grep", "test", testFile}})
	})
	expected := "this is a test\ntest again\nanother test line\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	output = captureOutput(func() {
		exec.Execute([][]string{{"grep", "-w", "test", testFile}})
	})
	expected = "this is a test\ntest again\nanother test line\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	output = captureOutput(func() {
		exec.Execute([][]string{{"grep", "-i", "TEST", testFile}})
	})
	expected = "this is a test\ntest again\nanother test line\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	output = captureOutput(func() {
		exec.Execute([][]string{{"grep", "-A", "1", "grep", testFile}})
	})
	expected = "grep is cool\ntest again\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}

	output = captureOutput(func() {
		exec.Execute([][]string{{"grep", "hello", "nofile.txt"}})
	})
	if output == "" {
		t.Errorf("Expected error message for missing file")
	}
}
