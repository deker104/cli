package substitutor

import (
	"os"
	"strconv"
	"strings"
)

// Substitute заменяет переменные окружения в аргументах.
func Substitute(args []string, lastExitCode int) []string {
	var result []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "$") {
			if arg == "$?" {
				result = append(result, strconv.Itoa(lastExitCode))
			} else {
				varName := arg[1:]
				result = append(result, os.Getenv(varName))
			}
		} else {
			result = append(result, arg)
		}
	}
	return result
}
