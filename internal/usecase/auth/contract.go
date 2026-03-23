package auth

import (
	"context"
	"time"

	"cmpnion.space/internal/domain"
)

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type TokenRepository interface {
	Store(ctx context.Context, userID int64, token string, ttl time.Duration) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

type TokenManager interface {
	Generate(userID int64, login string, ttl time.Duration) (string, error)
	Validate(token string) (int64, error)
}
