package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"silent/src/internal/domain"
)

type AiApi struct {
	Hash     string `gorm:"column:hash;primaryKey;type:text"`
	Requests int    `gorm:"column:requests;type:int;default:0"`
}

func NewAiApiModel(hash string) (*AiApi, error) {
	if hash == "" {
		return nil, errors.New("hash must not be empty")
	}
	return &AiApi{Hash: hash, Requests: 0}, nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (a *AiApi) ToDomain(token string) (*domain.AiApi, error) {
	if token == "" {
		return nil, errors.New("token must not be empty")
	}
	return &domain.AiApi{Token: token, Requests: a.Requests}, nil
}
