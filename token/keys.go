package token

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"time"

	database "github.com/DEVunderdog/user-management-gRPC/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/pbkdf2"
)

type JWTKeyResponse struct {
	ID int64
	PublicKey string
	PrivateKey string
	ExpiresAt time.Time
}

func InitializeJWTKeys(passphrase string, store database.Store, ctx context.Context) error {
	
	count, err := store.CountJWTKeys(ctx)

	if err != nil {
		return err
	}

	if count == 0 {
		return generateAndStoreKeys(passphrase, store, ctx)
	}

	return nil
}

func generateAndStoreKeys(passphrase string, store database.Store, ctx context.Context) error {
	privateKey, publicKey, err := generateRSAKeys()

	if err != nil {
		return err
	}

	encryptedPrivateKey, err := encryptPrivateKey(privateKey, []byte(passphrase))
	if err != nil {
		return err
	}

	publicKeyPEM, err := encodePublicKey(publicKey)
	if err != nil {
		return err
	}

	jwtKeyParams := database.CreateJWTKeyParams{
		PublicKey: base64.StdEncoding.EncodeToString(publicKeyPEM),
		PrivateKey: base64.StdEncoding.EncodeToString(encryptedPrivateKey),
		Algorithm: "RS256",
		IsActive: pgtype.Bool{
			Bool: true,
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time: time.Now().AddDate(0, 6, 0),
			Valid: true,
		},
	}

	_, err = store.CreateJWTKey(ctx, jwtKeyParams)
	if err != nil {
		return err
	}

	return nil

}

func GetActiveJWTKey(ctx context.Context, isActive bool, store database.Store) (*JWTKeyResponse, error){
	
	active := pgtype.Bool{
		Bool: isActive,
		Valid: true,
	}
	
	jwtStruct, err := store.GetActiveKey(ctx, active)
	if err != nil {
		return nil, err
	}

	data := &JWTKeyResponse{
		ID: jwtStruct.ID,
		PublicKey: jwtStruct.PublicKey,
		PrivateKey: jwtStruct.PrivateKey,
		ExpiresAt: jwtStruct.ExpiresAt.Time,
	}


	return data, nil
}

func GetPrivateKey( key string, passphrase []byte) (*rsa.PrivateKey, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return decryptPrivateKey(encryptedKey, passphrase)
}

func GetPublicKey(key string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyPEM)

	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type){
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("not an RSA public key")
	}
}

func generateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func encryptPrivateKey(privateKey *rsa.PrivateKey, passphrase []byte) ([]byte, error) {

	block := &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	blockBytes := pem.EncodeToMemory(block)
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key(passphrase, salt, 100000, 32, sha256.New)

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, nonce, blockBytes, nil)

	result := make([]byte, len(salt)+len(nonce)+len(cipherText))
	copy(result, salt)
	copy(result[len(salt):], nonce)
	copy(result[len(salt)+len(nonce):], cipherText)

	return result, nil
}

func decryptPrivateKey(encryptedKey []byte, passphrase []byte) (*rsa.PrivateKey, error) {
	if len(encryptedKey) < 28 {
		return nil, errors.New("invalid encrypted key format")
	}

	salt := encryptedKey[:16]
	nonce := encryptedKey[16:28]
	cipherText := encryptedKey[28:]

	key := pbkdf2.Key(passphrase, salt, 100000, 32, sha256.New)

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	pemBytes, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM Block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func encodePublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)

	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}