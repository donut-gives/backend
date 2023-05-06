package org_verification

import (
	"context"
	"donutbackend/db"
	"errors"
	"time"

	. "donutbackend/logger"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongodb-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type State int

func (s State) String() (string, error) {
	states := [...]string{
		"ANDAMAN AND NICOBAR ISLANDS",
		"ANDHRA PRADESH",
		"ARUNACHAL PRADESH",
		"ASSAM",
		"BIHAR",
		"CHANDIGARH",
		"CHHATTISGARH",
		"DADRA AND NAGAR HAVELI",
		"DAMAN AND DIU",
		"DELHI",
		"GOA",
		"GUJARAT",
		"HARYANA",
		"HIMACHAL PRADESH",
		"JAMMU AND KASHMIR",
		"JHARKHAND",
		"KARNATAKA",
		"KERALA",
		"LAKSHADWEEP",
		"MADHYA PRADESH",
		"MAHARASHTRA",
		"MANIPUR",
		"MEGHALAYA",
		"MIZORAM",
		"NAGALAND",
		"ODISHA",
		"PUDUCHERRY",
		"PUNJAB",
		"RAJASTHAN",
		"SIKKIM",
		"TAMIL NADU",
		"TELANGANA",
		"TRIPURA",
		"UTTAR PRADESH",
		"UTTARAKHAND",
		"WEST BENGAL",
	}
	if len(states) < int(s) {
		return "", errors.New("Invalid State")
	}

	return states[s], nil
}

type Tags int

func (s Tags) String() (string, error) {
	tags := [...]string{
		"children education",
		"animal welfare",
		"healthcare",
		"poverty",
		"sustainable development",
	}
	if len(tags) < int(s) {
		return "", errors.New("Invalid State")
	}

	return tags[s], nil
}

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

func Insert(org *Organization, StateNum int, TagsNum []int) (interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var st State = State(StateNum)

	stateString, err := st.String()
	if err != nil {
		return nil, err
	}

	var tagsString []string

	for _, tag := range TagsNum {
		var tg Tags = Tags(tag)
		tagString, err := tg.String()
		if err != nil {
			return nil, err
		}
		tagsString = append(tagsString, tagString)
	}

	(*org).Location = stateString
	(*org).Tags = tagsString
	(*org).Status = "PENDING"

	insertResult, err := organizationCollection.InsertOne(ctx, org)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
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
		if err == mongo.ErrNoDocuments {
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
	if findResult["status"] != "PENDING" {
		return nil, errors.New("Organization already " + findResult["status"].(string))
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
	if findResult["status"] != "PENDING" {
		return nil, errors.New("Organization already " + findResult["status"].(string))
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

func Find(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{
			{Key: "email", Value: email},
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
