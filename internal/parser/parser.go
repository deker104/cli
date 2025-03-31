package parser

import (
	"strings"
)

// Parse разбирает строку, учитывая кавычки, пайпы и пробелы между аргументами.
func Parse(input string) [][]string {
	var result [][]string
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '|' && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			if len(tokens) > 0 {
				result = append(result, tokens)
				tokens = nil
			}
			continue
		}

		if char == '"' || char == '\'' {
			if inQuotes {
				if char == quoteChar {
					inQuotes = false
					continue
				}
			} else {
				inQuotes = true
				quoteChar = char
				continue
			}
		}

		// Поддержка экранированных кавычек
		if char == '\\' && i+1 < len(input) && (input[i+1] == '"' || input[i+1] == '\'') {
			i++
			char = input[i]
		}

		if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	// НЕ разбиваем аргументы, если они были в кавычках
	if len(tokens) > 0 {
		result = append(result, tokens)
	}

	return result
}
