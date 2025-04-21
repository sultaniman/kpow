package enc

import (
	"os"

	"filippo.io/age"
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

	switch info.KeyKind {
	case server.Age:
		recipient, err := age.ParseX25519Recipient(string(content))
		if err != nil {
			log.Fatal().Err(err)
		}
		return NewAgeKey(recipient)
	case server.PGP:
		return NewPGPKey(info.Path, info.Password)
	default:
		log.Fatal().Msgf("Uknown key kind %v", info.KeyKind)
		return nil
	}
}
