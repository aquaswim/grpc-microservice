package subscriber

import (
	"context"
	eventv1 "gaman-microservice/notification-service/gen/event/v1"
	"gaman-microservice/notification-service/internal/config"
	"gaman-microservice/notification-service/internal/pkg/pubsub"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Subscriber interface {
	Listen() error
}

type subscriber struct {
	client pubsub.Client
	cfg    *config.Config
}

func New(
	client pubsub.Client,
	cfg *config.Config,
) Subscriber {
	return &subscriber{
		client: client,
		cfg:    cfg,
	}
}

func (s *subscriber) Listen() (err error) {
	err = s.client.Receive(s.cfg.UserForgotPasswordTopic, s.forgotPasswordHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to receive messages from forgot password queue")
	}

	return
}

func (s *subscriber) forgotPasswordHandler(ctx context.Context, msg pubsub.Message) error {
	l := log.Ctx(ctx)

	var event eventv1.UserForgotPassword
	if err := proto.Unmarshal(msg.GetData(), &event); err != nil {
		l.Error().Err(err).Msg("failed to unmarshal forgot password event")
		return err
	}
	// todo implement this
	l.Debug().Bytes("data", msg.GetData()).Msg("received forgot password event")

	msg.Ack(ctx)

	return nil
}
