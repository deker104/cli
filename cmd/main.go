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

	var lastExitCode int

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
			substituted = append(substituted, substitutor.Substitute(cmd, lastExitCode))
		}

		lastExitCode = exec.Execute(substituted)
		if lastExitCode == executor.ExitCode {
			break
		}
	}

	if lastExitCode != 0 && lastExitCode != executor.ExitCode {
		os.Exit(lastExitCode)
	}
}
