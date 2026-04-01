package grpcInterceptorUtil

import (
	"context"
	"gaman-microservice/api-gateway/constant"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"

	"google.golang.org/grpc/metadata"
)

func ValidateTokenFromContext(ctx context.Context, client userv1.AuthServiceClient) (*userv1.ValidateTokenResponse, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	token := ""
	if ok {
		tokens := md.Get(constant.MetadataKeuAuth)
		if len(tokens) > 0 {
			token = tokens[0]
		}
	}

	return client.ValidateToken(ctx, &userv1.ValidateTokenRequest{
		Token: token,
	})
}
