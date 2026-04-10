package service

import (
	"context"
	emailTemplate "gaman-microservice/notification-service/email_template"
	"gaman-microservice/notification-service/internal/client/email"
	"gaman-microservice/notification-service/internal/entity"

	"github.com/rs/zerolog/log"
)

type EmailService interface {
	SendForgotPasswordEmail(ctx context.Context, data *entity.ForgotPasswordNotificationData) error
	SendResetPasswordSuccessEvent(ctx context.Context, data entity.ResetPasswordSuccess) error
}

type emailService struct {
	mailClient email.Client
}

func NewEmailService(
	mailClient email.Client,
) EmailService {
	return &emailService{
		mailClient: mailClient,
	}
}

func (e emailService) SendForgotPasswordEmail(ctx context.Context, data *entity.ForgotPasswordNotificationData) error {
	l := log.Ctx(ctx)

	emailContent, err := emailTemplate.RenderTemplate("forgot_password.gohtml", data)
	if err != nil {
		l.Error().Err(err).Msg("failed to render forgot password email template")
		return err
	}

	sendEmailRes, err := e.mailClient.SendEmail(ctx, &email.SendEmailReq{
		Subject:  "Forgot Password",
		ToEmail:  data.Email,
		ToName:   data.Username,
		BodyHtml: emailContent,
		BodyText: "Use this link to reset your password",
	})
	if err != nil {
		l.Error().Err(err).Msg("failed to send forgot password email")
		return err
	}
	l.Info().Str("email_id", sendEmailRes.EmailId).Msg("forgot password email sent")
	return nil
}

func (e emailService) SendResetPasswordSuccessEvent(ctx context.Context, data entity.ResetPasswordSuccess) error {
	l := log.Ctx(ctx)
	//TODO implement me
	l.Info().Any("data", data).Msg("sending reset password success event")
	return nil
}
