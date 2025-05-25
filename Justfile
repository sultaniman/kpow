default: build

clean:
	go mod tidy
	rm -f kpow

build:
	go build -o kpow

fmt:
	gofmt -w .

test:
	go test -v ./...

styles:
	bunx @tailwindcss/cli -m -i ./styles/kpow.css -o ./server/public/kpow.min.css --watch
