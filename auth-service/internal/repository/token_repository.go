package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"auth-service/api/pkg/model"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return model.InternalError()
	}
	return nil
}

func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return model.InternalError()
	}

	return nil
}