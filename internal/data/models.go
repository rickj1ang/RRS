package data

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Records RecordModel
	Users   UserModel
	Tokens  TokenModel
}

func NewModels(cl *mongo.Client) Models {
	return Models{
		Records: RecordModel{CL: cl},
		Users:   UserModel{CL: cl},
		Tokens:  TokenModel{CL: cl},
	}
}
