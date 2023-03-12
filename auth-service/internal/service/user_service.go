package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"auth-service/api/pkg/model"
)

type userService struct {
	UserRepository model.UserRepository
}

type USConfig struct {
	UserRepository model.UserRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

func (s *userService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("Unable to signup user %v", u.UID)
		return model.InternalError()
	}
	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	return nil
}
