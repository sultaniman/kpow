package mailer_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sultaniman/kpow/server/mailer"
	"go.uber.org/mock/gomock"
)

func TestSendMessageFallback(t *testing.T) {
	inbox := t.TempDir()
	msg := mailer.NewMessage("subj", "cipher", "hash123")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	failing := mailer.NewMockMailer(ctrl)
	failing.EXPECT().Send(msg).Return(errors.New("boom"))

	err := mailer.SendMessage(msg, failing, nil, inbox)
	assert.Error(t, err)

	path := filepath.Join(inbox, "kpow-"+msg.Hash+".json")
	data, readErr := os.ReadFile(path)
	assert.NoError(t, readErr)

	var saved mailer.Message
	assert.NoError(t, json.Unmarshal(data, &saved))
	assert.Equal(t, "mailer", saved.Method)
	assert.Equal(t, 0, saved.Retries)
}

func TestWebhookMailerSend(t *testing.T) {
	type payload struct {
		Subject          string `json:"subject"`
		EncryptedMessage string `json:"content"`
		Hash             string `json:"hash"`
	}

	received := make(chan payload, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var p payload
		json.NewDecoder(r.Body).Decode(&p)
		received <- p
	}))
	defer srv.Close()

	wh, err := mailer.NewWebhookMailer(srv.URL)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mailer.NewMockMailer(ctrl)
	mock.EXPECT().Send(gomock.Any()).Return(nil)

	msg := mailer.NewMessage("subj", "cipher", "hash123")
	assert.NoError(t, mailer.SendMessage(msg, mock, wh, t.TempDir()))
	assert.Equal(t, msg.Subject, (<-received).Subject)
}
