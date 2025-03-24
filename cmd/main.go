package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deker104/cli/internal/env"
	"github.com/deker104/cli/internal/executor"
	"github.com/deker104/cli/internal/parser"
)

func main() {
	// Инициализация окружения
	envManager := env.NewEnvManager()
	exec := executor.NewExecutor(envManager)

	fmt.Println("Simple CLI. Type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		// Парсинг команды
		tokens := parser.Parse(line)
		if len(tokens) == 0 {
			continue
		}

		// Выполнение команды
		exitCode := exec.Execute(tokens)
		if exitCode == executor.ExitCode {
			break
		}
	}
}
