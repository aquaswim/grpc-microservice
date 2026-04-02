package pubsub

import "context"

type Message interface {
	GetData() []byte
	GetID() string
	Ack(ctx context.Context)
}

type Client interface {
	Receive(topic string, fn func(ctx context.Context, msg Message) error) error
	Stop() error
}
