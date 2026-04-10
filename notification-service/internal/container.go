package internal

import (
	"gaman-microservice/notification-service/internal/client/email"
	"gaman-microservice/notification-service/internal/client/email/mailpit"
	"gaman-microservice/notification-service/internal/config"
	globalLogger "gaman-microservice/notification-service/internal/pkg/global_logger"
	loggedHttpclient "gaman-microservice/notification-service/internal/pkg/logged_httpclient"
	"gaman-microservice/notification-service/internal/pkg/pubsub"
	"gaman-microservice/notification-service/internal/pkg/pubsub/rabbitmq"
	"gaman-microservice/notification-service/internal/service"
	"gaman-microservice/notification-service/internal/subscriber"
	"net/http"

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

	container.MustSingleton(c, loggedHttpclient.New)

	// 3rd party client
	container.MustSingleton(c, func(cfg *config.Config, client *http.Client) email.Client {
		return mailpit.New(client, cfg.MailpitUrl)
	})

	// services
	container.MustSingleton(c, service.NewEmailService)

	container.MustSingleton(c, subscriber.New)

	return c
}
