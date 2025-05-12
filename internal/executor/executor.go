package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/deker104/cli/internal/env"
	"github.com/spf13/pflag"
)

// ExitCode — специальный код выхода, используемый командой `exit`.
const ExitCode = 127 // выбрано как стандартный "ошибочный" код завершения

// Executor — структура, выполняющая команды CLI.
type Executor struct {
	env *env.EnvManager
}

// NewExecutor создает Executor с поддержкой переменных окружения.
func NewExecutor(env *env.EnvManager) *Executor {
	return &Executor{env: env}
}

// Execute выполняет список команд в пайпе. Возвращает код завершения.
func (e *Executor) Execute(pipeline [][]string) int {
	if len(pipeline) == 1 {
		return e.runSingleCommand(pipeline[0])
	}
	return e.runPipeline(pipeline)
}

// runSingleCommand выполняет одну команду (встроенную или внешнюю).
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
	case "grep":
		return e.runGrep(tokens[1:])
	default:
		return e.runExternalCommand(tokens)
	}
	return 0
}

// runPipeline — выполняет пайп (несколько команд, соединённых через `|`).
func (e *Executor) runPipeline(pipeline [][]string) int {
	var commands []*exec.Cmd
	var pipes []*os.File

	// Создание exec.Cmd объектов для всех команд
	for _, cmdArgs := range pipeline {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Env = os.Environ()
		commands = append(commands, cmd)
	}

	// Связывание stdout одной команды с stdin следующей
	for i := 0; i < len(commands)-1; i++ {
		readEnd, writeEnd, err := os.Pipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "pipe error: %v\n", err)
			return 1
		}
		commands[i].Stdout = writeEnd
		commands[i+1].Stdin = readEnd
		pipes = append(pipes, readEnd, writeEnd)
	}

	// Последняя команда пишет в стандартный вывод
	commands[len(commands)-1].Stdout = os.Stdout
	commands[len(commands)-1].Stderr = os.Stderr

	// Запуск команд
	for _, cmd := range commands {
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "command start error: %v\n", err)
			return 1
		}
	}

	// Закрытие всех пайпов
	for _, pipe := range pipes {
		pipe.Close()
	}

	// Ожидание завершения всех команд
	finalCode := 0
	for i, cmd := range commands {
		err := cmd.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				code := exitErr.ExitCode()
				fmt.Fprintf(os.Stderr, "command %d exited with code %d: %v\n", i, code, err)
				finalCode = code
			} else {
				fmt.Fprintf(os.Stderr, "command %d failed: %v\n", i, err)
				finalCode = 1
			}
		} else if finalCode == 0 {
			finalCode = cmd.ProcessState.ExitCode()
		}
	}

	return finalCode
}

// runCat реализует встроенную команду `cat`.
func (e *Executor) runCat(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "cat: missing file operand")
		return 1
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cat: %v\n", err)
		return 1
	}
	fmt.Print(string(data))
	return 0
}

// runWc реализует команду `wc` — подсчёт строк, слов и байт.
func (e *Executor) runWc(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "wc: missing file operand")
		return 1
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "wc: %v\n", err)
		return 1
	}
	lines := strings.Count(string(data), "\n")
	words := len(strings.Fields(string(data)))
	bytes := len(data)
	fmt.Printf("%d %d %d %s\n", lines, words, bytes, args[0])
	return 0
}

// runExternalCommand выполняет внешнюю системную команду.
func (e *Executor) runExternalCommand(tokens []string) int {
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

// GrepOptions — параметры команды grep.
type GrepOptions struct {
	Pattern    string
	Filename   string
	WordMatch  bool
	IgnoreCase bool
	AfterLines int
}

// runGrep парсит аргументы и вызывает executeGrep.
func (e *Executor) runGrep(args []string) int {
	flags := pflag.NewFlagSet("grep", pflag.ContinueOnError)
	wordMatch := flags.BoolP("word", "w", false, "Искать только целые слова")
	ignoreCase := flags.BoolP("ignore-case", "i", false, "Игнорировать регистр")
	afterLines := flags.IntP("after", "A", 0, "Количество строк после совпадения")

	err := flags.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "grep: error parsing flags")
		return 1
	}

	remainingArgs := flags.Args()
	if len(remainingArgs) < 2 {
		fmt.Fprintln(os.Stderr, "grep: usage: grep [options] pattern file")
		return 1
	}

	options := GrepOptions{
		Pattern:    remainingArgs[0],
		Filename:   remainingArgs[1],
		WordMatch:  *wordMatch,
		IgnoreCase: *ignoreCase,
		AfterLines: *afterLines,
	}
	return e.executeGrep(options)
}

// executeGrep выполняет саму логику поиска по регулярному выражению.
func (e *Executor) executeGrep(opts GrepOptions) int {
	file, err := os.Open(opts.Filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grep: %v\n", err)
		return 1
	}
	defer file.Close()

	pattern := opts.Pattern
	if opts.WordMatch {
		pattern = `\b` + pattern + `\b`
	}
	if opts.IgnoreCase {
		pattern = `(?i)` + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grep: invalid regex pattern: %v\n", err)
		return 1
	}

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	linesAfter := 0

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			linesAfter = max(linesAfter, opts.AfterLines)
			buffer.WriteString(line + "\n")
		} else if linesAfter > 0 {
			buffer.WriteString(line + "\n")
			linesAfter--
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "grep: error reading file: %v\n", err)
		return 1
	}

	fmt.Print(buffer.String())
	return 0
}

// max возвращает максимум из двух целых чисел.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
