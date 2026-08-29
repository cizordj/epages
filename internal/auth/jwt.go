package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Raw       string
	ExpiresAt time.Time
}

func New(raw string) (*Token, error) {
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(raw, claims)
	if err != nil {
		return nil, err
	}

	var expiresAt time.Time
	if exp, err := claims.GetExpirationTime(); err == nil && exp != nil {
		expiresAt = exp.Time
	}

	return &Token{
		Raw:       raw,
		ExpiresAt: expiresAt,
	}, nil
}

func (t *Token) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(t.ExpiresAt)
}
