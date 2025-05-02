package enc

import (
	"bytes"
	"errors"
	"io"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/rs/zerolog/log"
)

type AgeKey struct {
	recipient age.Recipient
	Password  string
}

func (k *AgeKey) Encrypt(message string) (string, error) {
	// try with recipient if set
	if k.recipient != nil {
		return k.withRecipient(message)
	}

	// else try with password
	if k.Password != "" {
		recipient, err := age.NewScryptRecipient(k.Password)
		if err != nil {
			log.Error().Err(err)
			return "", err
		}
		k.recipient = recipient
		return k.withRecipient(message)
	}

	return "", errors.New("no recipient or password provided")
}

func (k *AgeKey) withRecipient(message string) (string, error) {
	buf := &bytes.Buffer{}
	armorWriter := armor.NewWriter(buf)
	writer, err := age.Encrypt(armorWriter, k.recipient)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	if _, err := io.WriteString(writer, message); err != nil {
		log.Error().Msgf("Failed to write to encrypted message: %v", err)
		return "", err
	}

	if err := writer.Close(); err != nil {
		log.Error().Msgf("Failed to close encrypted message: %v", err)
		return "", err
	}

	if err := armorWriter.Close(); err != nil {
		log.Error().Msgf("Failed to close armor: %v", err)
		return "", err
	}

	return buf.String(), nil
}

func NewAgeKey(recipient age.Recipient, password string) *AgeKey {
	return &AgeKey{
		recipient,
		password,
	}
}
