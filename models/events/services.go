package events

import (
	"context"
	"donutBackend/db"

	//. "donutBackend/logger"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var eventsCollection = new(mongo.Collection)

func init() {
	eventsCollection = db.Get().Collection("events")

	
}

func GetEventById(id string) (Event, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOne()

	eventId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Event{}, err
	}
	var findResult Event
	err = eventsCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: eventId}},
		opts,
	).Decode(&findResult)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return Event{}, errors.New("No event found")
		}
		return Event{}, err
	}
	return findResult, nil
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
		if(mongo.IsDuplicateKeyError(err)){
			return nil, errors.New("Event already exists")
		}
		return nil, err
	}
	return result.InsertedID, nil
}

func DeleteEvent(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := eventsCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	if(result.DeletedCount == 0){
		return errors.New("Event not found")
	}
	return nil
}

func InsertEvent(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := eventsCollection.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	if(result.InsertedID == nil){
		return errors.New("Event not inserted")
	}
	return nil
}
