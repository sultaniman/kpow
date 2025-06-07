package mailer

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type DummyMailer struct{}

func (m *DummyMailer) Send(message Message) error {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFA500")).
		Padding(1, 2).
		Align(lipgloss.Left)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555")).
		Render("💌 A NEW MESSAGE")

	subject := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(message.Subject)

	body := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(message.EncryptedMessage)

	combined := fmt.Sprintf("%s\n\n%s\n\n%s", title, subject, body)
	fmt.Println(border.Render(combined))

	return nil
}

func NewDummyMailer() (*DummyMailer, error) {
	return &DummyMailer{}, nil
}
