package users

import (
	"context"
	"donutbackend/db"
	. "donutbackend/logger"
	events "donutbackend/models/volunteer"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersCollection = new(mongo.Collection)

func init() {
	usersCollection = db.Get().Collection("users")
	//create email index
	indexview := usersCollection.Indexes()
	index := mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := indexview.CreateOne(context.Background(), index)
	if err != nil {
		Logger.Errorf("Error creating index for users collection: %v", err)
	}

}

// Insert : Create a new user
func Insert(user *GoogleUser) (interface{}, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult bson.M

	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: user.Email}},
		opts,
	).Decode(&findResult)

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
			return id, nil
		}
		return nil, err
	}

	id := findResult["_id"]
	return id, nil
}

func Find(id string) (GoogleUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult GoogleUser
	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: id}},
		opts,
	).Decode(&findResult)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return GoogleUser{}, errors.New("User not found")
		}
		return GoogleUser{}, err
	}

	return findResult, nil
}

func GetUserProfile(email string) (GoogleUserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult GoogleUserProfile
	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return GoogleUserProfile{}, errors.New("User not found")
		}
		return GoogleUserProfile{}, err
	}

	return findResult, nil
}

func GetEvents(email string) ([]events.Opportunity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult GoogleUser
	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return findResult.Events, nil
}

func CheckEventExists(user GoogleUser, event events.Opportunity) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult GoogleUser
	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: user.Email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		return false, err
	}

	for _, e := range findResult.Events {
		if e.Id == event.Id {
			return true, nil
		}
	}
	return false, nil
}

func AddEvent(user GoogleUser, event events.Opportunity) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)

	filter := bson.D{{Key: "email", Value: user.Email}}
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "volunteer", Value: event},
		}},
	}

	var updatedDocument bson.M
	err := usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		option,
	).Decode(&updatedDocument)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEvent(email string, eventId string) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOne()

	var find_result bson.M

	err := usersCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&find_result)
	if err != nil {
		return err
	}

	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)

	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "volunteer", Value: eventId},
		}},
	}

	var updatedDocument bson.M
	err = usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		option,
	).Decode(&updatedDocument)
	if err != nil {
		return err
	}

	return nil
}

func AddBookmark(user GoogleUser, bookmark events.Opportunity) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)

	filter := bson.D{{Key: "email", Value: user.Email}}
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "bookmarks", Value: bookmark},
		}},
	}

	var updatedDocument bson.M
	err := usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		option,
	).Decode(&updatedDocument)
	if err != nil {
		return err
	}

	return nil
}

//Update : Insert a new transaction
func InsertTransaction(userId string, transaction *Transaction) (interface{}, error) {
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOne()

	var find_result bson.M

	err = usersCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: userIdObj}},
		opts,
	).Decode(&find_result)
	if err != nil {
		return nil, err
	}

	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)

	filter := bson.D{{Key: "_id", Value: userIdObj}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "transactions", Value: transaction},
		}},
	}

	var updatedDocument bson.M
	err = usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		option,
	).Decode(&updatedDocument)
	if err != nil {
		return nil, err
	}

	return updatedDocument, nil
}

//Update : Transaction Status
func UpdatePaymentStatus(userId string, transactionId string, status string) (interface{}, error) {
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOne()

	var find_result bson.M
	err = usersCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: userIdObj}},
		opts,
	).Decode(&find_result)
	if err != nil {
		return nil, err
	}

	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)
	option.SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "elem.Id", Value: transactionId}},
		},
	})

	filter := bson.D{
		{Key: "_id", Value: userIdObj},
		{Key: "transactions.Id", Value: transactionId},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "transactions.$[elem].status", Value: status},
		}},
	}

	var updatedDocument bson.M
	err = usersCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		option,
	).Decode(&updatedDocument)
	if err != nil {
		return nil, err
	}

	return updatedDocument, nil
}

//Update Wallet Balance
func UpdateWalletBalance(userId string, amount float64) (interface{}, error) {
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOne()

	var find_result bson.M
	err = usersCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: userIdObj}},
		opts,
	).Decode(&find_result)
	if err != nil {
		return nil, err
	}

	option := options.FindOneAndUpdate()
	option.SetReturnDocument(options.After)

	filter := bson.D{{Key: "_id", Value: userIdObj}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "wallet", Value: amount},
		}},
	}

	var updatedDocument bson.M
	err = usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		option,
	).Decode(&updatedDocument)
	if err != nil {
		return nil, err
	}

	return updatedDocument, nil
}
