package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

var (
	ErrInvalidToken     = errors.New("token is invalid")
	ErrExpiredToken     = errors.New("token has expired")
	ErrWrongAudience    = errors.New("wrong token audience")
	ErrWrongIssuer      = errors.New("wrong token issuer")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"username"`
	Role      string    `json:"role"`
	Type      string    `json:"token_type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	jwt.RegisteredClaims
}

func NewPayload(email string, role string, duration time.Duration, issuer string, audience string, notBefore time.Duration, tokenType string) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:    tokenID,
		Email: email,
		Role:  role,
		Type: tokenType,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			ID:        uuid.New().String(),
		},
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
