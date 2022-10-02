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
			org.Verified = "false"
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

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}}
	var updatedDocument Organization
	err := organizationCollection.FindOneAndUpdate(
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
