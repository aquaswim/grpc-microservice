package email

import "context"

type Client interface {
	SendEmail(ctx context.Context, req *SendEmailReq) (*SendEmailRes, error)
}
