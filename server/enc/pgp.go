package enc

import (
	"errors"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

type PGPKey struct {
	Path     string
	Password string
}

func (k *PGPKey) Encrypt(message string) (string, error) {
	if k.Path == "" && k.Password == "" {
		return "", errors.New("path or password expected")
	}

	if k.Path != "" {
		return k.withPubKey(message)
	} else if k.Password != "" {
		return k.withPassword(message)
	}

	return "", errors.New(`
		Please specify path and passphrase (if set)
		or set password for password encryption
	`)
}

func (k *PGPKey) withPubKey(message string) (string, error) {
	pubkey, err := os.ReadFile(k.Path)
	if err != nil {
		return "", err
	}

	publicKey, err := crypto.NewKeyFromArmored(string(pubkey))
	if err != nil {
		return "", err
	}

	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().Recipient(publicKey).New()
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

func NewPGPKey(path string, password string) *PGPKey {
	return &PGPKey{
		path,
		password,
	}
}
