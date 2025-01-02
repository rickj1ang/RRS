package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

const ScopeAuthentication = "authentication"

type Token struct {
	Plaintext string    `json:"token" bson:"token"`
	Hash      []byte    `json:"-" bson:"hash"`
	UserEmail string    `json:"-" bson:"email"`
	Expiry    time.Time `json:"expiry" bson:"expiry"`
	Scope     string    `json:"-" bson:"scope"`
}

func generateToken(userEmail string, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserEmail: userEmail,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	CL *mongo.Client
}

func (t TokenModel) New(userEmail string, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userEmail, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err

}

func (t TokenModel) Insert(token *Token) error {
	coll := connectRRStokens(t)
	_, err := coll.InsertOne(context.TODO(), token)
	if err != nil {
		return err
	}
	return nil
}
