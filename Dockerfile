# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /feature-flag-service ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl \
    && curl -fsSL https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
       | tar -xz -C /usr/local/bin migrate

WORKDIR /app

COPY --from=builder /feature-flag-service .
COPY migrations ./migrations
COPY entrypoint.sh ./entrypoint.sh

RUN chmod +x /app/entrypoint.sh

ENV PORT=8080
ENV MIGRATIONS_PATH=file://migrations
ENV DATABASE_URL=postgres://flaguser:flagpass@postgres:5432/feature_flags?sslmode=disable
ENV SKIP_MIGRATIONS=true

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
