package server

type Kind string

const (
	PGP Kind = "pgp"
	Age Kind = "age"
)

// For kind=PGP and unset path
// password pgp encryption is used.
type KeyInfo struct {
	Path     string
	KeyKind  Kind
	Password string
}

type Config struct {
	Title     string
	Port      int
	LogLevel  string
	KeyInfo   KeyInfo
	MailerDSN string
}
