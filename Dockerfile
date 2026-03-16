FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o kpow

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
RUN adduser -D -h /app kpow

WORKDIR /app
COPY --from=builder /build/kpow .

USER kpow
EXPOSE 8080

ENTRYPOINT ["./kpow"]
CMD ["start", "--host", "0.0.0.0"]
