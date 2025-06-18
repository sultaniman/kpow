package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	rsaPrivateKey *rsa.PrivateKey
)

const (
	secretRSAMessage    = "secret rsa message"
	rsaPubkeyPath       = "testkeys/public_rsa.pem"
	rsaPrivkeyPath      = "testkeys/private_rsa.pem"
	rsaEncryptedMessage = `qF4kahdKQWJtPK4NlXtfxCMJSavlc0KBZxROUboI6OX/nnDykEaGbL/eyEsIzVmDzbi/Z0l+AqwM98JCdFUBoiZFLy5LWCd+F4SxiGSZUn88KMqiWBIoggonbsbRJi8sWsr9X0FCJwgU7lNonTHNXsi7WZayhuuVhG5H6y16KtPHUT2CO1utr/l2ZwrRljIHL2TPW7bMzyag7IkytYZDNZg4aVDvLYDiamYCuaUcOu3pAkrookDkVZwck5RSpo+acuzo+iYiV6BT4cRnQ2tBQRDFJfcHDvbr/S3oinP02q/oMrmymwdB4oir3i3m2ZTgENPxI1+SIJrlNO+pu9dtfg==`
)

func TestRSADecryptWithKey(t *testing.T) {
	loadRSAKeys()

	// Test message is encrypted with sha256 hashing algorithm
	hash := sha256.New()
	ciphertext, _ := base64.StdEncoding.DecodeString(rsaEncryptedMessage)
	result, err := rsa.DecryptOAEP(hash, rand.Reader, rsaPrivateKey, ciphertext, nil)
	assert.NoError(t, err)
	assert.Equal(t, secretRSAMessage, string(result))
}

func TestRSAEncryptAndDecryptWithKey(t *testing.T) {
	loadRSAKeys()
	// Encrypt
	pubKeyBytes, _ := os.ReadFile(rsaPubkeyPath)
	rsaKey, _ := NewRSAKey(pubKeyBytes)
	encryptedMessage, err := rsaKey.Encrypt(secretRSAMessage)
	assert.NoError(t, err)
	encryptedMessage = strings.Replace(encryptedMessage, "-----BEGIN RSA ENCRYPTED MESSAGE-----", "", 1)
	encryptedMessage = strings.Replace(encryptedMessage, "-----END RSA ENCRYPTED MESSAGE-----", "", 1)
	encryptedMessage = strings.TrimSpace(encryptedMessage)

	// Decrypt
	hash := sha256.New()
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage)
	assert.NoError(t, err)
	result, err := rsa.DecryptOAEP(hash, rand.Reader, rsaPrivateKey, ciphertext, nil)
	assert.NoError(t, err)
	assert.Equal(t, secretRSAMessage, string(result))
}

func loadRSAKeys() {
	privateKeyBytes, _ := os.ReadFile(rsaPrivkeyPath)
	privateBlock, _ := pem.Decode(privateKeyBytes)
	privateKey, _ := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	rsaPrivateKey, _ = privateKey.(*rsa.PrivateKey)
}
