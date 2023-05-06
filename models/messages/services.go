package messages

import (
	"context"
	"github.com/donut-gives/backend/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var messageCollection = new(mongo.Collection)

func init() {
	messageCollection = db.Get().Collection("messages")
}

func Insert(message Message) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := messageCollection.InsertOne(ctx, message)

	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	//fmt.Println("Inserted a single user: ", result.InsertedID)
	id := result.InsertedID
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
}
