package server

import (
	"fmt"

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
