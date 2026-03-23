package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cmpnion.space/internal/domain"
)

type Service struct {
	users    UserRepository
	tokens   TokenRepository
	hasher   PasswordHasher
	jwt      TokenManager
	tokenTTL time.Duration
}

func NewService(
	users UserRepository,
	tokens TokenRepository,
	hasher PasswordHasher,
	jwt TokenManager,
	tokenTTL time.Duration,
) *Service {
	return &Service{
		users:    users,
		tokens:   tokens,
		hasher:   hasher,
		jwt:      jwt,
		tokenTTL: tokenTTL,
	}
}

func (s *Service) SignUp(ctx context.Context, login, password string) (string, error) {
	login = strings.TrimSpace(login)
	if login == "" || password == "" {
		return "", domain.ErrValidation
	}

	if _, err := s.users.GetByLogin(ctx, login); err == nil {
		return "", domain.ErrLoginTaken
	} else if !errors.Is(err, domain.ErrNotFound) {
		return "", fmt.Errorf("check login: %w", err)
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := domain.User{
		Login:        login,
		PasswordHash: hash,
	}
	created, err := s.users.Create(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrLoginTaken) {
			return "", domain.ErrLoginTaken
		}
		return "", fmt.Errorf("create user: %w", err)
	}

	token, err := s.jwt.Generate(created.ID, created.Login, s.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	if err := s.tokens.Store(ctx, created.ID, token, s.tokenTTL); err != nil {
		return "", fmt.Errorf("store token: %w", err)
	}

	return token, nil
}

func (s *Service) SignIn(ctx context.Context, login, password string) (string, error) {
	login = strings.TrimSpace(login)
	if login == "" || password == "" {
		return "", domain.ErrValidation
	}

	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", fmt.Errorf("get user: %w", err)
	}

	if !s.hasher.Compare(user.PasswordHash, password) {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(user.ID, user.Login, s.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	if err := s.tokens.Store(ctx, user.ID, token, s.tokenTTL); err != nil {
		return "", fmt.Errorf("store token: %w", err)
	}

	return token, nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (int64, error) {
	userID, err := s.jwt.Validate(token)
	if err != nil {
		return 0, fmt.Errorf("validate token: %w", err)
	}

	return userID, nil
}

func (s *Service) GetMe(ctx context.Context, userID int64) (UserView, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return UserView{}, fmt.Errorf("get user: %w", err)
	}

	return UserView{
		ID:    user.ID,
		Login: user.Login,
	}, nil
}
