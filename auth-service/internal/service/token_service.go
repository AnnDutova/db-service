package service

import (
	"context"
	"crypto/rsa"
	"log"

	"auth-service/api/pkg/model"
)

type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivateKey            *rsa.PrivateKey
	PublicKey             *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		PrivateKey:            c.PrivateKey,
		PublicKey:             c.PublicKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
		TokenRepository:       c.TokenRepository,
	}
}

func (s *tokenService) NewToken(ctx context.Context, u *model.User, prevTokenID string) (*model.Token, error) {
	idToken, err := generateIDToken(u, s.PrivateKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	if err = s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	if prevTokenID != "" {
		if err = s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
			return nil, model.InternalError()
		}
	}

	return &model.Token{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
