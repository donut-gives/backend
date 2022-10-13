package orgVerification

import (
	"context"
	"donutBackend/db"
	"errors"
	"time"

	. "donutBackend/logger"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var organizationCollection = new(mongo.Collection)

func init() {
	organizationCollection = db.Get().Collection("organizationsList")

	// Create a unique index on the email field
	indexView := organizationCollection.Indexes()
	mod := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := indexView.CreateOne(context.Background(), mod)
	if err != nil {
		Logger.Errorf("Error creating email index in organizationsList collection: %v", err)
	}
}

func Insert(org *Organization) (interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	(*org).Status = "PENDING"

	insertResult, err := organizationCollection.InsertOne(ctx, org)
	if err != nil {
		if(mongo.IsDuplicateKeyError(err)){
			return nil, errors.New("Organization already exists")
		}
		return nil, err
	}

	return insertResult.InsertedID, nil
}

func Get(email string) (*Organization, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult Organization

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		if(err == mongo.ErrNoDocuments){
			return nil, errors.New("Organization does not exist")
		}
		return nil, err
	}

	return &findResult, nil
}

func Verify(email string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	option := options.FindOne()
	var findResult bson.M

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		option,
	).Decode(&findResult)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Organization does not exist")
		}
		return nil, err
	}
	if(findResult["status"] != "PENDING"){
		return nil, errors.New("Organization already "+findResult["status"].(string))
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "VERIFIED"}}}}
	var updatedDocument Organization
	err = organizationCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedDocument)

	if err != nil {
		return nil, err
	}

	return &updatedDocument, nil
}

func Reject(email string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	option := options.FindOne()
	var findResult bson.M

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		option,
	).Decode(&findResult)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Organization does not exist")
		}
		return nil, err
	}
	if(findResult["status"] != "PENDING"){
		return nil, errors.New("Organization already "+findResult["status"].(string))
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: "REJECTED"}}}}
	var updatedDocument Organization
	err = organizationCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&updatedDocument)

	if err != nil {
		return nil, err
	}

	return &updatedDocument, nil
}

func Find(email string) (Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{
			{Key: "email", Value: email},
			{Key: "status", Value: "VERIFIED"},
		},
		opts,
	).Decode(&findResult)
	if err != nil {
		
		if err == mongo.ErrNoDocuments {
			return Organization{}, errors.New("No such organization found")
		}
		return Organization{}, err
	}
	
	return findResult, nil
}

