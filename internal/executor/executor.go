package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/deker104/cli/internal/env"
)

const ExitCode = 127

// Executor выполняет команды.
type Executor struct {
	env *env.EnvManager
}

// NewExecutor создает новый Executor.
func NewExecutor(env *env.EnvManager) *Executor {
	return &Executor{env: env}
}

// Execute выполняет команду.
func (e *Executor) Execute(tokens []string) int {
	cmd := tokens[0]

	switch cmd {
	case "echo":
		fmt.Println(strings.Join(tokens[1:], " "))
	case "pwd":
		dir, _ := os.Getwd()
		fmt.Println(dir)
	case "exit":
		return ExitCode
	case "cat":
		return e.runCat(tokens[1:])
	case "wc":
		return e.runWc(tokens[1:])
	default:
		return e.runExternalCommand(tokens)
	}
	return 0
}

func (e *Executor) runCat(args []string) int {
	if len(args) == 0 {
		fmt.Println("cat: missing file operand")
		return 1
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("cat: %v\n", err)
		return 1
	}
	fmt.Print(string(data))
	return 0
}

func (e *Executor) runWc(args []string) int {
	if len(args) == 0 {
		fmt.Println("wc: missing file operand")
		return 1
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("wc: %v\n", err)
		return 1
	}
	lines := strings.Count(string(data), "\n")
	words := len(strings.Fields(string(data)))
	bytes := len(data)
	fmt.Printf("%d %d %d %s\n", lines, words, bytes, args[0])
	return 0
}

func (e *Executor) runExternalCommand(tokens []string) int {
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
