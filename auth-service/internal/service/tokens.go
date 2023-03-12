package service

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"time"

	"auth-service/api/pkg/model"
)

type IDTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.RegisteredClaims
}

func GenerateIDToken(u *model.User, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := jwt.NewNumericDate(time.Now())
	tokenExp := jwt.NewNumericDate(unixTime.Add(time.Duration(exp) * time.Second)) // 15 minutes from current time

	claims := IDTokenCustomClaims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: tokenExp,
			IssuedAt:  unixTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return ss, nil
}

type RefreshToken struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}

type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

type idTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.RegisteredClaims
}

func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)
	tokenID, err := uuid.NewRandom()

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &RefreshToken{
		SS:        ss,
		ID:        tokenID.String(),
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

func validateIDToken(tokenString string, key *rsa.PublicKey) (*idTokenCustomClaims, error) {
	claims := &idTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*idTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}
