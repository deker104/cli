package env

import "os"

// EnvManager управляет переменными окружения CLI.
// Использует свою map для локальных переменных и системную os.Getenv для остального.
type EnvManager struct {
	vars map[string]string
}

// NewEnvManager создает новый экземпляр менеджера переменных окружения.
func NewEnvManager() *EnvManager {
	return &EnvManager{vars: make(map[string]string)}
}

// Get возвращает значение переменной окружения.
// Сначала ищет в собственной карте, потом в системных переменных.
func (e *EnvManager) Get(name string) string {
	if val, exists := e.vars[name]; exists {
		return val
	}
	return os.Getenv(name)
}

// Set устанавливает (или переопределяет) переменную окружения.
func (e *EnvManager) Set(name, value string) {
	e.vars[name] = value
}
