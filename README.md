# Gaman Microservice

Microservice boilerplate

## Project Structure

This repository is organized into several key directories:

- **[api-gateway/](api-gateway/)**: A gRPC-Gateway that provides a RESTful HTTP interface to internal gRPC services. It handles routing, middleware, and interceptors.
- **[user-service/](user-service/)**: The core User microservice implementing business logic, data persistence (PostgreSQL), and authentication. It follows Hexagonal Architecture.
- **[protos/](protos/)**: Centralized Protobuf definitions (API contracts) managed with [Buf](https://buf.build/). It serves as the single source of truth for service interfaces.

## Deployment

For instructions on how to deploy this microservice stack to MicroK8s, please refer to the [k8s/README.md](k8s/README.md).

### Independent Development

Each service (`api-gateway`, `user-service`) has its own `Dockerfile` and can be built and developed independently.
- To build `api-gateway`: `docker build -t api-gateway ./api-gateway`
- To build `user-service`: `docker build -t user-service ./user-service`

The `./k8s` folder contains the necessary Kubernetes configurations to run the entire stack, including a single instance of PostgreSQL.

**Author:** [aquaswim](https://github.com/aquaswim)
