package token

import (
	"time"
)

type Maker interface {
	GenerateToken(duration, noteBefore time.Duration, email, role, issuer, audience, tokenType string) (*string, *Payload, error)

	VerifyToken(token string, audience, issuer string) (*Payload, error)
}