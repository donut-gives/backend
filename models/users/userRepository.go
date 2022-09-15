package users

import (
	"context"
	//"fmt"
	."donutBackend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

type UsersRepository struct {}
const usersCollectionName = "users"

var config Config

var usersCollection = new(mongo.Collection)

func init() {
	config.Read()
	// Connect to DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Database.Uri))
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	database := client.Database(config.Database.DatabaseName)
	
	usersCollection = database.Collection(usersCollectionName)
	
}


// Create a new user
func (p *UsersRepository) Insert(user *GoogleUser) (interface{}, error) {
	//fmt.Println("Inserting user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var find_result bson.M
	err := usersCollection.FindOne(
		context.TODO(),
		bson.D{{Key: "email", Value: user.Email}},
		opts,
	).Decode(&find_result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
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
		return nil, err
	}
	//fmt.Printf("found document %v", find_result)

	id:=find_result["_id"]
	string_id:=id.(primitive.ObjectID).Hex()
	return string_id, nil
}

