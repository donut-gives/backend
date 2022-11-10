package admin

import (
	"context"
	"donutBackend/db"
	. "donutBackend/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var adminCollection = new(mongo.Collection)

func init() {
	adminCollection = db.Get().Collection("admin")

	// Create a unique index on the email field
	indexView := adminCollection.Indexes()
	mod := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := indexView.CreateOne(context.Background(), mod)
	if err != nil {
		Logger.Errorf("Error creating email index in admin collection: %v", err)
	}
}

func Find(email string) (bool,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult bson.M
	err := adminCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}
