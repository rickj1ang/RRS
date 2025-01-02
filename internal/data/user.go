package data

import (
	"context"
	"errors"
	"time"

	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type UserModel struct {
	CL *mongo.Client
}

type UserLevel int8

// return string of user level
func (u UserLevel) String() string {
	switch {
	// regular user
	case u == 1:
		return "normal"
	// un activate user
	case u == 0:
		return "die"
	// manager can propose a post to delete a recommand_book
	case u == 2:
		return "manager"
	// me
	case u == 3:
		return "lord"
	// baned user, do something really bad
	case u == -1:
		return "devil"
	// do not know who
	default:
		return "unknow"
	}
}

type User struct {
	//email
	Email string `json:"email" bson:"email"`
	//password
	Password password `json:"password" bson:"password"`
	//nickname
	Name string `json:"name" bson:"name"`
	//createAt
	CreateAt time.Time `json:"create_at" bson:"create_at"`
	//records
	Records []primitive.ObjectID `json:"records,omitempty" bson:"records,omitempty"`
	//recommand book
	RecommandBook primitive.ObjectID `json:"recommand_book,omitempty" bson:"recommand_book,omitempty"`
	//Level
	Level UserLevel `json:"user_level" bson:"user_level"`
}

// if we do not have palintext It will be nil, because it is a pointer
type password struct {
	Plaintext *string `json:"plaintext" bson:"plaintext"`
	Hash      []byte  `json:"hash" bson:"hash"`
}

func (p *password) Set(plaintestPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintestPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintestPassword
	p.Hash = hash

	return nil
}

func (p *password) Match(plaintestPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintestPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must provide email")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a ture email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password can not be empty")
	v.Check(len(password) >= 8, "password", "password must longer than 8 bytes")
	v.Check(len(password) <= 66, "password", "password must less than 66 bytes")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "every user musr have a name")
	v.Check(len(user.Name) < 100, "name", "name should be short and stong")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	// a very important situation
	if user.Password.Hash == nil {
		panic("miss a user's password hash!")
	}
}

func (u UserModel) Insert(user *User) (string, error) {
	if u.IsEmailExist(user.Email) {
		return "", ErrDuplicateEmail
	}
	user.CreateAt = time.Now()
	user.Level = 0
	coll := connectRRSusers(u)
	res, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}
	stringID := res.InsertedID.(primitive.ObjectID).Hex()

	return stringID, nil
}

func (u UserModel) IsEmailExist(email string) bool {
	coll := connectRRSusers(u)

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	res := coll.FindOne(context.TODO(), filter)
	switch {
	case errors.Is(res.Err(), mongo.ErrNoDocuments):
		return false
	default:
		return true
	}
}

func (u UserModel) Get(key string, value any) (*User, error) {
	coll := connectRRSusers(u)
	filter := bson.D{primitive.E{Key: key, Value: value}}
	var user User
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserModel) Update(key string, value any, user *User) error {
	coll := connectRRSusers(u)
	doc, err := structToDoc(user)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: key, Value: value}}
	update := bson.D{primitive.E{Key: "$set", Value: doc}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)

	return err
}

func (u UserModel) Delete(id primitive.ObjectID) error {
	coll := connectRRSusers(u)
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}
