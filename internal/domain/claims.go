package domain

import "time"

type Claims struct {
	UserID    int64
	Login     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
