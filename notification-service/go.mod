module gaman-microservice/notification-service

go 1.25.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1
	github.com/caarlos0/env/v11 v11.4.0
	github.com/golobby/container/v3 v3.3.2
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/rs/zerolog v1.35.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool google.golang.org/protobuf/cmd/protoc-gen-go
