package entity

import "time"

type TokenData struct {
	Id       string
	Username string
}

type TokenWithExpiry struct {
	Token  string
	Expiry time.Time
}
