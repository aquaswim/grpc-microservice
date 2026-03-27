package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password string
	Email    string
}

func NewUserWithAutoId() *User {
	return &User{
		ID: uuid.Must(uuid.NewV7()).String(),
	}
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type PasswordResetToken struct {
	ID        int64
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}
