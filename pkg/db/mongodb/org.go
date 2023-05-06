package mongodb

import (
	"context"
	"donutbackend/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orgCollection = new(mongo.Collection)

func createOrgIndex() {
	orgCollection = get().Collection("orgs")
	index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := orgCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		logger.Logger.Errorf("Error creating email index in organizations collection: %v", err)
	}
}

func init() {
	createOrgIndex()
}

func GetOrgCollection() *mongo.Collection {
	return orgCollection
}
