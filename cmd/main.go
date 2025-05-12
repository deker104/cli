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

// main — точка входа в CLI-интерпретатор.
// Читает строки из stdin, разбивает на команды, подставляет переменные, запускает исполнение.
func main() {
	envManager := env.NewEnvManager()        // Менеджер переменных окружения
	exec := executor.NewExecutor(envManager) // Исполнитель команд

	fmt.Println("Simple CLI. Type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)

	var lastExitCode int // Код возврата последней команды

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		// Разбор строки на команды и аргументы
		tokens := parser.Parse(line)
		if len(tokens) == 0 {
			continue
		}

		// Подстановка переменных окружения и $?
		var substituted [][]string
		for _, cmd := range tokens {
			substituted = append(substituted, substitutor.Substitute(cmd, lastExitCode))
		}

		// Выполнение команды или пайпа
		lastExitCode = exec.Execute(substituted)

		// Завершение при exit
		if lastExitCode == executor.ExitCode {
			break
		}
	}

	// Завершаем CLI с последним кодом ошибки (если был)
	if lastExitCode != 0 && lastExitCode != executor.ExitCode {
		os.Exit(lastExitCode)
	}
}
