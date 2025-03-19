# CLI Interpreter

![Build and test workflow status](https://github.com/deker104/convex-hull/actions/workflows/ci.yml/badge.svg?branch=master)

## Описание
Этот проект представляет собой CLI-интерпретатор, поддерживающий команды:
- `cat [FILE]` — вывести содержимое файла
- `echo` — вывести аргументы
- `wc [FILE]` — посчитать строки, слова и байты
- `pwd` — вывести текущую директорию
- `exit` — выйти из интерпретатора
- Поддержка пайпов (`|`) и переменных окружения

## Установка
Требуется **Go 1.23**.  
Склонируйте репозиторий:
```sh
git clone https://github.com/deker104/cli.git
cd cli-interpreter
```
Запустите сборку:
```sh
go build -o cli-interpreter ./cmd
```
Запустите:
```sh
./cli-interpreter
```

## Запуск тестов
```sh
go test ./...
```

## CI/CD
Проект использует **GitHub Actions** для автоматического тестирования.
