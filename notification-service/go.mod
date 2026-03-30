module gaman-microservice/notification-service

go 1.25.0

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/rs/zerolog v1.35.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

tool google.golang.org/protobuf/cmd/protoc-gen-go
