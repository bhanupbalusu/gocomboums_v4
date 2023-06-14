package service

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/chacha20poly1305"
)

type Claims struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf"`
}

func GenerateKey() []byte {
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error()) // handle error in production code
	}
	return key
}

func StoreKeyInFile(key []byte) {
	err := os.WriteFile("../../pkg/utils/keyfile", key, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func RetrieveKeyInFile() []byte {
	key, err := os.ReadFile("../../pkg/utils/keyfile")
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func GeneratePasetoToken(claims *Claims, key []byte) (string, error) {
	// Ensure the key is the correct size
	if len(key) != chacha20poly1305.KeySize {
		return "", errors.New("incorrect key size")
	}

	// Initialize a Paseto V2 object
	v2 := paseto.NewV2()

	// Set the expiration time in the token
	now := time.Now()
	claims.ExpiresAt = now.Add(15 * time.Minute).Unix() // Expires in 24 hours
	claims.IssuedAt = now.Unix()
	claims.NotBefore = now.Unix()

	// Serialize the claims to JSON
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Encrypt and encode the token
	token, err := v2.Encrypt(key, jsonClaims, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyAndExtractClaims(token string, key []byte) (*Claims, error) {
	// Ensure the key is the correct size
	if len(key) != chacha20poly1305.KeySize {
		return nil, errors.New("incorrect key size")
	}

	// Initialize a Paseto V2 object
	v2 := paseto.NewV2()

	var claims Claims
	err := v2.Decrypt(token, key, &claims, nil)
	if err != nil {
		return nil, err
	}

	// Validate the claims
	now := time.Now()
	if now.Unix() < claims.NotBefore {
		return nil, errors.New("token is not valid yet")
	}
	if now.Unix() > claims.ExpiresAt {
		return nil, errors.New("token is expired")
	}

	return &claims, nil
}
