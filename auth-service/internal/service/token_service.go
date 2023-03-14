package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/google/uuid"

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
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
			return nil, model.InternalError()
		}
	}

	idToken, err := GenerateIDToken(u, s.PrivateKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	if err = s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, model.InternalError()
	}

	return &model.Token{
		IDToken: model.IDToken{SS: idToken},
		RefreshToken: model.RefreshToken{
			SS:  refreshToken.SS,
			ID:  refreshToken.ID,
			UID: u.UID,
		},
	}, nil
}

func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, s.PublicKey)

	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, model.UnauthorizedError("Unable to verify user from idToken")
	}

	return claims.User, nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	if err != nil {
		log.Printf("Unable to validate or parse refreshToken - Error: %v\n", err)
		return nil, model.UnauthorizedError("Unable to verify user from refreshToken")
	}

	tokenUUID, err := uuid.Parse(claims.ID)
	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %v\n", err)
		return nil, model.UnauthorizedError("Unable to verify user from refreshToken")
	}

	return &model.RefreshToken{
		ID:  tokenUUID,
		SS:  tokenString,
		UID: claims.UID,
	}, nil
}

func (s *tokenService) Signout(ctx context.Context, uid uuid.UUID) error {
	return s.TokenRepository.DeleteUserRefreshToken(ctx, uid.String())
}
