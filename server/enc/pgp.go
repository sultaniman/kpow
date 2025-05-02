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
	if k.Key == nil && k.Password == "" {
		return "", errors.New("path or password expected")
	}

	if k.Key != nil {
		return k.withPubKey(message)
	} else if k.Password != "" {
		return k.withPassword(message)
	}

	return "", errors.New("Please specify path or password")
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

func NewPGPKey(key *crypto.Key, password string) *PGPKey {
	return &PGPKey{
		key,
		password,
	}
}
