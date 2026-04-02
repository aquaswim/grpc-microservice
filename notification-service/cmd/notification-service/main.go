package main

import (
	"gaman-microservice/notification-service/internal"
	"gaman-microservice/notification-service/internal/pkg/utils"
	"gaman-microservice/notification-service/internal/subscriber"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	c := internal.InitContainer()

	log.Info().Msg("Starting notification service")

	subs := utils.Resolve[subscriber.Subscriber](c)

	go func() {
		err := subs.Listen()
		if err != nil {
			log.Error().Err(err).Msg("Failed to start subscriber")
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Received shutdown signal, stopping notification service...")

	// Stop subscribe
	err := subs.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to stop subscriber")
	}

	log.Info().Msg("Notification service stopped")

}
