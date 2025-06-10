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
	jsonMessage, err := json.Marshal(message)
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

	if parts.Scheme != "https" {
		return nil, errors.New("webhook url should be https only")
	}

	return &WebhookMailer{
		endpoint: webhookUrl,
	}, nil
}
