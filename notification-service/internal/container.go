package internal

import (
	"gaman-microservice/notification-service/internal/config"
	globalLogger "gaman-microservice/notification-service/internal/pkg/global_logger"
	"gaman-microservice/notification-service/internal/pkg/pubsub"
	"gaman-microservice/notification-service/internal/pkg/pubsub/rabbitmq"
	"gaman-microservice/notification-service/internal/subscriber"

	"github.com/golobby/container/v3"
	amqp "github.com/rabbitmq/amqp091-go"
)

func InitContainer() container.Container {
	c := container.New()

	// config
	container.MustSingleton(c, config.Load)

	container.MustCall(c, func(cfg *config.Config) {
		globalLogger.Setup(&globalLogger.Config{
			LogPretty: cfg.LogPretty,
			LogLevel:  cfg.LogLevel,
		})
	})

	container.MustSingleton(c, func(cfg *config.Config) (pubsub.Client, error) {
		conn, err := amqp.Dial(cfg.RabbitMQUrl)
		if err != nil {
			return nil, err
		}

		return rabbitmq.New(conn, cfg.RabbitMqExchange)
	})

	container.MustSingleton(c, subscriber.New)

	return c
}
