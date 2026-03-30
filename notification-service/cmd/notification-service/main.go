package main

import (
	"gaman-microservice/notification-service/internal"
	"gaman-microservice/notification-service/internal/pkg/utils"
	"gaman-microservice/notification-service/internal/subscriber"

	"github.com/rs/zerolog/log"
)

func main() {
	c := internal.InitContainer()

	log.Info().Msg("Starting notification service")

	subs := utils.Resolve[subscriber.Subscriber](c)

	err := subs.Listen()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start subscriber")
	}
}
