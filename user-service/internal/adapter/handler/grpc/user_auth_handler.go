package grpc

import (
	"context"
	userv1 "gaman-microservice/user-service/gen/user/v1"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/port/in"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
)

type UserAuthHandler struct {
	userv1.UnimplementedAuthServiceServer
	userAuthUseCase   in.UserAuthUseCase
	manageUserUseCase in.ManageUserUseCase
}

func NewUserAuthHandler(userAuthUseCase in.UserAuthUseCase, manageUserUseCase in.ManageUserUseCase) *UserAuthHandler {
	return &UserAuthHandler{
		userAuthUseCase:   userAuthUseCase,
		manageUserUseCase: manageUserUseCase,
	}
}

func (h *UserAuthHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	_, token, err := h.userAuthUseCase.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userv1.LoginResponse{
		Token:     token.Token,
		ExpiresAt: token.Expiry.Format(time.RFC3339),
	}, nil
}

func (h *UserAuthHandler) ValidateToken(ctx context.Context, request *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	token := request.GetToken()
	if token == "" {
		// get token from metadata
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if authToken, ok := md["authorization"]; ok {
				token = strings.TrimPrefix(authToken[0], "Bearer ")
			}
		}
	}

	// check if token still empty
	if token == "" {
		return nil, appError.ErrValidation.New("token is empty")
	}

	tokenData, err := h.userAuthUseCase.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &userv1.ValidateTokenResponse{
		Token:    token,
		Id:       tokenData.Id,
		Username: tokenData.Username,
	}, nil
}

func (h *UserAuthHandler) GetMyProfile(ctx context.Context, _ *userv1.GetMyProfileRequest) (*userv1.GetMyProfileResponse, error) {
	validateTokenResponse, err := h.ValidateToken(ctx, &userv1.ValidateTokenRequest{})
	if err != nil {
		return nil, err
	}

	userData, err := h.manageUserUseCase.GetUser(ctx, validateTokenResponse.Id)
	if err != nil {
		return nil, err
	}

	return &userv1.GetMyProfileResponse{
		User: convertUserEntityToGrpcUser(userData),
	}, nil
}
