package grpc

import (
	"context"
	commonv1 "gaman-microservice/user-service/gen/common/v1"
	appError "gaman-microservice/user-service/internal/domain/app_error"

	"google.golang.org/grpc/metadata"
)

func getAuthDataFromCtx(ctx context.Context) (*commonv1.TokenPayload, error) {
	out := &commonv1.TokenPayload{}
	userId := getFirst(metadata.ValueFromIncomingContext(ctx, "x-user-id"))
	if userId != nil {
		out.Id = *userId
	} else {
		return nil, appError.ErrUnauthorized.New("user id is empty")
	}
	userName := getFirst(metadata.ValueFromIncomingContext(ctx, "x-user-username"))
	if userName != nil {
		out.Username = *userName
	} else {
		return nil, appError.ErrUnauthorized.New("user name is empty")
	}

	return out, nil
}

func getFirst[T any](arr []T) *T {
	if len(arr) > 0 {
		return &arr[0]
	}
	return nil
}
