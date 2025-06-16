package enc

import (
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

type PGPKey struct {
	Key *crypto.Key
}

func (k *PGPKey) Encrypt(message string) (string, error) {
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

func NewPGPKey(pubkeyBytes []byte) (KeyLike, error) {
	publicKey, err := crypto.NewKeyFromArmored(string(pubkeyBytes))
	if err != nil {
		return nil, err
	}

	return &PGPKey{publicKey}, nil
}
