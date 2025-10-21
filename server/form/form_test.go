package form

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sultaniman/kpow/server/enc"
	"github.com/sultaniman/kpow/server/mailer"
	"go.uber.org/mock/gomock"
)

func TestEncryptAndSend_EncryptionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := enc.NewMockKeyLike(ctrl)
	key.EXPECT().Encrypt("body").Return("", errors.New("boom"))

	fd := &FormData{Message: MessageForm{Subject: "sub", Content: "body"}}

	errCh := fd.EncryptAndSend(nil, nil, key, t.TempDir())
	err := <-errCh
	assert.Error(t, err)
	assert.Equal(t, MessageForm{}, fd.Message)
}

func TestEncryptAndSend_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := enc.NewMockKeyLike(ctrl)
	key.EXPECT().Encrypt("body").Return("cipher", nil)

	mail := mailer.NewMockMailer(ctrl)
	mail.EXPECT().Send(gomock.Any()).Return(errors.New("boom"))

	fd := &FormData{Message: MessageForm{Subject: "sub", Content: "body"}}

	errCh := fd.EncryptAndSend(mail, nil, key, t.TempDir())
	err := <-errCh
	assert.Error(t, err)
	assert.Equal(t, MessageForm{}, fd.Message)
}

func TestEncryptAndSend_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := enc.NewMockKeyLike(ctrl)
	key.EXPECT().Encrypt("body").Return("cipher", nil)

	mail := mailer.NewMockMailer(ctrl)
	mail.EXPECT().Send(gomock.Any()).Return(nil)

	fd := &FormData{Message: MessageForm{Subject: "sub", Content: "body"}}

	errCh := fd.EncryptAndSend(mail, nil, key, t.TempDir())
	err := <-errCh
	assert.NoError(t, err)
	assert.Equal(t, MessageForm{}, fd.Message)
}
