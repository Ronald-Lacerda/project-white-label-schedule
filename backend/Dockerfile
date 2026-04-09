# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /schedule-api ./cmd/api

# Install golang-migrate for entrypoint use
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata bash

WORKDIR /app

# Binary
COPY --from=builder /schedule-api .

# Migrations
COPY --from=builder /root/go/bin/migrate /usr/local/bin/migrate
COPY migrations/ ./migrations/

# Uploads directory
RUN mkdir -p ./uploads

COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
