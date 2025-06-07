package enc

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/stretchr/testify/assert"
)

const (
	agePublicKey           = "age124ule2zpfku4tp604awadgq9myqsm7fw87kygl2ge7lwzrw9xawqg7xdyk"
	agePrivateKey          = "AGE-SECRET-KEY-1JWYKFGEH6FFDSS2XJA80GYYTKE390X4LE6REJ3HUP69XDG7R446QHSRLDJ"
	preEncryptedAgeMessage = `
-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSA3WnpnL0h6OHp6bDEySHFq
azlWWW9EUmxNUXk4a29NOElUWERFSXJaNXhVClV4d3pocXNGL2ZNZnRwV0RvY3NT
STUzRHBmRzFJcEhhdmtFTGZabkcyZmcKLS0tIFhvNnJReVBnOXc4Zm1FOFdZdHhB
SzlEblMralZqcmloSXRucHNSV2Fqc1EK3KFfOX1Ln968kq1tX1iaI+9RoSqekVOF
na03n83y9DttvF2XOw==
-----END AGE ENCRYPTED FILE-----`
)

// Sanity check for key encryption
func TestAgeDecryptWithKey(t *testing.T) {
	_, identity := getAgeCredentials()
	decryptedMessage := decrypAgeMessage(preEncryptedAgeMessage, identity)
	assert.Equal(t, decryptedMessage, "hello")
}

func TestAgeEncryptDecryptWithKey(t *testing.T) {
	recipient, identity := getAgeCredentials()
	ageKey, _ := NewAgeKey(recipient)
	encryptedMessage, _ := ageKey.Encrypt("message")
	assert.Contains(t, encryptedMessage, "BEGIN AGE ENCRYPTED")
	decryptedMessage := decrypAgeMessage(encryptedMessage, identity)
	assert.Equal(t, decryptedMessage, "message")
}

func getAgeCredentials() (age.Recipient, age.Identity) {
	recipient, _ := age.ParseX25519Recipient(agePublicKey)
	identity, _ := age.ParseX25519Identity(agePrivateKey)
	return recipient, identity
}

func decrypAgeMessage(message string, identity age.Identity) string {
	out := &bytes.Buffer{}
	f := strings.NewReader(message)
	armorReader := armor.NewReader(f)
	reader, _ := age.Decrypt(armorReader, identity)
	io.Copy(out, reader)
	return out.String()
}
