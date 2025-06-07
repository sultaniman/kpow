package enc

import (
	"os"

	"filippo.io/age"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
)

type KeyLike interface {
	Encrypt(message string) (string, error)
}

func LoadKey(info *config.KeyInfo) (KeyLike, error) {
	content, err := os.ReadFile(info.Path)
	if err != nil {
		log.Fatal().Err(err)
	}

	switch info.Kind {
	case config.Age:
		recipient, err := age.ParseX25519Recipient(string(content))
		if err != nil {
			log.Fatal().Err(err)
		}
		return NewAgeKey(recipient)
	case config.PGP:
		pubkey, err := os.ReadFile(info.Path)
		if err != nil {
			log.Fatal().Err(err)
		}

		publicKey, err := crypto.NewKeyFromArmored(string(pubkey))
		if err != nil {
			log.Fatal().Err(err)
		}

		return NewPGPKey(publicKey)
	default:
		log.Fatal().Msgf("Uknown key kind %v", info.Kind)
		return nil, nil
	}
}
