package container

import (
	"context"
	"gaman-microservice/user-service/internal/adapter/auth"
	eventProducer "gaman-microservice/user-service/internal/adapter/event_producer"
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/adapter/repository/postgres"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/infrastructure/pgsql"
	"gaman-microservice/user-service/internal/infrastructure/pubsub"
	"gaman-microservice/user-service/internal/infrastructure/pubsub/rabbitmq"
	"gaman-microservice/user-service/internal/port/out"
	"gaman-microservice/user-service/internal/usecase"
	"os"

	"github.com/golobby/container/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() container.Container {
	c := container.New()

	// Config
	container.MustSingleton(c, config.Load)

	// setup logging
	container.MustCall(c, func(cfg *config.Config) {
		if cfg.PrettyLog {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			log.Warn().Msgf("Pretty logging is enabled, this must only be used in local!")
		}
	})

	// DB
	container.MustSingleton(c, func(cfg *config.Config) (*pgxpool.Pool, error) {
		return pgsql.Connect(context.Background(), cfg.DatabaseUrl)
	})

	// pubsub
	container.MustSingleton(c, func(cfg *config.Config) (pubsub.Client, error) {
		conn, err := amqp.Dial(cfg.RabbitMQUrl)
		if err != nil {
			return nil, err
		}
		return rabbitmq.New(conn, cfg.RabbitMqExchange)
	})

	// Token Manager
	container.MustSingleton(c, func(cfg *config.Config) (out.TokenManager, error) {
		return auth.NewPasetoManager(cfg.TokenPrivateKey, cfg.TokenPublicKey, cfg.GetTokenExpiryDuration())
	})

	// event producer
	container.MustSingleton(c, eventProducer.NewPubsubProducer)

	// Repository
	container.MustSingleton(c, postgres.NewUserRepository)
	container.MustSingleton(c, postgres.NewPasswordResetTokenRepository)

	// UseCase
	container.MustSingleton(c, usecase.NewUserAuthUseCase)
	container.MustSingleton(c, usecase.NewManageUserUseCase)
	container.MustSingleton(c, usecase.NewUserForgotPasswordUseCase)

	// Handler
	container.MustSingleton(c, usergrpc.NewUserManageHandler)
	container.MustSingleton(c, usergrpc.NewUserAuthHandler)
	container.MustSingleton(c, usergrpc.NewUserForgotPasswordHandler)

	return c
}
