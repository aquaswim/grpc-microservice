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
	Close() error
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
	go s.mustCreateListener(s.cfg.UserForgotPasswordTopic, s.forgotPasswordHandler)

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
	l.Debug().
		Bytes("data", msg.GetData()).
		Str("token", event.GetResetToken()).
		Msg("received forgot password event")

	msg.Ack(ctx)

	return nil
}

func (s *subscriber) Close() error {
	return s.client.Stop()
}
