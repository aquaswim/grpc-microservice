package rabbitmq

import (
	"context"
	"gaman-microservice/notification-service/internal/pkg/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type message struct {
	msg amqp.Delivery
}

func (m message) GetData() []byte {
	return m.msg.Body
}

func (m message) GetID() string {
	return m.msg.MessageId
}

func (m message) Ack(ctx context.Context) {
	err := m.msg.Ack(false)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to ack the message")
	}
}

func wrapMessage(msg amqp.Delivery) pubsub.Message {
	return message{
		msg: msg,
	}
}
