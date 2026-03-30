package entity

import "time"

type UserForgotPasswordData struct {
	User      *User
	Token     string
	ExpiredAt time.Time
}
