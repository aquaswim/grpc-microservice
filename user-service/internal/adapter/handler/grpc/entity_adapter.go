package grpc

import (
	userv1 "gaman-microservice/user-service/gen/user/v1"
	"gaman-microservice/user-service/internal/domain/entity"
)

func convertUserEntityToGrpcUser(user *entity.User) (out *userv1.User) {
	if user == nil {
		return
	}
	out = &userv1.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return
}
