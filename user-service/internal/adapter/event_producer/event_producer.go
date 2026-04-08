package eventProducer

import (
	"context"
	eventv1 "gaman-microservice/user-service/gen/event/v1"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/infrastructure/pubsub"
	"gaman-microservice/user-service/internal/port/out"

	"google.golang.org/protobuf/proto"
)

type producer struct {
	cfg    *config.Config
	client pubsub.Client
}

func NewPubsubProducer(
	cfg *config.Config,
	client pubsub.Client,
) out.EventProducer {
	return &producer{
		cfg:    cfg,
		client: client,
	}
}

func (p producer) ForgotPassword(ctx context.Context, data *entity.UserForgotPasswordData) error {
	if data.User == nil {
		return appError.ErrInternal.New("user data is nil")
	}
	msg, err := proto.Marshal(&eventv1.UserForgotPassword{
		ResetToken: data.Token,
		ExpiredAt:  data.ExpiredAt.Unix(),
		UserId:     data.User.ID,
		Username:   data.User.Username,
		Email:      data.User.Email,
	})
	if err != nil {
		return appError.ErrInternal.Wrap(err, "failed to create forgot password message")
	}

	return p.client.Publish(ctx, p.cfg.UserForgotPasswordTopic, msg)
}

func (p producer) UserResetPasswordDone(ctx context.Context, data *entity.UserResetPasswordDoneData) error {
	msg, err := proto.Marshal(&eventv1.UserResetPasswordDone{
		UserId:   data.UserID,
		Username: data.Username,
		Email:    data.Email,
	})
	if err != nil {
		return appError.ErrInternal.Wrap(err, "failed to create reset password done message")
	}
	return p.client.Publish(ctx, p.cfg.UserResetPasswordDoneTopic, msg)
}
