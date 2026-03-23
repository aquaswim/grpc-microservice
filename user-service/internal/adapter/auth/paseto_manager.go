package auth

import (
	"context"
	"fmt"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/out"
	"time"

	"aidanwoods.dev/go-paseto"
)

type pasetoManager struct {
	symmetricKey   paseto.V4SymmetricKey
	expiryDuration time.Duration
}

func NewPasetoManager(secret string, expiryDuration time.Duration) (out.TokenManager, error) {
	// Paseto V4 symmetric key must be 32 bytes (256 bits)
	// If the secret is shorter, we can pad it or use its hash.
	// To be strictly correct according to aidanwoods.dev/go-paseto, we should use V4SymmetricKeyFromBytes.
	// However, if we want to allow user-defined strings, we can use a fixed-size hash.
	// For simplicity, let's assume the secret is correctly formatted for now, or just provide a helper.

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(fmt.Sprintf("%-32s", secret)[:32]))
	if err != nil {
		return nil, appError.ErrInternal.Wrap(err, "failed at generate SymmetricKeyFromBytes")
	}

	return &pasetoManager{
		symmetricKey:   key,
		expiryDuration: expiryDuration,
	}, nil
}

func (p *pasetoManager) Generate(_ context.Context, tokenData *entity.TokenData) (string, time.Time, error) {
	expirationTime := time.Now().Add(p.expiryDuration)

	token := paseto.NewToken()
	token.SetExpiration(expirationTime)
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	token.SetString("id", tokenData.Id)
	token.SetString("username", tokenData.Username)

	return token.V4Encrypt(p.symmetricKey, nil), expirationTime, nil
}

func (p *pasetoManager) Validate(_ context.Context, tokenStr string) (*entity.TokenData, error) {
	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(p.symmetricKey, tokenStr, nil)
	if err != nil {
		return nil, appError.ErrUnauthorized.Wrap(err, "invalid token")
	}

	id, err := token.GetString("id")
	if err != nil {
		return nil, appError.ErrUnauthorized.Wrap(err, "token missing id claim")
	}

	username, err := token.GetString("username")
	if err != nil {
		return nil, appError.ErrUnauthorized.Wrap(err, "token missing username claim")
	}

	return &entity.TokenData{
		Id:       id,
		Username: username,
	}, nil
}
