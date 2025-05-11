# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o osu-ha ./cmd/server

# Final image
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/osu-ha .
COPY .env.example .env

EXPOSE 8081

CMD ["./osu-ha"]
