package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword string, password string) bool
}

type SHA256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) *SHA256Hasher {
	return &SHA256Hasher{salt: salt}
}

func (h *SHA256Hasher) Hash(password string) ([]byte, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return nil, err
	}

	return hash.Sum([]byte(h.salt)), nil
}

func (h *SHA256Hasher) Verify(hashedPassword string, password string) bool {
	expectedHash, err := h.Hash(password)
	if err != nil {
		return false
	}

	return hmac.Equal([]byte(hashedPassword), expectedHash)
}
