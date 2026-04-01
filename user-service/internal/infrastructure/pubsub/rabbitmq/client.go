package rabbitmq

import (
	"context"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/infrastructure/pubsub"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type client struct {
	conn     *amqp.Connection
	exchange string
}

func (c client) Publish(ctx context.Context, topic string, message []byte) error {
	// channel
	ch, err := c.conn.Channel()
	if err != nil {
		return appError.ErrInternal.Wrap(err, "failed to create channel")
	}
	defer channelCloserWithLog(ch)

	// send
	log.Ctx(ctx).Info().
		Str("topic", topic).
		Str("exchange", c.exchange).
		Msg("publishing message")
	err = ch.PublishWithContext(
		ctx,
		c.exchange,
		topic,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			//ContentType: "text/plain",
			MessageId: uuid.NewString(),
			Timestamp: time.Now(),
			Body:      message,
		},
	)
	if err != nil {
		return appError.ErrInternal.Wrap(err, "error publishing message")
	}
	return nil
}

func New(conn *amqp.Connection, exchange string) (pubsub.Client, error) {
	l := log.With().Caller().Logger()
	ch, err := conn.Channel()
	if err != nil {
		l.Error().Err(err).Msg("failed to create channel")
		return nil, err
	}
	defer channelCloserWithLog(ch)

	// exchange stuff
	l.Debug().Str("exchange", exchange).Msg("declaring exchange")
	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		false,    // durability
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		l.Error().Err(err).Msg("failed to declare exchange")
		return nil, err
	}
	return &client{
		conn:     conn,
		exchange: exchange,
	}, nil
}
