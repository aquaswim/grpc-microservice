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
| `TOKEN_PRIVATE_KEY` | Private key for asymmetric token signing (Required) | - |
| `TOKEN_PUBLIC_KEY` | Public key for asymmetric token verification (Optional, can be derived from private key if empty) | - |
| `TOKEN_EXPIRY_MINUTES` | Token expiration time in minutes | `60` |
| `RESET_TOKEN_EXPIRY_MINUTES` | Password reset token expiration time in minutes | `10` |

### 4. Run the Application

Once you have generated the code, ran migrations, and configured your `.env` file, you can start the service:

```bash
go run cmd/server/main.go
```

The service will be accessible at the address specified in `TCP_LISTENER_URL`.

### 5. Generate Asymmetric Keys (Optional)

We use PASETO V4 Public (asymmetric) for token signing. You can generate a new private and public key pair using:

```bash
go run cmd/token-gen/main.go
```

The output will provide both `private` and `public` keys in hex format, which you can use in your `.env` file.

## Project Structure

- `cmd/`: Application entry points.
    - `server/`: The gRPC server.
    - `token-gen/`: Utility to generate PASETO asymmetric keys.
- `db/`: Database migrations and schema.
- `gen/`: Generated Go and gRPC code.
- `internal/`: Internal application logic (Hexagonal Architecture).
    - `adapter/`: External adapters (gRPC handlers, repositories).
    - `domain/`: Core domain entities and errors.
    - `infrastructure/`: Infrastructure concerns (config, database connectors).
    - `port/`: Input and output ports (interfaces).
    - `usecase/`: Application business logic.

**Author:** [aquaswim](https://github.com/aquaswim)
