package users

import (
	"context"
	"donutBackend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var usersCollection = new(mongo.Collection)

func init() {
	usersCollection = db.Get().Collection("users")
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
			stringId := id.(primitive.ObjectID).Hex()
			return stringId, nil
		}
		return nil, err
	}

	id := findResult["_id"]
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
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
