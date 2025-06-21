package mailer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type WebhookMailer struct {
	endpoint string
}

func (m *WebhookMailer) Send(message Message) error {
	payload := struct {
		Subject          string `json:"subject"`
		EncryptedMessage string `json:"content"`
		Hash             string `json:"hash"`
	}{
		Subject:          message.Subject,
		EncryptedMessage: message.EncryptedMessage,
		Hash:             message.Hash,
	}

	jsonMessage, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	response, err := http.Post(m.endpoint, "application/json", bytes.NewReader(jsonMessage))
	if err != nil {
		return err
	}
	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("webhook request error status=%d", response.StatusCode)
	}

	return nil
}

func NewWebhookMailer(webhookUrl string) (Mailer, error) {
	if webhookUrl == "" {
		return nil, nil
	}

	parts, err := url.Parse(webhookUrl)
	if err != nil {
		return nil, err
	}

	hostname := parts.Hostname()
	isLoopback := hostname == "localhost" || hostname == "127.0.0.1"
	if parts.Scheme != "https" && !isLoopback {
		return nil, errors.New("webhook url should be https only")
	}

	return &WebhookMailer{
		endpoint: webhookUrl,
	}, nil
}
