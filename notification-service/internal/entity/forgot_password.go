package entity

import "time"

type ForgotPasswordNotificationData struct {
	Token     string
	ExpiredAt time.Time
	Username  string
	Email     string
}

type ResetPasswordSuccess struct {
	UserId   string
	Username string
	Email    string
}
