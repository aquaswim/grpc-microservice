package subscriber

import (
	"context"
	eventv1 "gaman-microservice/notification-service/gen/event/v1"
	"gaman-microservice/notification-service/internal/config"
	"gaman-microservice/notification-service/internal/entity"
	"gaman-microservice/notification-service/internal/pkg/pubsub"
	"gaman-microservice/notification-service/internal/service"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Subscriber interface {
	Listen() error
	Close() error
}

type subscriber struct {
	client       pubsub.Client
	cfg          *config.Config
	emailService service.EmailService
}

func New(
	client pubsub.Client,
	cfg *config.Config,
	emailService service.EmailService,
) Subscriber {
	return &subscriber{
		client:       client,
		cfg:          cfg,
		emailService: emailService,
	}
}

func (s *subscriber) Listen() (err error) {
	go s.mustCreateListener(s.cfg.UserForgotPasswordTopic, s.forgotPasswordHandler)
	go s.mustCreateListener(s.cfg.UserResetPasswordDoneTopic, s.resetPasswordDoneHandler)

	return
}

func (s *subscriber) Close() error {
	return s.client.Stop()
}

func (s *subscriber) forgotPasswordHandler(ctx context.Context, msg pubsub.Message) error {
	l := log.Ctx(ctx)

	// note: implement nack if error if needed
	defer msg.Ack(ctx)

	// validate the message if needed
	var payload eventv1.UserForgotPassword
	if err := proto.Unmarshal(msg.GetData(), &payload); err != nil {
		l.Error().Err(err).Msg("failed to unmarshal forgot password event")
		return err
	}

	forgotPasswordData := &entity.ForgotPasswordNotificationData{
		Token:     payload.GetResetToken(),
		Username:  payload.GetUsername(),
		Email:     payload.GetEmail(),
		ExpiredAt: time.Unix(payload.GetExpiredAt(), 0),
	}

	err := s.emailService.SendForgotPasswordEmail(ctx, forgotPasswordData)
	if err != nil {
		l.Error().Err(err).Msg("failed to send forgot password email")
		return err
	}

	return nil
}

func (s *subscriber) resetPasswordDoneHandler(ctx context.Context, msg pubsub.Message) error {
	l := log.Ctx(ctx)

	defer msg.Ack(ctx)

	var payload eventv1.UserResetPasswordDone
	if err := proto.Unmarshal(msg.GetData(), &payload); err != nil {
		l.Error().Err(err).Msg("failed to unmarshal forgot password event")
		return err
	}

	eventData := entity.ResetPasswordSuccess{
		UserId:   payload.GetUserId(),
		Username: payload.GetUsername(),
		Email:    payload.GetEmail(),
	}

	return s.emailService.SendResetPasswordSuccessEvent(ctx, eventData)
}
