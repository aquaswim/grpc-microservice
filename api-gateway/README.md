# API Gateway

This project is a gRPC API Gateway for the **Gaman Microservice** ecosystem. It provides a RESTful HTTP interface to internal gRPC services using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

We use [Devbox](https://www.jetify.com/devbox) to manage our development environment. It ensures that everyone is using the same version of Go, Buf, and other tools.

## Getting Started

### 1. Generate Code

We use [Buf](https://buf.build/) for protocol buffer management and code generation. To generate the Go and gRPC-Gateway code, run:

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
| `PRETTY_LOG` | Enable console-friendly logging (zerolog ConsoleWriter) | `false` |
| `LISTEN_ADDR` | The address the HTTP server will listen on | `:8080` |
| `USER_SVC_ADDR` | The address of the User gRPC service (Required) | `localhost:50051` |
| `REDIS_ADDR` | The address of the Redis server | `localhost:6379` |
| `REDIS_DB` | The Redis database to use | `0` |
| `REDIS_PASS` | The password for the Redis server | (empty) |
| `REDIS_USER` | The username for the Redis server | (empty) |

### 3. Run the Application

Once you have generated the code and configured your `.env` file, you can start the gateway:

```bash
go run main.go
```

The gateway will be accessible at the address specified in `LISTEN_ADDR`.

## Project Structure
- `config/`: Configuration loading logic.
- `gen/`: Generated Go and gRPC-Gateway code.
- `interceptor/`: gRPC interceptors (unary and stream).
- `middleware/`: HTTP middlewares.
- `main.go`: Application entry point.

**Author:** [aquaswim](https://github.com/aquaswim)
