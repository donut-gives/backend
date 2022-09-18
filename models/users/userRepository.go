package users

import (
	"context"
	. "donutBackend/config"
	. "donutBackend/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type UsersRepository struct{}

const usersCollectionName = "users"

var usersCollection = new(mongo.Collection)

func init() {
	// Connect to DB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Configs.DB.Url))
	if err != nil {
		Logger.Fatal(err)
	}
	// defer client.Disconnect(ctx)
	segments := strings.Split(Configs.DB.Url, string('/'))
	database := client.Database(segments[len(segments)-1])

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
			id := result.InsertedID
			string_id := id.(primitive.ObjectID).Hex()
			return string_id, nil
		}
		return nil, err
	}
	//fmt.Printf("found document %v", find_result)

	id := find_result["_id"]
	string_id := id.(primitive.ObjectID).Hex()
	return string_id, nil
}
