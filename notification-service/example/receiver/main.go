package main

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("testing receiver")

	conn, err := amqp.Dial("amqp://admin:password@192.168.1.201/")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to RabbitMQ")
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close RabbitMQ connection")
			return
		}
		log.Info().Msg("closed RabbitMQ connection")
	}()
	log.Info().Msg("connected to RabbitMQ")

	// channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create RabbitMQ channel")
	}
	defer func() {
		err = ch.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close RabbitMQ channel")
			return
		}
		log.Info().Msg("closed RabbitMQ channel")
	}()
	log.Info().Msg("created RabbitMQ channel")
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set QoS")
	}
	log.Info().Msg("set QoS success")

	// exchange stuff
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		false,    // durability
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to declare exchange")
	}
	log.Info().Msg("exchange declared")

	// queue stuff
	log.Info().Msg("declaring queue")
	q, err := ch.QueueDeclare(
		"log_receiver",
		true,  // durability
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to declare queue")
	}
	log.Info().Str("name", q.Name).Msg("queue declared")

	// bind the queue
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bind queue")
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
		log.Fatal().Err(err).Msg("failed to consume queue")
	}

	for msg := range msgs {
		log.Info().Any("msg", msg).Str("body", string(msg.Body)).Msg("received message")
		err := msg.Ack(false)
		if err != nil {
			log.Error().Err(err).Msg("failed to ack message")
		}
	}
}
