# Dockerfile for notification-service
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service cmd/notification-service/main.go

# Run stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary
COPY --from=builder /app/notification-service .

# Command to run the binary
CMD ["./notification-service"]
