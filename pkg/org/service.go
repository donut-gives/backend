package org

import (
	"context"
	"errors"
	"github.com/donut-gives/backend/pkg/db/mongodb"
	"github.com/donut-gives/backend/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Service interface {
	FindWithId(id string) (*Organization, error)
	FindWithUsername(username string) (*Organization, error)
	FindWithEmail(email string) (*Organization, error)
	CreateOrg(profile Profile) (string, error)
	SetPassword(id string, password string) error
	VerifyEmailAndPassword(email string, password string) (bool, error)
	VerifyUsernameAndPassword(username string, password string) (bool, error)
	UpdatePassword(id string, password string) error
	UpdateProfile(profile Profile) error
	UpdateUsername(id string, username string) error
	UpdateEmail(id string, email string) error
	UpdateStory(username string, story Story) error
	UpdateImpacts(username string, impacts []Impact) error
	UpdateReferences(username string, references []Reference) error
	UpdateMessages(username string, messages []Message) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

var orgs = mongodb.GetOrgCollection()

func (service service) FindWithId(id string) (*Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var org *Organization
	err := orgs.FindOne(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
		},
		opts,
	).Decode(org)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (service service) FindWithUsername(username string) (*Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var org *Organization
	err := orgs.FindOne(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
		opts,
	).Decode(org)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (service service) FindWithEmail(email string) (*Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var org *Organization
	err := orgs.FindOne(
		ctx,
		bson.D{
			{Key: "email", Value: email},
		},
		opts,
	).Decode(org)
	if err != nil {
		return nil, err
	}
	return org, nil
}

// CreateOrg creates a new organization and returns the id of the new organization
func (service service) CreateOrg(profile Profile) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doc, err := profile.Map()
	if err != nil {
		return "", err
	}
	opts := options.InsertOne()
	result, errOnInsert := orgs.InsertOne(ctx, doc, opts)
	if errOnInsert != nil {
		return "", errOnInsert
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Inserted ID is not an ObjectID")
	}
	return id.String(), nil
}

func (service service) SetPassword(id string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedPassword, err := utils.Hash(password)
	if err != nil {
		return err
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: hashedPassword},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) VerifyEmailAndPassword(email string, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var org *Organization
	err := orgs.FindOne(
		ctx,
		bson.D{
			{Key: "email", Value: email},
		},
		opts,
	).Decode(org)
	if err != nil {
		return false, err
	}
	result := utils.Verify(password, org.Password)
	return result, nil
}

func (service service) VerifyUsernameAndPassword(username string, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var org *Organization
	err := orgs.FindOne(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
		opts,
	).Decode(org)
	if err != nil {
		return false, err
	}
	result := utils.Verify(password, org.Password)
	return result, nil
}

func (service service) UpdatePassword(id string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: password},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateProfile(profile Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doc, err := profile.Map()
	if err != nil {
		return err
	}
	update := bson.D{
		{Key: "$set", Value: doc},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "username", Value: profile.Username},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateUsername(id string, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "username", Value: username},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateEmail(id string, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "email", Value: email},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateStory(username string, story Story) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "stats.story", Value: story},
		},
		}}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateImpacts(username string, impacts []Impact) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "stats.impacts", Value: impacts},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateReferences(username string, references []Reference) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "stats.references", Value: references},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}

func (service service) UpdateMessages(username string, messages []Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "stats.messages", Value: messages},
		}},
	}
	opts := options.FindOneAndUpdate()
	_, errOnUpdate := orgs.FindOneAndUpdate(
		ctx,
		bson.D{
			{Key: "username", Value: username},
		},
		update,
		opts,
	).DecodeBytes()
	if errOnUpdate != nil {
		return errOnUpdate
	}
	return nil
}
