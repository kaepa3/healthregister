package db

import (
	"context"
	"log"

	"github.com/kaepa3/healthplanet"
	"go.mongodb.org/mongo-driver/bson"
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

func RegisterDB(data *healthplanet.JsonResponce, opt *RegisterOption, ctx context.Context) error {
	log.Println("connect:", opt.Url)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(opt.Url))
	if err != nil {
		return err
	}
	coll := client.Database(dbName).Collection(collectionName)
	for _, v := range data.Data {
		filter := bson.D{
			{Key: "date", Value: v.Date},
			{Key: "tag", Value: v.Tag},
		}
		count, err := coll.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}
		if count == 0 {
			_, err := coll.InsertOne(ctx, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
