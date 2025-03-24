package parser

import (
	"strings"
)

// Parse разбирает строку, учитывая кавычки и экранированные символы.
func Parse(input string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		char := input[i]

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

		if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		// Экранирование кавычек
		if char == '\\' && i+1 < len(input) && (input[i+1] == '"' || input[i+1] == '\'') {
			i++
			char = input[i]
		}

		current.WriteByte(char)
	}

	// Если кавычки остались незакрытыми, включаем их в последний токен
	if inQuotes {
		tokens = append(tokens, string(quoteChar)+current.String())
	} else if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
