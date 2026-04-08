package entity

import "time"

type ForgotPasswordNotificationData struct {
	Token     string
	ExpiredAt time.Time
	Username  string
	Email     string
}
