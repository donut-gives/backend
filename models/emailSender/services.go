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

func Find(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult EmailSender
	err := emailSenderCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetToken(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult EmailSender
	err := emailSenderCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return findResult.Token, nil
}

func UpdateToken(email string, token string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOneAndUpdate()
	res := emailSenderCollection.FindOneAndUpdate(
		ctx,
		bson.D{{Key: "email", Value: email}},
		bson.D{{Key: "token", Value: token}},
		opts,
	)
	if res == nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		return false
	}
	return true
}

func GetEmail() (string, error) {
	//get all documents
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find()
	var findResult []EmailSender
	cur, err := emailSenderCollection.Find(
		ctx,
		bson.D{},
		opts,
	)
	if err != nil {
		return "", err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var elem EmailSender
		err := cur.Decode(&elem)
		if err != nil {
			return "", err
		}
		findResult = append(findResult, elem)
	}
	if err := cur.Err(); err != nil {
		return "", err
	}
	return findResult[0].Email, nil
}
