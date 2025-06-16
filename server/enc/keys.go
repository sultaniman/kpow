package enc

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/config"
)

type KeyLike interface {
	Encrypt(message string) (string, error)
}

func LoadKey(info *config.KeyInfo) (KeyLike, error) {
	pubkeyBytes, err := os.ReadFile(info.Path)
	if err != nil {
		log.Fatal().Err(err)
	}

	switch info.Kind {
	case config.Age:
		return NewAgeKey(pubkeyBytes)
	case config.PGP:
		return NewPGPKey(pubkeyBytes)
	case config.RSA:
		return NewRSAKey(pubkeyBytes)
	default:
		log.Fatal().Msgf("Uknown key kind %v", info.Kind)
		return nil, nil
	}
}
