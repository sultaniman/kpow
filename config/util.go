package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func LogErrors(errors []error) {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

	items := list.
		New(errors).
		Enumerator(list.Dash).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)

	fmt.Println(items)
}

func IsLocalhost(host string) bool {
	return strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1")
}
