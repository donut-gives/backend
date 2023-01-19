package weblinks

import (
	"context"
	"donutBackend/db"
	"errors"
	"strings"
	"time"

	. "donutBackend/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var weblinksCollection = new(mongo.Collection)
var baseUrl = "https://donut.com"
func init() {
	weblinksCollection = db.Get().Collection("weblinks")

	// Create a unique index on the email field
	indexView := weblinksCollection.Indexes()
	mod := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := indexView.CreateOne(context.Background(), mod)
	if err != nil {
		Logger.Errorf("Error creating email index in organizationsList collection: %v", err)
	}
	
}

//get link by id
func GetLink(id string) (Link, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOne()

	linkId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Link{}, err
	}
	var findResult Link
	err = weblinksCollection.FindOne(
		ctx,
		bson.D{{Key: "_id", Value: linkId}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Link{}, errors.New("No link found")
		}
		return Link{}, err
	}
	return findResult, nil
}
//increment link count
func IncrementLinkCount(id string) (error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	linkId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	var findResult Link
	err = weblinksCollection.FindOneAndUpdate(
		ctx,
		bson.D{{Key: "_id", Value: linkId}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "count", Value: 1}}}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("No link found")
		}
		return err
	}
	return nil
}

//get all links
func GetLinks() ([]Link, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.Find()

	var findResult []Link
	cur, err := weblinksCollection.Find(
		ctx,
		bson.D{},
		opts,
	)	
	if err != nil {
		return []Link{}, err
	}
	// check EOF
	for cur.Next(ctx) {
		var result Link
		err := cur.Decode(&result)
		if err != nil {
			return []Link{}, err
		}
		findResult = append(findResult, result)
	}
	if err := cur.Err(); err != nil {
		return []Link{}, err
	}
	cur.Close(ctx)
	return findResult, nil
}

//add or update link
func AddOrUpdateLink(link Link) (Link,bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	var linkId primitive.ObjectID
	var err error

	if link.Id == "" {
		linkId = primitive.NewObjectID()
	} else {
		linkId, err = primitive.ObjectIDFromHex(link.Id)
		if err != nil {
			return  Link{},false,err
		}
	}

	var findResult Link
	err = weblinksCollection.FindOneAndUpdate(
		ctx,
		bson.D{{Key: "_id", Value: linkId}},
		bson.D{{Key: "$set", Value: link}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// link does not exist, so create it
			link.Count = 1
			result, err := weblinksCollection.InsertOne(ctx, link)
			if err != nil {
				return Link{},false, err
			}
			oid,ok := result.InsertedID.(primitive.ObjectID)
			if !ok {
				return Link{},false, errors.New("Error converting to object id")
			}
			link.Id = oid.Hex()
			return link,true, nil
		}
		return Link{},false, err
	}
	return findResult,false, nil
}

func AddLink(name string) (Link, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	//linkId := primitive.NewObjectID()
	//cretae link
	link := Link{
		Name: name,
		Count: 1,
	}
	result, err := weblinksCollection.InsertOne(ctx, link)
	if err != nil {
		//check if link already exists
		if strings.Contains(err.Error(), "duplicate key error") {
			return Link{}, errors.New("Link already exists")
		}
		return Link{}, err
	}
	oid,ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return Link{}, errors.New("Error converting to object id")
	}
	link.Id = oid.Hex()
	return link, nil
}

func UpdateLink(id string,name string) (Link, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)

	linkId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Link{}, err
	}

	//update naem
	link := Link{
		Name: name,
	}

	var findResult Link
	err = weblinksCollection.FindOneAndUpdate(
		ctx,
		bson.D{{Key: "_id", Value: linkId}},
		bson.D{{Key: "$set", Value: link}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Link{}, errors.New("No link found")
		}
		return Link{}, err
	}
	return findResult, nil
}

//delete link
func DeleteLink(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOneAndDelete()

	linkId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = weblinksCollection.FindOneAndDelete(
		ctx,
		bson.D{{Key: "_id", Value: linkId}},
		opts,
	).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("No link found")
		}
		return err
	}
	return nil
}