package jwt

import (
	"fmt"
	"time"

	"cmpnion.space/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
}

func NewManager(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

func (m *Manager) Generate(userID int64, login string, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":   userID,
		"login": login,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Validate(token string) (int64, error) {
	claims, err := m.Parse(token)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func (m *Manager) Parse(token string) (*domain.Claims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, domain.ErrUnauthorized
	}

	mapClaims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	sub, ok := mapClaims["sub"].(float64)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	login, _ := mapClaims["login"].(string)

	iatFloat, _ := mapClaims["iat"].(float64)
	expFloat, _ := mapClaims["exp"].(float64)

	claims := &domain.Claims{
		UserID:    int64(sub),
		Login:     login,
		IssuedAt:  time.Unix(int64(iatFloat), 0),
		ExpiresAt: time.Unix(int64(expFloat), 0),
	}

	return claims, nil
}
