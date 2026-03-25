package grpc

import (
	"context"
	commonv1 "gaman-microservice/user-service/gen/common/v1"
	userv1 "gaman-microservice/user-service/gen/user/v1"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/in"
)

type UserManageHandler struct {
	userv1.UnimplementedManageServiceServer
	manageUserUseCase in.ManageUserUseCase
}

func NewUserManageHandler(userAuthUseCase in.UserAuthUseCase, manageUserUseCase in.ManageUserUseCase) *UserManageHandler {
	return &UserManageHandler{
		manageUserUseCase: manageUserUseCase,
	}
}

func (h *UserManageHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	user := entity.NewUserWithAutoId()
	user.Username = req.GetUsername()
	user.Email = req.GetEmail()
	if err := user.SetPassword(req.GetPassword()); err != nil {
		return nil, err
	}

	createdUser, err := h.manageUserUseCase.CreateUser(ctx, user)
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

func (h *UserManageHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := h.manageUserUseCase.GetUser(ctx, req.GetId())
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

func (h *UserManageHandler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	user := &entity.User{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	updatedUser, err := h.manageUserUseCase.UpdateUser(ctx, user)
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

func (h *UserManageHandler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	err := h.manageUserUseCase.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteUserResponse{
		Success: true,
	}, nil
}

func (h *UserManageHandler) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	var limit uint64 = 10
	var cursor string

	if req.GetPagination() != nil {
		if req.GetPagination().GetPageSize() > 0 {
			limit = uint64(req.GetPagination().GetPageSize())
		}
		cursor = req.GetPagination().GetPageToken()
	}

	users, err := h.manageUserUseCase.ListUsers(ctx, limit, cursor)
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
