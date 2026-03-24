# Dockerfile for api-gateway
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum from the current directory (api-gateway folder)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway main.go

# Run stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary
COPY --from=builder /app/api-gateway .

# Expose the port the gateway listens on
EXPOSE 8080

# Command to run the binary
CMD ["./api-gateway"]
