package events

import (
	"context"
	"donutBackend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var eventsCollection = new(mongo.Collection)

func init() {
	eventsCollection = db.Get().Collection("events")
}

func GetEvents() ([]Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.Find()
	var findResult []Event
	cur, err := eventsCollection.Find(
		ctx,
		bson.D{},
		opts,
	)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var result Event
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		findResult = append(findResult, result)
	}
	return findResult, nil
}

func AddEvent(event *Event) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := eventsCollection.InsertOne(ctx, event)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func DeleteEvent(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := eventsCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	return nil
}
