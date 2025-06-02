package config

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

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

const advices = `
Your encryption key is generated from a passphrase.

🛡️  This is only as strong as the passphrase itself.
💥  Weak or reused passwords can be brute-forced if
    someone gets your encrypted data.

🔑  For stronger protection, use a real public key.
✅  Use a long and unique passphrase (preferably generated).`

func WarnAboutPassphrase() {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFA500")).
		Padding(1, 2).
		Align(lipgloss.Left)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555")).
		Render("⚠️  SECURITY WARNING: PASSWORD-DERIVED KEY IN USE")

	body := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(advices)

	combined := fmt.Sprintf("%s\n\n%s", title, body)
	fmt.Println(border.Render(combined))
}
