package pubsub

import "context"

type Client interface {
	Publish(ctx context.Context, topic string, message []byte) error
}
