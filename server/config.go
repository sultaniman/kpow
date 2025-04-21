package server

type Mailer struct {
	Host        string
	Port        int
	User        string
	Password    string
	Retries     int
	SendTo      string
	FromAddress string
}

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
	Title    string
	Port     int
	LogLevel string
	KeyInfo  KeyInfo
}
