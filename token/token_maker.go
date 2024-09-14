package token

import (
	"crypto/rsa"
	"crypto/subtle"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/DEVunderdog/user-management-gRPC/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrWrongAudience    = errors.New("wrong token audience")
	ErrWrongIssuer      = errors.New("wrong token issuer")
	ErrInvalidTokenType = errors.New("invalid token type")
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"token_type"`
	Nonce  string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(privateKey *rsa.PrivateKey, userID uint, email string, config utils.Config) (string, string, error){

	if userID == 0 || email == "" {
		return "", "", fmt.Errorf("invalid input: userID or email is empty")
	}

	nonce := uuid.New().String()

	claims := &Claims{
		UserID: userID,
		Email: email,
		Type: AccessTokenType,
		Nonce: nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: config.Issuer,
			Subject: strconv.Itoa(int(userID)),
			Audience: jwt.ClaimStrings{config.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15*time.Minute + time.Duration(rand.Intn(60))*time.Second)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID: uuid.New().String(),
		},
	}

	refreshClaims := &Claims{
		UserID: userID,
		Email: email,
		Type: RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: config.Issuer,
			Subject: strconv.Itoa(int(userID)),
			Audience: jwt.ClaimStrings{config.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID: uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func ValidateToken(tokenString string, publicKey *rsa.PublicKey, config utils.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error){
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired){
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Type != AccessTokenType && claims.Type != RefreshTokenType{
		return nil, ErrInvalidToken
	}

	if !verifyAudience(claims.Audience, config.Audience) {
		return nil, ErrWrongAudience
	}

	if !verifyIssuer(claims.Issuer, config.Issuer) {
		return nil, ErrWrongIssuer
	}

	if claims.Issuer != config.Issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	if claims.Subject != strconv.Itoa(int(claims.UserID)){
		return nil, fmt.Errorf("subject does not match UserID")
	}
	
	if claims.Type == RefreshTokenType{
		if claims.NotBefore.After(time.Now().Add(15 * time.Minute)){
			return nil, errors.New("token not yet valid")
		}
	}

	if claims.NotBefore.After(time.Now()) {
		return nil, errors.New("token not yet valid")
	}

	return claims, nil
}


func verifyAudience(tokenAud jwt.ClaimStrings, expectedAud string) bool {
	for _, aud := range tokenAud{
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error){
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

