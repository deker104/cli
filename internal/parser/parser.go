package parser

import "strings"

// ParseCommand разбирает входную строку и возвращает команду и аргументы
func ParseCommand(input string) []string {
	return strings.Fields(input)
}