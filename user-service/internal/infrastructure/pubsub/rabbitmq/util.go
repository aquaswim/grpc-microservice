package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func channelCloserWithLog(channel *amqp.Channel) {
	err := channel.Close()
	if err != nil {
		log.Err(err).Caller().Msg("failed to close channel")
	}
}
