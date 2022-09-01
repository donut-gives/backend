package repository

import (
	"context"
	"fmt"
	"donutBackend/models"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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


// Get all Places
// func (p *UsersRepository) FindAll() ([]models.Place, error) {
// 	var users []models.GoogleUser

// 	findOptions := options.Find()
// 	findOptions.SetLimit(100)

// 	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
// 	// Finding multiple documents returns a cursor
// 	cur, err := usersCollection.Find(ctx, bson.D{{}}, findOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Iterate through the cursor
// 	defer cur.Close(ctx)
// 	for cur.Next(ctx) {
// 		var result models.Place
// 		err := cur.Decode(&result)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		places = append(places, result)
// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	return places, err
// }

// Create a new Place
func (p *UsersRepository) Insert(user *models.GoogleUser) (interface{}, error) {
	fmt.Println("Inserting user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := usersCollection.InsertOne(ctx, user)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	//fmt.Println("Inserted a single user: ", result.InsertedID)
	return result.InsertedID, err
}

// Delete an existing Place
// func (p *UsersRepository) Delete(id string) error {
// 	objectId, err := primitive.ObjectIDFromHex(id)
// 	filter := bson.M{"_id": objectId}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	result, err := usersCollection.DeleteOne(ctx, filter)
// 	fmt.Println("Deleted a single document: ", result.DeletedCount)
// 	return err
// }