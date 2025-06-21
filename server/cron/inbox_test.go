package cron

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sultaniman/kpow/server/mailer"
	"go.uber.org/mock/gomock"
)

// TestInboxCleanerMailerAndWebhook verifies that messages saved with the
// "mailer" method trigger both mail and webhook delivery attempts. The inbox
// file should remain because only webhook-only messages are removed on success.
func TestInboxCleanerMailerAndWebhook(t *testing.T) {
	inbox := t.TempDir()
	msg := mailer.Message{Subject: "s", EncryptedMessage: "c", Hash: "m1", Method: "mailer"}
	assert.NoError(t, msg.Save(inbox))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mail := mailer.NewMockMailer(ctrl)
	webhook := mailer.NewMockMailer(ctrl)

	mail.EXPECT().Send(gomock.Any()).Return(nil)
	webhook.EXPECT().Send(gomock.Any()).Return(nil)

	cleaner := InboxCleaner(inbox, mail, webhook)
	cleaner()

	_, err := os.Stat(filepath.Join(inbox, "kpow-"+msg.Hash+".json"))
	assert.NoError(t, err)
}

// TestInboxCleanerWebhookSuccessRemovesFile ensures webhook-only messages are
// deleted after a successful webhook delivery attempt.
func TestInboxCleanerWebhookSuccessRemovesFile(t *testing.T) {
	inbox := t.TempDir()
	msg := mailer.Message{Subject: "s", EncryptedMessage: "c", Hash: "m2", Method: "webhook"}
	assert.NoError(t, msg.Save(inbox))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	webhook := mailer.NewMockMailer(ctrl)
	webhook.EXPECT().Send(gomock.Any()).Return(nil)

	cleaner := InboxCleaner(inbox, nil, webhook)
	cleaner()

	_, err := os.Stat(filepath.Join(inbox, "kpow-"+msg.Hash+".json"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

// TestInboxCleanerWebhookFailureKeepsFile checks that inbox files remain when
// webhook delivery fails so the scheduler can retry later.
func TestInboxCleanerWebhookFailureKeepsFile(t *testing.T) {
	inbox := t.TempDir()
	msg := mailer.Message{Subject: "s", EncryptedMessage: "c", Hash: "m3", Method: "webhook"}
	assert.NoError(t, msg.Save(inbox))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	webhook := mailer.NewMockMailer(ctrl)
	webhook.EXPECT().Send(gomock.Any()).Return(errors.New("boom"))

	cleaner := InboxCleaner(inbox, nil, webhook)
	cleaner()

	_, err := os.Stat(filepath.Join(inbox, "kpow-"+msg.Hash+".json"))
	assert.NoError(t, err)
}
