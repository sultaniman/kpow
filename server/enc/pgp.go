package enc

import (
	"errors"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

type PGPKey struct {
	Key      *crypto.Key
	Password string
}

func (k *PGPKey) Encrypt(message string) (string, error) {
	if k.Key != nil {
		return k.withPubKey(message)
	}

	return k.withPassword(message)
}

func (k *PGPKey) withPubKey(message string) (string, error) {
	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().Recipient(k.Key).New()
	if err != nil {
		return "", err
	}

	pgpMessage, err := encHandle.Encrypt([]byte(message))
	if err != nil {
		return "", err
	}

	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return "", err
	}

	return string(armored), nil
}

func (k *PGPKey) withPassword(message string) (string, error) {
	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().Password([]byte(k.Password)).New()
	if err != nil {
		return "", err
	}

	pgpMessage, err := encHandle.Encrypt([]byte(message))
	if err != nil {
		return "", err
	}

	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return "", err
	}

	return string(armored), nil
}

func NewPGPKey(key *crypto.Key, password string) (*PGPKey, error) {
	if key == nil && password == "" {
		return nil, errors.New("path or password expected")
	}

	return &PGPKey{
		key,
		password,
	}, nil
}
