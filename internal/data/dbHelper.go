package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func structToDoc(v any) (bson.D, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	var doc bson.D
	err = bson.Unmarshal(data, &doc)
	return doc, err
}

func connectRRSrecords(r RecordModel) *mongo.Collection {
	return r.CL.Database("RRS").Collection("records")
}

func connectRRSusers(u UserModel) *mongo.Collection {
	return u.CL.Database("RRS").Collection("users")
}

func connectRRStokens(t TokenModel) *mongo.Collection {
	return t.CL.Database("RRS").Collection("tokens")
}
