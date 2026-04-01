package auth

import (
	"context"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"gaman-microservice/user-service/internal/domain/entity"
	"gaman-microservice/user-service/internal/port/out"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
)

type pasetoManager struct {
	privateKey     paseto.V4AsymmetricSecretKey
	publicKey      paseto.V4AsymmetricPublicKey
	expiryDuration time.Duration
}

func NewPasetoManager(privateKey string, publicKey string, expiryDuration time.Duration) (out.TokenManager, error) {
	tm := &pasetoManager{
		expiryDuration: expiryDuration,
	}

	passetoPk, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKey)
	if err != nil {
		return nil, appError.ErrInternal.Wrap(err, "failed at generate V4AsymmetricPrivateKeyFromHex")
	}
	tm.privateKey = passetoPk

	if publicKey != "" {
		tm.publicKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publicKey)
		if err != nil {
			return nil, appError.ErrInternal.Wrap(err, "failed at generate V4AsymmetricPublicKeyFromHex")
		}
	} else {
		log.Info().Msg("public key is empty, generate new one from private key")
		tm.publicKey = passetoPk.Public()
	}

	return tm, nil
}

func (p *pasetoManager) Generate(_ context.Context, tokenData *entity.TokenData) (string, time.Time, error) {
	expirationTime := time.Now().Add(p.expiryDuration)

	token := paseto.NewToken()
	token.SetExpiration(expirationTime)
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())

	token.SetString("id", tokenData.Id)
	token.SetString("username", tokenData.Username)

	return token.V4Sign(p.privateKey, nil), expirationTime, nil
}

func (p *pasetoManager) Validate(_ context.Context, tokenStr string) (*entity.TokenData, error) {
	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(p.publicKey, tokenStr, nil)
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
