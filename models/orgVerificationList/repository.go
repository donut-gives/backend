package orgVerification

import (
	"context"
	"donutBackend/db"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var organizationCollection = new(mongo.Collection)

func init() {
	organizationCollection = db.Get().Collection("organizationsList")
}

func Insert(org *Organization) (interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult bson.M

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: org.Email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			org.Status = "PENDING"
			result, err := organizationCollection.InsertOne(ctx, org)
			if err != nil {
				return nil, err
			}

			id := result.InsertedID
			stringId := id.(primitive.ObjectID).Hex()
			return stringId, nil
		}
		return nil, err
	}

	return nil, errors.New("Email already exists")
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
	if(findResult["verified"] != "PENDING"){
		return nil, errors.New("Organization already "+findResult["verified"].(string))
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: "VERIFIED"}}}}
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
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: "REJECTED"}}}}
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

func Find(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult bson.M
	err := organizationCollection.FindOne(
		ctx,
		bson.D{
			{Key: "email", Value: email},
			{Key: "verified", Value: "true"},
		},
		opts,
	).Decode(&findResult)
	if err != nil {
		
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	
	return true, nil
}

