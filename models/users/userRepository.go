package users

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type UsersRepository struct {}
const usersCollectionName = "users"

var usersCollection = new(mongo.Collection)

func init() {

	// Connect to DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	database := client.Database("donut")
	
	usersCollection = database.Collection(usersCollectionName)
	
}


// Create a new user
func (p *UsersRepository) Insert(user *GoogleUser) (interface{}, error) {
	fmt.Println("Inserting user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	//fmt.Println("Inserted a single user: ", result.InsertedID)
	id:=result.InsertedID
	string_id:=id.(primitive.ObjectID).Hex()
	return string_id, nil
}

