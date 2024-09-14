package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

func HashPassword(password string) (string, error) {

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	peppered := []byte(password + string(salt))

	hash, err := bcrypt.GenerateFromPassword(peppered, bcryptCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(append(salt, hash...)), nil
}

func CheckPassword(password string, hashedPassword string) (bool, error) {

	decoded, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, err
	}

	salt, hash := decoded[:16], decoded[16:]

	peppered := []byte(password + string(salt))

	err = bcrypt.CompareHashAndPassword(hash, peppered)

	if err == bcrypt.ErrMismatchedHashAndPassword{
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}