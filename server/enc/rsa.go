package enc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

const (
	encryptedMessageFormat = `
-----BEGIN RSA ENCRYPTED MESSAGE-----
%s
-----END RSA ENCRYPTED MESSAGE-----`
)

type RSAKey struct {
	pubkey *rsa.PublicKey
}

func (k *RSAKey) Encrypt(message string) (string, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, k.pubkey, []byte(message), nil)
	if err != nil {
		return "", err
	}

	b64Encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return fmt.Sprintf(encryptedMessageFormat, b64Encoded), nil
}

func NewRSAKey(pubkeyBytes []byte) (KeyLike, error) {
	parsedKey, err := x509.ParsePKIXPublicKey(pubkeyBytes)
	if err != nil {
		return nil, err
	}

	pubkey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("provided key is not rsa key")
	}

	return &RSAKey{
		pubkey: pubkey,
	}, nil
}
