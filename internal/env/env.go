package env

import "os"

// EnvManager управляет переменными окружения.
type EnvManager struct {
	vars map[string]string
}

// NewEnvManager создает новый менеджер окружения.
func NewEnvManager() *EnvManager {
	return &EnvManager{vars: make(map[string]string)}
}

// Get получает значение переменной окружения.
func (e *EnvManager) Get(name string) string {
	if val, exists := e.vars[name]; exists {
		return val
	}
	return os.Getenv(name)
}

// Set устанавливает переменную окружения.
func (e *EnvManager) Set(name, value string) {
	e.vars[name] = value
}
