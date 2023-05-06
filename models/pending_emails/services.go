package pending_emails

import (
	"context"
	"errors"
	"github.com/donut-gives/backend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var pendingEmailCollection = new(mongo.Collection)

func init() {
	pendingEmailCollection = db.Get().Collection("pending_email")
	_, err := pendingEmailCollection.Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		return
	}
}

func Insert(pending PendingEmail) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := pendingEmailCollection.InsertOne(ctx, pending)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("Already In Waitlist")
		}
		return nil, err
	}
	// fmt.Println("Inserted a single user: ", result.InsertedID)
	id := result.InsertedID
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
}
