# User Service

This project is the **User Service** for the **Gaman Microservice** ecosystem.

We use [Devbox](https://www.jetify.com/devbox) to manage our development environment, ensuring consistent versions of Buf, protoc, and other tools.

## Getting Started

### 1. Generate Code

We use [Buf](https://buf.build/) for protocol buffer management and code generation. To generate the Go and gRPC code, run:

```bash
devbox run gen
```

Or simply:

```bash
buf generate
```

### 2. Database Migrations

This service uses [dbmate](https://github.com/amacneil/dbmate) for database migrations. To run migrations, ensure your `DATABASE_URL` is set (see Configuration) and run:

```bash
dbmate up
```

### 3. Configuration

Copy the `.env.example` file to `.env` and adjust the values as needed:

```bash
cp .env.example .env
```

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PRETTY_LOG` | Enable console-friendly logging (zerolog ConsoleWriter) | `false` |
| `TCP_LISTENER_URL` | The address the gRPC server will listen on | `:50051` |
| `DATABASE_URL` | PostgreSQL connection string (Required) | - |
| `TOKEN_SECRET` | Secret key for token generation (Required) | - |
| `TOKEN_EXPIRY_MINUTES` | Token expiration time in minutes | `60` |

### 4. Run the Application

Once you have generated the code, ran migrations, and configured your `.env` file, you can start the service:

```bash
go run cmd/server/main.go
```

The service will be accessible at the address specified in `TCP_LISTENER_URL`.

## Project Structure

- `cmd/`: Application entry points.
- `db/`: Database migrations and schema.
- `gen/`: Generated Go and gRPC code.
- `internal/`: Internal application logic (Hexagonal Architecture).
    - `adapter/`: External adapters (gRPC handlers, repositories).
    - `domain/`: Core domain entities and errors.
    - `infrastructure/`: Infrastructure concerns (config, database connectors).
    - `port/`: Input and output ports (interfaces).
    - `usecase/`: Application business logic.

**Author:** [aquaswim](https://github.com/aquaswim)
