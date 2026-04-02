package rabbitmq

import (
	"context"
	"fmt"
	"gaman-microservice/notification-service/internal/pkg/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type client struct {
	conn     *amqp.Connection
	exchange string
}

func (c client) Receive(topic string, fn func(ctx context.Context, msg pubsub.Message) error) error {
	// format gonna be
	//  main exchange
	//   routing: topic
	//   queue-name: q-<topic>-<service-name>
	l := log.With().
		Str("exchange", c.exchange).
		Str("topic", topic).
		Logger()

	// create channel
	ch, err := c.conn.Channel()
	if err != nil {
		l.Error().Err(err).Msg("failed to create channel")
		return err
	}
	err = ch.Qos(1, 0, false)
	if err != nil {
		l.Error().Err(err).Msg("failed to set qos")
		return err
	}

	// queue-stuff
	queueName := c.generateQueueName(topic)
	l.Debug().Str("queueName", queueName).Msg("declaring queue")
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durability
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		l.Error().Err(err).Msg("failed to declare queue")
		return err
	}

	// bind queue
	err = ch.QueueBind(
		q.Name,     // queue name
		topic,      // routing key
		c.exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to bind queue")
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to consume queue")
		return err
	}

	l.Info().Msg("started consuming messages...")

	for msg := range msgs {
		ll := l.
			With().
			Str("msq_id", msg.MessageId).
			Logger()
		ctx := ll.WithContext(context.TODO())
		err := fn(ctx, wrapMessage(msg))
		if err != nil {
			ll.Error().Err(err).Msg("error processing message")
			continue
		}
	}

	return nil
}

func (c client) generateQueueName(topic string) string {
	return fmt.Sprintf("q-%s-notification-svc", topic)
}

func New(conn *amqp.Connection, exchange string) (pubsub.Client, error) {
	l := log.With().Caller().Logger()
	ch, err := conn.Channel()
	if err != nil {
		l.Error().Err(err).Msg("failed to create channel")
		return nil, err
	}
	defer ch.Close()

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

func (c client) Stop() error {
	log.Info().Msg("closing rabbitmq connection")
	return c.conn.Close()
}
