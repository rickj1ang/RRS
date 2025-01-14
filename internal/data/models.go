package data

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Records RecordModel
	Users   UserModel
	Tokens  TokenRedis
	Notify  NotifyQue
}

func NewModels(cl *mongo.Client, rd *redis.Client, conn *amqp.Connection) Models {
	return Models{
		Records: RecordModel{CL: cl},
		Users:   UserModel{CL: cl},
		Tokens:  TokenRedis{CL: rd},
		Notify:  NotifyQue{Conn: conn},
	}
}
