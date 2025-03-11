package executor

import (
	"fmt"
)

// ExecuteCommand выполняет команду
func ExecuteCommand(command string, args []string) {
	fmt.Println("Executing:", command, args)
}