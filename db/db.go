package db

import (
	"context"
	"log"

	"github.com/kaepa3/healthplanet"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegisterOption struct {
	Url string
}

const (
	dbName         = "hellsee"
	collectionName = "healthplanet"
)

type Record struct {
	Id      primitive.ObjectID `bson:"_id"`
	UserID  string             `bson:"userid"`
	Date    string             `bson:"date"`
	KeyData string             `bson:"keydata"`
	Model   string             `bson:"model"`
	Tag     string             `bson:"tag"`
}

func createRecord(id string, d healthplanet.Data) Record {
	rec := Record{
		Id:      primitive.NewObjectID(),
		UserID:  id,
		Date:    d.Date,
		KeyData: d.KeyData,
		Model:   d.Model,
		Tag:     d.Tag,
	}
	return rec
}

func RegisterDB(data *healthplanet.JsonResponce, opt *RegisterOption, clientID string, ctx context.Context) error {
	log.Println("insertdb:", opt.Url)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(opt.Url))
	if err != nil {
		return err
	}
	coll := client.Database(dbName).Collection(collectionName)
	for _, v := range data.Data {
		rec := createRecord(clientID, v)
		log.Println("yarude:", v)
		filter := bson.D{
			{Key: "userid", Value: rec.UserID},
			{Key: "date", Value: rec.Date},
			{Key: "tag", Value: rec.Tag},
		}
		count, err := coll.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}
		if count == 0 {
			_, err := coll.InsertOne(ctx, rec)
			if err != nil {
				return err
			}
		}
		log.Println("count:", count)
	}
	return nil
}
