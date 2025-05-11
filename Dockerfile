# -------- Build stage --------
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Add go modules first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Build the binary
RUN go build -o osu-ha ./cmd/server

# -------- Runtime stage --------
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/osu-ha .

# Copy example env and config (user can override at runtime)
COPY .env.example .env
COPY config/config.yaml /config.yaml

# Expose default port (can be overridden)
EXPOSE 8081

# Entry point
CMD ["./osu-ha"]
