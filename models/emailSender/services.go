package emailsender

import (
	"context"
	"donutBackend/db"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	//"errors"
)

var emailSenderCollection = new(mongo.Collection)

func init() {
	emailSenderCollection = db.Get().Collection("email_sender")
	_, err := emailSenderCollection.Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		return
	}
}

func InsertOrUpdateOne(sender EmailSender) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{{Key: "email", Value: sender.Email}}
	update := bson.D{{Key: "$set", Value: sender}}
	opts := options.Update().SetUpsert(true)
	res, err := emailSenderCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	if res.UpsertedCount == 1 {
		return res.UpsertedID, nil
	}
	return nil, nil
}

func SetDeactivated(email string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: "FALSE"}}}}
	_, err := emailSenderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}