package waitlist

import (
	"context"
	"donutBackend/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var waitlistCollection = new(mongo.Collection)

func init() {
	waitlistCollection = db.Get().Collection("waitlist")
}

func Insert(user WaitlistedUser) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := waitlistCollection.InsertOne(ctx, user)

	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	// fmt.Println("Inserted a single user: ", result.InsertedID)
	id := result.InsertedID
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
}
