package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenRedis struct {
	CL *redis.Client
}

type PassHash struct {
	plaintext string
	hash      [32]byte
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func generateHash() (*PassHash, error) {
	randomBytes := make([]byte, 16)
	token := &PassHash{}

	_, err := rand.Read(randomBytes)
	if err != nil {
		return token, err
	}

	token.plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	token.hash = sha256.Sum256([]byte(token.plaintext))
	return token, nil
}

func (t TokenRedis) GiveToken(id primitive.ObjectID) (*string, error) {
	token, err := generateHash()
	if err != nil {
		return nil, err
	}

	err = t.CL.Set(context.TODO(), string(token.hash[:]), id.Hex(), 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &token.plaintext, nil
}

func (t TokenRedis) GetIdByToken(plaintext string) (primitive.ObjectID, error) {
	hash := sha256.Sum256([]byte(plaintext))
	idStr, err := t.CL.Get(context.TODO(), string(hash[:])).Result()
	if err != nil {
		return primitive.NilObjectID, err
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return id, nil
}
