package users

import (
	"context"
	"donutBackend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var usersCollection = new(mongo.Collection)

func init() {
	usersCollection = db.Get().Collection("users")
}

// Insert : Create a new user
func Insert(user *GoogleUser) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult bson.M

	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: user.Email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			result, err := usersCollection.InsertOne(ctx, user)
			if err != nil {
				//log.Fatal(err)
				return nil, err
			}
			//fmt.Println("Inserted a single user: ", result.InsertedID)
			id := result.InsertedID
			stringId := id.(primitive.ObjectID).Hex()
			return stringId, nil
		}
		return nil, err
	}

	id := findResult["_id"]
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
}
