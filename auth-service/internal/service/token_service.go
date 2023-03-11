package service

import (
	"context"
	"crypto/rsa"
	"log"

	"auth-service/api/pkg/model"
)

type TokenService struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

type TSConfig struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &TokenService{
		PrivateKey:    c.PrivateKey,
		PublicKey:     c.PublicKey,
		RefreshSecret: c.RefreshSecret,
	}
}

func (s *TokenService) NewToken(ctx context.Context, u *model.User, prevTokenID string) (*model.Token, error) {
	idToken, err := generateIDToken(u, s.PrivateKey)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	// TODO: store refresh tokens by calling TokenRepository methods

	return &model.Token{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
