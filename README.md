# Gaman Microservice

Microservice boilerplate

## Project Structure

This repository is organized into several key directories:

- **[api-gateway/](api-gateway/)**: A gRPC-Gateway that provides a RESTful HTTP interface to internal gRPC services. It handles routing, middleware, and interceptors.
- **[user-service/](user-service/)**: The core User microservice implementing business logic, data persistence (PostgreSQL), and authentication. It follows Hexagonal Architecture.
- **[notification-service/](notification-service/)**: A service that handles sending notifications.
- **[protos/](protos/)**: Centralized Protobuf definitions (API contracts) managed with [Buf](https://buf.build/). It serves as the single source of truth for service interfaces.

**Author:** [aquaswim](https://github.com/aquaswim)

## General Recomendations

### Move each key directory to its own repository

To make your project more maintainable, it's recommended to move each key directory to its own repository.
This is what you need to do to achieve that:

1. Create a new repository for each key directory.
2. except for protos, update the buf.gen.yaml file to point to the protos in the new repository.