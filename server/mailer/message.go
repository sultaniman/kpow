package mailer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

type Message struct {
	Subject          string `json:"subject"`
	EncryptedMessage string `json:"content"`
	Hash             string `json:"hash"`
}

func (m *Message) Save(basepath string) error {
	filepath := path.Join(basepath, m.Hash)
	messageBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(fmt.Sprintf("kpow-%s.json", filepath), messageBytes, 0644)
	}

	return nil
}

func NewMessage(subject string, encryptedMessage string, hash string) Message {
	return Message{
		Subject:          subject,
		EncryptedMessage: encryptedMessage,
		Hash:             hash,
	}
}
