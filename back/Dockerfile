# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install git (for go mod), and ca-certs
RUN apk add --no-cache git ca-certificates

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.19
WORKDIR /app

# Add curl for docker healthcheck
RUN apk add --no-cache ca-certificates curl && update-ca-certificates

COPY --from=builder /app/server /app/server

# Environment variables (overridden by docker-compose)
ENV HTTP_PORT=8080

EXPOSE 8080

CMD ["/app/server"]
