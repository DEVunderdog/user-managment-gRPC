package token

import (
	"crypto/rsa"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJWTMaker(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (Maker, error) {
	return &JWTMaker{
		publicKey: publicKey,
		privateKey: privateKey,
	}, nil
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"token_type"`
	Nonce  string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

func (maker *JWTMaker) GenerateToken(duration, noteBefore time.Duration, email, role, issuer, audience, tokenType string) (*string, *Payload, error) {

	payload, err := NewPayload(email, role, duration, issuer, audience, noteBefore, tokenType)
	if err != nil {
		return nil, payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token, err := jwtToken.SignedString(maker.privateKey)

	return &token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string, audience, issuer string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, ErrInvalidToken
		}
		return maker.publicKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok || !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	if payload.Type != AccessTokenType && payload.Type != RefreshTokenType {
		return nil, ErrInvalidTokenType
	}

	if !verifyAudience(payload.Audience, audience) {
		return nil, ErrWrongAudience
	}

	if !verifyIssuer(payload.Issuer, issuer) {
		return nil, ErrWrongIssuer
	}

	if payload.Type == RefreshTokenType {
		if payload.IssuedAt.After(time.Now().Add(15 * time.Minute)) {
			return nil, errors.New("token not yet valid")
		}
	}

	return payload, nil
}

func verifyAudience(tokenAud jwt.ClaimStrings, expectedAud string) bool {
	for _, aud := range tokenAud {
		if subtle.ConstantTimeCompare([]byte(aud), []byte(expectedAud)) == 1 {
			return true
		}
	}

	return false
}

func verifyIssuer(tokenIss string, expectedIss string) bool {
	return subtle.ConstantTimeCompare([]byte(tokenIss), []byte(expectedIss)) == 1
}

func GetExpirationTime(tokenString string, publicKey *rsa.PublicKey) (*time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.ExpiresAt.Time, nil
	}

	return nil, errors.New("invalid token")
}
