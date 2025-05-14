package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/deker104/cli/internal/env"
	"github.com/spf13/pflag"
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
	case "grep":
		return e.runGrep(tokens[1:])
	case "cd":
		return e.runCd(tokens[1:])
	case "ls":
		return e.runLs(tokens[1:])
	default:
		return e.runExternalCommand(tokens)
	}
	return 0
}

func (e *Executor) runPipeline(pipeline [][]string) int {
	var commands []*exec.Cmd
	var pipes []*os.File

	// Создание команд
	for _, cmdArgs := range pipeline {
		commands = append(commands, exec.Command(cmdArgs[0], cmdArgs[1:]...))
	}

	// Установка пайпов между командами
	for i := 0; i < len(commands)-1; i++ {
		readEnd, writeEnd, err := os.Pipe()
		if err != nil {
			fmt.Printf("pipe error: %v\n", err)
			return 1
		}
		commands[i].Stdout = writeEnd
		commands[i+1].Stdin = readEnd
		pipes = append(pipes, readEnd, writeEnd)
	}

	// Последняя команда выводит в stdout
	commands[len(commands)-1].Stdout = os.Stdout

	// Запуск всех команд
	for _, cmd := range commands {
		cmd.Env = os.Environ()
		if err := cmd.Start(); err != nil {
			fmt.Printf("command start error: %v\n", err)
			return 1
		}
	}

	// Закрытие всех пайпов в родителе
	for _, pipe := range pipes {
		pipe.Close()
	}

	// Ожидание завершения всех команд
	var lastExitCode int
	for _, cmd := range commands {
		err := cmd.Wait()
		if err != nil {
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

// runLs — встроенная команда `ls`
func (e *Executor) runLs(args []string) int {
	dir := "."
	if len(args) == 1 {
		dir = args[0]
	}
	data, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("ls: %v\n", err)
		return 1
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name() < data[j].Name()
	})
	for _, elem := range data {
		fmt.Println(elem.Name())
	}
	return 0
}

// runCd — встроенная команда `cd`
func (e *Executor) runCd(args []string) int {
	var dir string
	if len(args) == 0 {
		dir = os.Getenv("HOME")
		if dir == "" {
			dir = os.Getenv("USERPROFILE")
		}
		if dir == "" {
			fmt.Println("cd: HOME/USERPROFILE not set")
			return 1
		}
	} else if len(args) == 1 {
		if args[0] == "-" {
			// Если аргумент "-", переходим в предыдущую директорию
			dir = os.Getenv("OLDPWD")
			if dir == "" {
				cwd, err := os.Getwd()
				if err != nil {
					fmt.Println("cd: cannot determine current directory")
					return 1
				}
				dir = cwd
			}
		} else {
			dir = args[0]
		}
	} else {
		fmt.Println("cd: too many arguments")
		return 1
	}

	// Сохраняем текущее значение директории в переменную окружения OLDPWD
	currentDir, _ := os.Getwd()
	os.Setenv("OLDPWD", currentDir)

	if err := os.Chdir(dir); err != nil {
		fmt.Printf("cd: %v\n", err)
		return 1
	}

	newDir, _ := os.Getwd()
	fmt.Println("Changed directory to:", newDir)

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

// GrepOptions хранит параметры grep
type GrepOptions struct {
	Pattern    string
	Filename   string
	WordMatch  bool
	IgnoreCase bool
	AfterLines int
}

// runGrep парсит флаги и вызывает executeGrep
func (e *Executor) runGrep(args []string) int {
	flags := pflag.NewFlagSet("grep", pflag.ContinueOnError)

	wordMatch := flags.BoolP("word", "w", false, "Искать только целые слова")
	ignoreCase := flags.BoolP("ignore-case", "i", false, "Игнорировать регистр")
	afterLines := flags.IntP("after", "A", 0, "Количество строк после совпадения")

	err := flags.Parse(args)
	if err != nil {
		fmt.Println("grep: error parsing flags")
		return 1
	}

	remainingArgs := flags.Args()
	if len(remainingArgs) < 2 {
		fmt.Println("grep: usage: grep [options] pattern file")
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

// executeGrep выполняет поиск
func (e *Executor) executeGrep(opts GrepOptions) int {
	file, err := os.Open(opts.Filename)
	if err != nil {
		fmt.Printf("grep: %v\n", err)
		return 1
	}
	defer file.Close()

	// Формируем шаблон регулярного выражения
	pattern := opts.Pattern
	if opts.WordMatch {
		pattern = `\b` + pattern + `\b`
	}
	if opts.IgnoreCase {
		pattern = `(?i)` + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("grep: invalid regex pattern: %v\n", err)
		return 1
	}

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	linesAfter := 0
	currentLine := 0

	// Читаем строки и фильтруем вывод
	for scanner.Scan() {
		line := scanner.Text()

		if re.MatchString(line) {
			linesAfter = max(linesAfter, opts.AfterLines) // пересечение областей
			buffer.WriteString(line + "\n")
		} else if linesAfter > 0 {
			buffer.WriteString(line + "\n")
			linesAfter--
		}

		currentLine++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("grep: error reading file: %v\n", err)
		return 1
	}

	fmt.Print(buffer.String())
	return 0
}
