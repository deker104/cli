package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deker104/cli/internal/env"
	"github.com/deker104/cli/internal/executor"
	"github.com/deker104/cli/internal/parser"
	"github.com/deker104/cli/internal/substitutor"
)

func main() {
	envManager := env.NewEnvManager()
	exec := executor.NewExecutor(envManager)

	fmt.Println("Simple CLI. Type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)

	exitCode := 0 // <— вот тут инициализируем

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		tokens := parser.Parse(line)
		if len(tokens) == 0 {
			continue
		}

		var substituted [][]string
		for _, cmd := range tokens {
			substituted = append(substituted, substitutor.Substitute(cmd, 0))
		}

		exitCode = exec.Execute(substituted)
		if exitCode == executor.ExitCode {
			break
		}
	}

	os.Exit(exitCode)
}
