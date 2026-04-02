package service

import (
	"context"
	emailTemplate "gaman-microservice/notification-service/email_template"
	"gaman-microservice/notification-service/internal/entity"

	"github.com/rs/zerolog/log"
)

type EmailService interface {
	SendForgotPasswordEmail(ctx context.Context, data *entity.ForgotPasswordNotificationData) error
}

type emailService struct {
}

func NewEmailService() EmailService {
	return &emailService{}
}

func (e emailService) SendForgotPasswordEmail(ctx context.Context, data *entity.ForgotPasswordNotificationData) error {
	l := log.Ctx(ctx)

	emailContent, err := emailTemplate.RenderTemplate("forgot_password", data)
	if err != nil {
		l.Error().Err(err).Msg("failed to render forgot password email template")
		return err
	}

	// send email
	l.Info().Str("email_content", emailContent).Msg("sending forgot password email")
	return nil
}
