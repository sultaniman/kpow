package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func NewWebhookMailer(endpoint string) *WebhookMailer {
	return &WebhookMailer{
		endpoint: endpoint,
	}
}
