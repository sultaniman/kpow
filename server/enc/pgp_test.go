package enc

import (
	"os"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	secretMessage       = "secret message"
	gpgPubkey           = "testkeys/pubkey.pub"
	gpgPrivkey          = "testkeys/priv.gpg"
	pgpEncryptedMessage = `
-----BEGIN PGP MESSAGE-----

hQEMA7anQ2ruHCCRAQf+NGzPlYqNoGj005a/8SPlPKMQhLQPaRHq+U0hG4mnUvM2
APm0VHTMqbWd8hzg/8GcRu7qiCzoSPpujlEN+928Wlqh1+rWrFJY/EKms/fTJebq
7jJ6XWVaVH3DBXhduP7eyiiJi0lP7tprQ9hu+9Pjo7Ydxk6rbMYVIx5OjHobjYoR
mnhQgmnGLXCN4Gkr944KOkXl6la5APWspXHM+EmDNmmDuNSXO4YbfGZY80EaKyNl
OW3IOQ5IMR+V5BXhSofqOv4kdbJZVlh614IKNQWfNKCxGPPyP0kFlZM0CmzXLsuJ
cNd6dDdx5/F9xBm3PndOC6BKdbgLbl/RB5b/63AfvtJAAQHE1iiC1stSrN3JcE96
Zl914EjbJXjbUWDWwxU0g0sWnCKbCQ6IgDST/KP25eM+pfHPy7TYjTWbgPtCeiPN
/A==
=nDfZ
-----END PGP MESSAGE-----`
)

var (
	gpgPublicKey  *crypto.Key
	gpgPrivateKey *crypto.Key
	pgp           = crypto.PGP()
)

func init() {
	privateKeyBytes, err := os.ReadFile(gpgPrivkey)
	if err != nil {
		panic(err)
	}

	pubKeyBytes, err := os.ReadFile(gpgPubkey)
	if err != nil {
		panic(err)
	}

	privateKeyInstance, err := crypto.NewPrivateKeyFromArmored(string(privateKeyBytes), nil)
	if err != nil {
		panic(err)
	}

	publicKeyInstance, err := crypto.NewKeyFromArmored(string(pubKeyBytes))
	if err != nil {
		panic(err)
	}

	gpgPrivateKey = privateKeyInstance
	gpgPublicKey = publicKeyInstance
}

func TestPGPDecryptWithKey(t *testing.T) {
	decHandle, _ := pgp.
		Decryption().
		DecryptionKey(gpgPrivateKey).
		New()

	decResult, _ := decHandle.Decrypt([]byte(pgpEncryptedMessage), crypto.Armor)
	assert.Equal(t, "hello", string(decResult.Bytes()))
}

func TestPGPEncryptAndDecryptWithKey(t *testing.T) {
	pgpKey, _ := NewPGPKey(gpgPublicKey)
	encryptedMessage, _ := pgpKey.Encrypt(secretMessage)
	decHandle, _ := pgp.
		Decryption().
		DecryptionKey(gpgPrivateKey).
		New()

	decResult, _ := decHandle.Decrypt([]byte(encryptedMessage), crypto.Armor)
	assert.Equal(t, secretMessage, string(decResult.Bytes()))
}
