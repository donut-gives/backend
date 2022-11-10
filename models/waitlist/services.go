package waitlist

import (
	"context"
	"donutBackend/db"
	. "donutBackend/logger"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var waitlistCollection = new(mongo.Collection)

func init() {
	waitlistCollection = db.Get().Collection("waitlist")

	// Create a unique index on the email field
	indexView := waitlistCollection.Indexes()
	mod := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := indexView.CreateOne(context.Background(), mod)
	if err != nil {
		Logger.Errorf("Error creating email index in waitlist collection: %v", err)
	}

}

func Insert(user WaitlistedUser) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := waitlistCollection.InsertOne(ctx, user)

	if err != nil {
		if(mongo.IsDuplicateKeyError(err)){
			return nil, errors.New("Already In Waitlist")
		}
		return nil, err
	}
	// fmt.Println("Inserted a single user: ", result.InsertedID)
	id := result.InsertedID
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
}
