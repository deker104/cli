package tests

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const inputFile = "../tests/test_input.txt"

// runCLI запускает CLI с заданной строкой и возвращает stdout, код возврата и ошибку запуска
func runCLI(input string) (string, int, error) {
	cmd := exec.Command("go", "run", "../cmd/main.go")
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(input + "\nexit\n")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return "", 0, err
		}
	}
	return out.String(), exitCode, nil
}

// stripPrompt убирает приглашения и заголовок
func stripPrompt(s string) string {
	lines := strings.Split(s, "\n")
	var clean []string
	for _, l := range lines {
		if strings.HasPrefix(l, "> ") {
			clean = append(clean, l[2:])
		} else if !strings.HasPrefix(l, "Simple CLI.") && l != "" {
			clean = append(clean, l)
		}
	}
	return strings.Join(clean, "\n") + "\n"
}

func TestCatCommand(t *testing.T) {
	out, _, err := runCLI("cat " + inputFile)
	if err != nil {
		t.Fatalf("cat command failed to run: %v", err)
	}
	clean := stripPrompt(out)
	if !strings.Contains(clean, "hello world") {
		t.Errorf("TestCatCommand failed.\nExpected output to contain: \"hello world\"\nActual:\n%s", clean)
	}
}

func TestEchoCommand(t *testing.T) {
	out, _, _ := runCLI("echo hello")
	clean := stripPrompt(out)
	expected := "hello\n\n"
	if clean != expected {
		t.Errorf("TestEchoCommand failed.\nExpected: %q\nGot: %q", expected, clean)
	}
}

func TestWcCommand(t *testing.T) {
	out, _, _ := runCLI("wc " + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, inputFile) {
		t.Errorf("TestWcCommand failed.\nExpected output to mention filename: %q\nActual: %s", inputFile, clean)
	}
}

func TestPwdCommand(t *testing.T) {
	out, _, _ := runCLI("pwd")
	clean := stripPrompt(out)
	if !strings.Contains(clean, "/") && !strings.Contains(clean, "\\") {
		t.Errorf("TestPwdCommand failed.\nExpected output to be a path.\nGot: %q", clean)
	}
}

func TestUnknownCommandExternal(t *testing.T) {
	os.Setenv("HOME", "/tmp/fakehome")
	out, _, _ := runCLI("echo $HOME")
	clean := stripPrompt(out)
	expected := "/tmp/fakehome\n\n"
	if clean != expected {
		t.Errorf("TestUnknownCommandExternal failed.\nExpected: %q\nGot: %q", expected, clean)
	}
}

func TestQuotedArguments(t *testing.T) {
	out, _, _ := runCLI(`echo "a b c"`)
	clean := stripPrompt(out)
	expected := "a b c\n\n"
	if clean != expected {
		t.Errorf("TestQuotedArguments failed.\nExpected: %q\nGot: %q", expected, clean)
	}
}

func TestGrepBasic(t *testing.T) {
	out, _, _ := runCLI(`grep test ` + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "this is a test") {
		t.Errorf("TestGrepBasic failed.\nExpected to find: \"this is a test\"\nGot: %s", clean)
	}
}

func TestGrepWordMatch(t *testing.T) {
	out, _, _ := runCLI(`grep -w test ` + inputFile)
	clean := stripPrompt(out)
	if strings.Contains(clean, "testing") || strings.Contains(clean, "contest") {
		t.Errorf("TestGrepWordMatch failed.\nUnexpected lines matched: %s", clean)
	}
}

func TestGrepIgnoreCase(t *testing.T) {
	out, _, _ := runCLI(`grep -i FOO ` + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "FOO bar") || !strings.Contains(clean, "foobar") {
		t.Errorf("TestGrepIgnoreCase failed.\nExpected matches not found.\nOutput: %s", clean)
	}
}

func TestGrepAfterLines(t *testing.T) {
	out, _, _ := runCLI(`grep -A 2 match ` + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "line2") {
		t.Errorf("TestGrepAfterLines failed.\nExpected line2 after match not found.\nOutput:\n%s", clean)
	}
}

func TestPipelineCommands(t *testing.T) {
	out, _, _ := runCLI(`cat ` + inputFile + ` | grep test | wc`)
	clean := stripPrompt(out)
	if strings.Contains(clean, "broken pipe") || strings.Contains(clean, "exit status") {
		t.Errorf("TestPipelineCommands failed.\nError in execution:\n%s", clean)
	}
	if !strings.Contains(clean, "0") && !strings.Contains(clean, "test_input.txt") {
		t.Errorf("TestPipelineCommands failed.\nExpected some count or filename in output.\nGot:\n%s", clean)
	}
}

func TestEnvSubstitution(t *testing.T) {
	os.Setenv("FOO_TEST_VAR", "VALUE_123")
	out, _, _ := runCLI(`echo $FOO_TEST_VAR`)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "VALUE_123") {
		t.Errorf("TestEnvSubstitution failed.\nExpected: \"VALUE_123\"\nGot: %q", clean)
	}
}

func TestStrongQuotesLiteral(t *testing.T) {
	out, _, _ := runCLI(`echo '$HOME'`)
	clean := stripPrompt(out)
	expected := "$HOME\n\n"
	if clean != expected {
		t.Errorf("Expected strong quotes literal, got: %q", clean)
	}
}

func TestWeakQuotesExpand(t *testing.T) {
	os.Setenv("FOO_TEST_VAR", "ABC123")
	out, _, _ := runCLI(`echo "$FOO_TEST_VAR"`)
	clean := stripPrompt(out)
	expected := "ABC123\n\n"
	if clean != expected {
		t.Errorf("Expected weak quotes expansion, got: %q", clean)
	}
}

func TestExternalCommandErrorCode(t *testing.T) {
	out, exitCode, err := runCLI("nonexistent_command")
	if err != nil && exitCode == 0 {
		t.Errorf("Unexpected error type: %v", err)
	}
	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for unknown command, got 0")
	}
	if !strings.Contains(out, "nonexistent_command") {
		t.Errorf("Expected command name in error output, got: %s", out)
	}
}

