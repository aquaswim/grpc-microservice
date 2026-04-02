package subscriber

import (
	"context"
	"gaman-microservice/notification-service/internal/pkg/pubsub"

	"github.com/rs/zerolog/log"
)

type listenerFn func(ctx context.Context, msg pubsub.Message) error

func (s *subscriber) mustCreateListener(topic string, listener listenerFn) {
	err := s.client.Receive(topic, listener)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("topic", topic).
			Msg("failed create listener for topic")
	}
}
