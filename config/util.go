package config

import (
	"fmt"
	"strings"

	"github.com/labstack/gommon/color"
)

func LogErrors(errors []error) {
	listItem := color.Cyan("-")
	for _, err := range errors {
		fmt.Printf("%s %s\n", listItem, color.Red(err.Error()))
	}
}

func IsLocalhost(host string) bool {
	return strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1")
}
