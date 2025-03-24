package parser

import (
	"strings"
)

// Parse разбирает команду, учитывая кавычки.
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

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
