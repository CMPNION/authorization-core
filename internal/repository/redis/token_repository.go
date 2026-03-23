package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client    *goredis.Client
	keyPrefix string
}

func NewTokenRepository(client *goredis.Client) *TokenRepository {
	return &TokenRepository{
		client:    client,
		keyPrefix: "auth:token",
	}
}

func (r *TokenRepository) Store(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s", r.keyPrefix, token)
	if err := r.client.Set(ctx, key, userID, ttl).Err(); err != nil {
		return fmt.Errorf("redis set token: %w", err)
	}
	return nil
}
