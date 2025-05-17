.PHONY: build

clean:
	go mod tidy
	rm -f kpow

build:
	go build -o kpow

fmt:
	gofmt -w .

test:
	go test -v ./...
