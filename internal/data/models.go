package data

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Records RecordModel
	Users   UserModel
	Tokens  TokenRedis
}

func NewModels(cl *mongo.Client, rd *redis.Client) Models {
	return Models{
		Records: RecordModel{CL: cl},
		Users:   UserModel{CL: cl},
		Tokens:  TokenRedis{CL: rd},
	}
}
