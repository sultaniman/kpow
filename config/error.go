package config

import "fmt"

type ConfigError struct {
	What    string
	Message string
}

// Error implements error.
func (ce ConfigError) Error() string {
	return fmt.Sprintf("%s %s", ce.What, ce.Message)
}

func newConfigError(what string, message string) ConfigError {
	return ConfigError{What: what, Message: message}
}
