# Gaman Microservice

Microservice boilerplate

## Project Structure

This repository is organized into several key directories:

- **[api-gateway/](api-gateway/)**: A gRPC-Gateway that provides a RESTful HTTP interface to internal gRPC services. It handles routing, middleware, and interceptors.
- **[user-service/](user-service/)**: The core User microservice implementing business logic, data persistence (PostgreSQL), and authentication. It follows Hexagonal Architecture.
- **[protos/](protos/)**: Centralized Protobuf definitions (API contracts) managed with [Buf](https://buf.build/). It serves as the single source of truth for service interfaces.

**Author:** [aquaswim](https://github.com/aquaswim)
