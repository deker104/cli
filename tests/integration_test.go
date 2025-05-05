package tests

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const inputFile = "../tests/test_input.txt"

// runCLI запускает CLI с заданной строкой и возвращает stdout
func runCLI(input string) (string, error) {
	cmd := exec.Command("go", "run", "../cmd/main.go")
	cmd.Stderr = os.Stderr
	stdin := strings.NewReader(input + "\nexit\n")
	cmd.Stdin = stdin

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	return out.String(), err
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
	out, err := runCLI("cat " + inputFile)
	if err != nil {
		t.Fatalf("cat command failed to run: %v", err)
	}
	clean := stripPrompt(out)
	if !strings.Contains(clean, "hello world") {
		t.Errorf("TestCatCommand failed.\nExpected output to contain: \"hello world\"\nActual:\n%s", clean)
	}
}

func TestEchoCommand(t *testing.T) {
	out, _ := runCLI("echo hello")
	clean := stripPrompt(out)
	expected := "hello\n\n"
	if clean != expected {
		t.Errorf("TestEchoCommand failed.\nExpected: %q\nGot: %q", expected, clean)
	}
}

func TestWcCommand(t *testing.T) {
	out, _ := runCLI("wc " + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, inputFile) {
		t.Errorf("TestWcCommand failed.\nExpected output to mention filename: %q\nActual: %s", inputFile, clean)
	}
}

func TestPwdCommand(t *testing.T) {
	out, _ := runCLI("pwd")
	clean := stripPrompt(out)
	if !strings.Contains(clean, "/") {
		t.Errorf("TestPwdCommand failed.\nExpected output to be a path.\nGot: %q", clean)
	}
}

func TestUnknownCommandExternal(t *testing.T) {
	os.Setenv("HOME", "/tmp/fakehome")
	out, _ := runCLI("echo $HOME")
	clean := stripPrompt(out)
	expected := "/tmp/fakehome\n\n"
	if clean != expected {
		t.Errorf("TestUnknownCommandExternal failed.\nExpected: %q\nGot: %q", expected, clean)
	}
}

func TestQuotedArguments(t *testing.T) {
	out, _ := runCLI(`echo "a b c"`)
	clean := stripPrompt(out)
	expected := "a b c\n\n"
	if clean != expected {
		t.Errorf("TestQuotedArguments failed.\nExpected: %q\nGot: %q", expected, clean)
	}
}

func TestGrepBasic(t *testing.T) {
	out, _ := runCLI(`grep test ` + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "this is a test") {
		t.Errorf("TestGrepBasic failed.\nExpected to find: \"this is a test\"\nGot: %s", clean)
	}
}

func TestGrepWordMatch(t *testing.T) {
	out, _ := runCLI(`grep -w test ` + inputFile)
	clean := stripPrompt(out)
	if strings.Contains(clean, "testing") || strings.Contains(clean, "contest") {
		t.Errorf("TestGrepWordMatch failed.\nUnexpected lines matched: %s", clean)
	}
}

func TestGrepIgnoreCase(t *testing.T) {
	out, _ := runCLI(`grep -i FOO ` + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "FOO bar") || !strings.Contains(clean, "foobar") {
		t.Errorf("TestGrepIgnoreCase failed.\nExpected matches not found.\nOutput: %s", clean)
	}
}

func TestGrepAfterLines(t *testing.T) {
	out, _ := runCLI(`grep -A 2 match ` + inputFile)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "line2") {
		t.Errorf("TestGrepAfterLines failed.\nExpected line2 after match not found.\nOutput:\n%s", clean)
	}
}

func TestPipelineCommands(t *testing.T) {
	out, _ := runCLI(`cat ` + inputFile + ` | grep test | wc`)
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
	out, _ := runCLI(`echo $FOO_TEST_VAR`)
	clean := stripPrompt(out)
	if !strings.Contains(clean, "VALUE_123") {
		t.Errorf("TestEnvSubstitution failed.\nExpected: \"VALUE_123\"\nGot: %q", clean)
	}
}
