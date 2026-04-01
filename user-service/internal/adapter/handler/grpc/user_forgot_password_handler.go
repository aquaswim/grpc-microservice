package grpc

import (
	"context"
	userv1 "gaman-microservice/user-service/gen/user/v1"
	"gaman-microservice/user-service/internal/port/in"
)

type UserForgotPasswordHandler struct {
	userv1.UnimplementedForgotPasswordServiceServer
	useCase in.UserForgotPasswordUseCase
}

func NewUserForgotPasswordHandler(useCase in.UserForgotPasswordUseCase) *UserForgotPasswordHandler {
	return &UserForgotPasswordHandler{
		useCase: useCase,
	}
}

func (h *UserForgotPasswordHandler) ForgotPassword(ctx context.Context, request *userv1.ForgotPasswordRequest) (*userv1.ForgotPasswordResponse, error) {
	err := h.useCase.ForgotPassword(ctx, request.GetEmail())
	if err != nil {
		return nil, err
	}

	return &userv1.ForgotPasswordResponse{
		Success: true,
		Message: "If the email exists, a reset token has been sent.",
	}, nil
}

func (h *UserForgotPasswordHandler) ValidateResetPasswordToken(ctx context.Context, request *userv1.ValidateResetPasswordTokenRequest) (*userv1.ValidateResetPasswordTokenResponse, error) {
	isValid, err := h.useCase.ValidateResetToken(ctx, request.GetToken())
	if err != nil {
		return nil, err
	}

	if !isValid {
		return &userv1.ValidateResetPasswordTokenResponse{
			Valid:   false,
			Message: "Invalid or expired token",
		}, nil
	}

	return &userv1.ValidateResetPasswordTokenResponse{
		Valid:   true,
		Message: "Token is valid",
	}, nil
}

func (h *UserForgotPasswordHandler) ResetPasswordWithToken(ctx context.Context, request *userv1.ResetPasswordWithTokenRequest) (*userv1.ResetPasswordWithTokenResponse, error) {
	err := h.useCase.ResetPassword(ctx, request.GetToken(), request.GetNewPassword())
	if err != nil {
		return nil, err
	}

	return &userv1.ResetPasswordWithTokenResponse{
		Success: true,
		Message: "Password has been reset successfully",
	}, nil
}
