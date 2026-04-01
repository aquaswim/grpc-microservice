# API Contracts

This is the central repository for API contracts (Protobuf definitions) for the Gaman microservices. It uses [Buf](https://buf.build/) to manage, lint, and format Protobuf files.

## Project Structure

- `proto/`: Contains the Protobuf definitions organized by service and version.
  - `proto/common/v1/`: Common Protobuf definitions (e.g., pagination).
  - `proto/event/v1/`: pubsub message definitions.
  - `proto/user/v1/`: User service Protobuf definitions.
- `buf.yaml`: Configuration for the Buf CLI.
- `devbox.json`: Devbox configuration for managing tools.

## Getting Started

We use [Devbox](https://github.com/jetify-com/devbox) to manage our development environment and tools.

## Available Scripts

We have defined several scripts in `devbox.json` to simplify common tasks:

### Linting

To lint the Protobuf files:

```bash
devbox run lint
```

### Formatting

To format the Protobuf files (this will overwrite changes in-place):

```bash
devbox run format
```

---

**Author:** [aquaswim](https://github.com/aquaswim)
