package grpc

import (
	"context"
	commonv1 "gaman-microservice/user-service/gen/common/v1"
	userv1 "gaman-microservice/user-service/gen/user/v1"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/port/in"
	"strings"
	"time"
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

	// check if token still empty
	if token == "" {
		return nil, appError.ErrValidation.New("token is empty")
	}

	// strips the "Bearer " prefix
	token = strings.TrimPrefix(token, "Bearer ")

	tokenData, err := h.userAuthUseCase.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &userv1.ValidateTokenResponse{
		Token: token,
		Data:  &commonv1.TokenPayload{Id: tokenData.Id, Username: tokenData.Username},
	}, nil
}

func (h *UserAuthHandler) GetMyProfile(ctx context.Context, _ *userv1.GetMyProfileRequest) (*userv1.GetMyProfileResponse, error) {
	// get userId from context
	uCtx, err := getAuthDataFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userData, err := h.manageUserUseCase.GetUser(ctx, uCtx.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.GetMyProfileResponse{
		User: convertUserEntityToGrpcUser(userData),
	}, nil
}

func (h *UserAuthHandler) UpdateMyProfile(ctx context.Context, req *userv1.UpdateMyProfileRequest) (*userv1.UpdateMyProfileResponse, error) {
	uCtx, err := getAuthDataFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userData, err := h.manageUserUseCase.GetUser(ctx, uCtx.GetId())
	if err != nil {
		return nil, err
	}
	if req.GetEmail() != "" {
		userData.Email = req.GetEmail()
	}
	if req.GetUsername() != "" {
		userData.Username = req.GetUsername()
	}
	if req.GetPassword() != "" {
		err = userData.SetPassword(req.GetPassword())
		if err != nil {
			return nil, appError.ErrInternal.Wrap(err, "error setting password")
		}
	}

	_, err = h.manageUserUseCase.UpdateUser(ctx, userData)
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateMyProfileResponse{
		Success: true,
	}, nil
}
