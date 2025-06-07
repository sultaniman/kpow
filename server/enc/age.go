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
	Recipient age.Recipient
}

func (k *AgeKey) Encrypt(message string) (string, error) {
	buf := &bytes.Buffer{}
	armorWriter := armor.NewWriter(buf)
	writer, err := age.Encrypt(armorWriter, k.Recipient)
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

func NewAgeKey(recipient age.Recipient) (*AgeKey, error) {
	if recipient == nil {
		return nil, errors.New("expected recipient expected got nil")
	}

	return &AgeKey{recipient}, nil
}
