ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

ENV GOTOOLCHAIN=auto

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o rune-seer ./cmd/main.go

FROM golang:1.24

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/rune-seer /app/rune-seer

WORKDIR /app

COPY ./public /app/public

EXPOSE 8080

CMD ["/app/rune-seer"]
