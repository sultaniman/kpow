package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const (
	encryptedMessageFormat = `
-----BEGIN RSA ENCRYPTED MESSAGE-----
%s
-----END RSA ENCRYPTED MESSAGE-----`
)

type RSAKey struct {
	pubkey         *rsa.PublicKey
	maxMessageSize int
}

func (k *RSAKey) Encrypt(message string) (string, error) {
	hash := sha256.New()

	msgBytes := []byte(message)
	if k.maxMessageSize > 0 && len(msgBytes) > k.maxMessageSize {
		msgBytes = msgBytes[:k.maxMessageSize]
	}

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, k.pubkey, msgBytes, nil)
	if err != nil {
		return "", err
	}

	b64Encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return fmt.Sprintf(encryptedMessageFormat, b64Encoded), nil
}

func NewRSAKey(pubkeyBytes []byte) (KeyLike, error) {
	block, _ := pem.Decode(pubkeyBytes)
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubkey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("provided key is not rsa key")
	}

	maxSize := pubkey.Size() - (2*sha256.Size + 2)

	return &RSAKey{
		pubkey:         pubkey,
		maxMessageSize: maxSize,
	}, nil
}
