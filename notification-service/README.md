# Notification Service

This project is the **Notification Service** for the **Gaman Microservice** ecosystem.

We use [Devbox](https://www.jetify.com/devbox) to manage our development environment, ensuring consistent versions of Buf, protoc, and other tools.

## Getting Started

### 1. Generate Code

We use [Buf](https://buf.build/) for protocol buffer management and code generation. To generate the Go code, run:

```bash
devbox run gen
```

Or simply:

```bash
buf generate
```

### 2. Configuration

Copy the `.env.example` file to `.env` and adjust the values as needed:

```bash
cp .env.example .env
```

| Variable | Description | Default |
| :--- | :--- | :--- |
| `LOG_PRETTY` | Enable console-friendly logging (zerolog ConsoleWriter) | `false` |
| `LOG_LEVEL` | Minimum log level (debug, info, error, etc.) | `info` |
| `RABBITMQ_URL` | RabbitMQ connection string (Required) | - |
| `RABBITMQ_EXCHANGE` | RabbitMQ exchange name (Required) | - |
| `TOPIC_USER_FORGOT_PASSWORD` | Topic name for forgot password events | `user-forgot-password` |

### 3. Run the Application

Once you have generated the code and configured your `.env` file, you can start the service:

```bash
go run cmd/notification-service/main.go
```

The service will start listening for messages on the configured topics.

## Project Structure

- `cmd/`: Application entry points.
    - `notification-service/`: The main application.
- `gen/`: Generated Go code from protocol buffers.
- `internal/`: Internal application logic.
    - `config/`: Configuration loading logic.
    - `entity/`: Domain entities.
    - `pkg/`: Reusable packages and utilities.
    - `service/`: Core business logic services.
    - `subscriber/`: Message queue event handlers and listeners.

**Author:** [aquaswim](https://github.com/aquaswim)
