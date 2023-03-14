package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
	Signin(ctx context.Context, u *User) error
}
type TokenService interface {
	NewToken(ctx context.Context, u *User, prevTokenID string) (*Token, error)
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(tokenString string) (*RefreshToken, error)
	Signout(ctx context.Context, uid uuid.UUID) error
}

type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
	DeleteUserRefreshToken(ctx context.Context, userID string) error
}
