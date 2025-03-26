package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/deker104/cli/internal/env"
)

const ExitCode = 127

type Executor struct {
	env *env.EnvManager
}

func NewExecutor(env *env.EnvManager) *Executor {
	return &Executor{env: env}
}

func (e *Executor) Execute(pipeline [][]string) int {
	if len(pipeline) == 1 {
		return e.runSingleCommand(pipeline[0])
	}
	return e.runPipeline(pipeline)
}

func (e *Executor) runSingleCommand(tokens []string) int {
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

func (e *Executor) runPipeline(pipeline [][]string) int {
	var lastExitCode int
	var prevPipe *os.File

	for i, command := range pipeline {
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Env = os.Environ()

		// Если не последняя команда, создаем пайп
		if i != len(pipeline)-1 {
			pipeOut, pipeIn, _ := os.Pipe()
			cmd.Stdout = pipeIn
			prevPipe = pipeOut
		} else {
			cmd.Stdout = os.Stdout
		}

		// Подключаем stdin
		if prevPipe != nil {
			cmd.Stdin = prevPipe
			prevPipe.Close()
		}

		// Запускаем команду
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			lastExitCode = 1
		} else {
			lastExitCode = cmd.ProcessState.ExitCode()
		}
	}

	return lastExitCode
}

// runCat — встроенная команда `cat`
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

// runWc — встроенная команда `wc`
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

// runExternalCommand — выполняет внешнюю программу
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
