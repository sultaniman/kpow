default: setup-tools check build

export TEST_KEYS_DIR := invocation_directory() + "/server/enc/testkeys"

setup-tools:
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    go install github.com/air-verse/air@latest

clean:
    go mod tidy
    rm -f kpow

build:
    go build -o kpow

fmt:
    gofmt -w .

test:
    go test -v ./...

check:
    gosec ./...

dev:
    air

styles:
    bunx @tailwindcss/cli -m -i ./styles/kpow.css -o ./server/public/kpow.min.css --watch

error-styles:
    bunx @tailwindcss/cli -m -i ./styles/errors.css -o ./server/public/errors.min.css --watch
