package data

import (
	"context"
	"time"

	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecordModel struct {
	CL *mongo.Client
}

// impl MarshalJSON()([]byte, err) method for futhur customize the json
type Record struct {
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	Title       string    `json:"title" bson:"title"`
	Writer      string    `json:"writer" bson:"writer,omitempty"`
	TotalPages  uint16    `json:"total_pages" bson:"total_pages"`
	CurrentPage uint16    `json:"current_page" bson:"current_page"`
	Progress    float32   `json:"progress" bson:"progress"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Genres      []string  `json:"genres,omitempty" bson:"genres,omitempty"`
}

func ValidateRecord(v *validator.Validator, record *Record) {
	v.Check(record.Title != "", "title", "must be provided")
	v.Check(len(record.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(record.TotalPages > 0, "page", "must have a positive total page")
	v.Check(record.CurrentPage >= 0, "page", "must have a positive current page")
	v.Check(record.CurrentPage <= record.TotalPages, "pages", "can not read more than total pages")

	v.Check(len(record.Genres) <= 3, "genres", "must not more than 3 genres")
	v.Check(validator.Unique(record.Genres), "genres", "genres can not be dupicate")
}

func (r RecordModel) Insert(record *Record) (string, error) {
	coll := r.CL.Database("RRS").Collection("records")

	res, err := coll.InsertOne(context.TODO(), record)
	if err != nil {
		return "Fail to insert", err
	}
	stringID := res.InsertedID.(primitive.ObjectID).Hex()

	return stringID, nil
}

// give an arbitary key-value pair return the record sruct and err
// key must be string val can be any
func (r RecordModel) Get(key string, value any) (*Record, error) {
	coll := r.CL.Database("RRS").Collection("records")

	filter := bson.D{{key, value}}

	var result Record
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r RecordModel) Update(record *Record) error {
	return nil
}

func (r RecordModel) Delete(id int64) error {
	return nil
}
