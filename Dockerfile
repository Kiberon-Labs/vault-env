FROM golang:1.24.3-alpine AS builder

WORKDIR /src
COPY . .

RUN apk add --no-cache git && \
    go mod download && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o "vault-env"

FROM alpine:3.21.3

WORKDIR /

COPY --from=builder "/src/vault-env" "/"

ENTRYPOINT ["/vault-env"]