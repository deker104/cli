package substitutor

import (
	"os"
	"strconv"
	"strings"

	"github.com/deker104/cli/internal/parser"
)

// Substitute заменяет переменные окружения в аргументах.
// Только если SubstituteEnv == true
func Substitute(args []parser.Token, lastExitCode int) []string {
	var result []string
	for _, arg := range args {
		if !arg.SubstituteEnv {
			result = append(result, arg.Value)
			continue
		}

		if strings.HasPrefix(arg.Value, "$") {
			if arg.Value == "$?" {
				result = append(result, strconv.Itoa(lastExitCode))
			} else {
				varName := arg.Value[1:]
				result = append(result, os.Getenv(varName))
			}
		} else {
			result = append(result, arg.Value)
		}
	}
	return result
}
