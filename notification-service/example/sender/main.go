package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("testing sender")

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

	// send msg
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := fmt.Sprintf("message: #%d", rand.Int63())
	err = ch.PublishWithContext(ctx,
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			//ContentType: "text/plain",
			Body: []byte(body),
		})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to publish message")
	}
	log.Info().Msg("published message")
}
