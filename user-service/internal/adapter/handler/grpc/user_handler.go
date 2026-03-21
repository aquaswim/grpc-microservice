package grpc

import (
	"context"
	"errors"
	commonv1 "gaman-microservice/user-service/gen/common/v1"
	userv1 "gaman-microservice/user-service/gen/user/v1"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/in"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
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
		Token:     token.Token,
		ExpiresAt: token.Expiry.Format(time.RFC3339),
	}, nil
}

func (h *UserHandler) ValidateToken(ctx context.Context, request *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
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
		return nil, errors.New("token is empty")
	}

	tokenData, err := h.userUseCase.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &userv1.ValidateTokenResponse{
		Token:    token,
		Id:       tokenData.Id,
		Username: tokenData.Username,
	}, nil
}

func (h *UserHandler) GetMyProfile(ctx context.Context, _ *userv1.GetMyProfileRequest) (*userv1.GetMyProfileResponse, error) {
	validateTokenResponse, err := h.ValidateToken(ctx, &userv1.ValidateTokenRequest{})
	if err != nil {
		return nil, err
	}

	getUserResponse, err := h.GetUser(ctx, &userv1.GetUserRequest{Id: validateTokenResponse.Id})
	if err != nil {
		return nil, err
	}

	return &userv1.GetMyProfileResponse{
		User: getUserResponse.User,
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	user := entity.NewUserWithAutoId()
	user.Username = req.GetUsername()
	user.Email = req.GetEmail()
	if err := user.SetPassword(req.GetPassword()); err != nil {
		return nil, err
	}

	createdUser, err := h.userUseCase.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{
		User: &userv1.User{
			Id:       createdUser.ID,
			Username: createdUser.Username,
			Email:    createdUser.Email,
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := h.userUseCase.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	user := &entity.User{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	updatedUser, err := h.userUseCase.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:       updatedUser.ID,
			Username: updatedUser.Username,
			Email:    updatedUser.Email,
		},
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	err := h.userUseCase.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteUserResponse{
		Success: true,
	}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	var limit uint64 = 10
	var cursor string

	if req.GetPagination() != nil {
		if req.GetPagination().GetPageSize() > 0 {
			limit = uint64(req.GetPagination().GetPageSize())
		}
		cursor = req.GetPagination().GetPageToken()
	}

	users, err := h.userUseCase.ListUsers(ctx, limit, cursor)
	if err != nil {
		return nil, err
	}

	var pbUsers []*userv1.User
	for _, u := range users {
		pbUsers = append(pbUsers, &userv1.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		})
	}

	var nextPageToken string
	if uint64(len(users)) == limit && len(users) > 0 {
		nextPageToken = users[len(users)-1].ID
	}

	return &userv1.ListUsersResponse{
		Users: pbUsers,
		Pagination: &commonv1.PaginationResponse{
			NextPageToken: nextPageToken,
		},
	}, nil
}
