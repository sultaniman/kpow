package enc

import (
	"os"

	"filippo.io/age"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/server"
)

type KeyLike interface {
	Encrypt(message string) (string, error)
}

func LoadKey(info *server.KeyInfo) KeyLike {
	content, err := os.ReadFile(info.Path)
	if err != nil {
		log.Fatal().Err(err)
	}
	// FIXME: move key & password validation logic from implementations to this function
	switch info.KeyKind {
	case server.Age:
		recipient, err := age.ParseX25519Recipient(string(content))
		if err != nil {
			log.Fatal().Err(err)
		}
		return NewAgeKey(recipient, info.Password)
	case server.PGP:
		pubkey, err := os.ReadFile(info.Path)
		if err != nil {
			log.Fatal().Err(err)
		}

		publicKey, err := crypto.NewKeyFromArmored(string(pubkey))
		if err != nil {
			log.Fatal().Err(err)
		}

		return NewPGPKey(publicKey, info.Password)
	default:
		log.Fatal().Msgf("Uknown key kind %v", info.KeyKind)
		return nil
	}
}
