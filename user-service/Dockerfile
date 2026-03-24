# Dockerfile for user-service
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service cmd/server/main.go

# Run stage
FROM alpine:latest

# Install migrate tool for running migrations (optional, can also be done manually)
# For now, we assume migrations are run or handled. 
# But let's copy the migrations folder just in case.
WORKDIR /root/

# Copy the pre-built binary
COPY --from=builder /app/user-service .
COPY --from=builder /app/db/migrations ./db/migrations

# Expose the gRPC port
EXPOSE 50051

# Command to run the binary
CMD ["./user-service"]
