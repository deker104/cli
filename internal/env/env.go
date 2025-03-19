package env

var variables = make(map[string]string)

// Set переменную окружения
func Set(name, value string) {
	variables[name] = value
}

// Get переменную окружения
func Get(name string) string {
	return variables[name]
}
