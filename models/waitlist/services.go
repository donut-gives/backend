package waitlist

import (
	"context"
	"donutBackend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"errors"
)

var waitlistCollection = new(mongo.Collection)

func init() {
	waitlistCollection = db.Get().Collection("waitlist")
	_, err := waitlistCollection.Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		return
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

func GetCount() (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	count, err := waitlistCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return count, nil
}
