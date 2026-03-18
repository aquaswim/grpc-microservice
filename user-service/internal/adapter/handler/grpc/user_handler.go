package grpc

import (
	"context"
	userv1 "gaman-microservice/user-service/gen/user/v1"
	"gaman-microservice/user-service/internal/port/in"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	userUseCase in.UserUseCase
}

func NewUserHandler(userUseCase in.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	_, token, err := h.userUseCase.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userv1.LoginResponse{
		Token:     token,
		ExpiresAt: "2026-03-18T23:59:59Z",
	}, nil
}

func (h *UserHandler) Profile(ctx context.Context, req *userv1.ProfileRequest) (*userv1.ProfileResponse, error) {
	user, err := h.userUseCase.GetProfile(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &userv1.ProfileResponse{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
