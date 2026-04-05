FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY . .

RUN set -eux; \
    if [ -f /app/backend/go.mod ]; then \
        cd /app/backend; \
    fi; \
    go mod download; \
    go mod tidy; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/happyhouse-api ./cmd/api

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /bin/happyhouse-api /usr/local/bin/happyhouse-api

ENV APP_PORT=8080
EXPOSE 8080

CMD ["happyhouse-api"]
