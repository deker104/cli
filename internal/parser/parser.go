package parser

import "strings"

// Token — аргумент команды, с указанием, разрешена ли подстановка переменных
type Token struct {
	Value         string
	SubstituteEnv bool
}

// Parse разбирает строку в пайпы и токены с учетом кавычек и подстановок
func Parse(input string) [][]Token {
	var result [][]Token
	var tokens []Token
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	substitute := true

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '|' && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, Token{current.String(), substitute})
				current.Reset()
			}
			if len(tokens) > 0 {
				result = append(result, tokens)
				tokens = nil
			}
			continue
		}

		if (char == '"' || char == '\'') && (!inQuotes || quoteChar == char) {
			if inQuotes {
				if char == quoteChar {
					inQuotes = false
					continue
				}
			} else {
				inQuotes = true
				quoteChar = char
				substitute = (char == '"') // только слабые кавычки позволяют подстановку
				continue
			}
		}

		if char == '\\' && i+1 < len(input) && (input[i+1] == '"' || input[i+1] == '\'') {
			i++
			char = input[i]
		}

		if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, Token{current.String(), substitute})
				current.Reset()
				substitute = true
			}
			continue
		}

		current.WriteByte(char)
	}

	if current.Len() > 0 {
		tokens = append(tokens, Token{current.String(), substitute})
	}
	if len(tokens) > 0 {
		result = append(result, tokens)
	}

	return result
}
