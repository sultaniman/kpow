.PHONY: build

build:
	go build -o kpow

fmt:
	gofmt -w .

test:
	go test -v ./...
