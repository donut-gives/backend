package organization

import (
	"context"
	"donutBackend/db"
	"donutBackend/models/orgVerificationList"
	"donutBackend/models/events"
	///"donutBackend/models/users"
	"errors"
	"time"

	. "donutBackend/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var organizationCollection = new(mongo.Collection)

func init() {
	organizationCollection = db.Get().Collection("organizations")

	// Create a unique index on the email field
	index := mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := organizationCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		Logger.Errorf("Error creating email index in organizations collection: %v", err)
	}
}

func SetPassword(org *Organization) (interface{}, error) {

	password,err := HashPassword(org.Password)
	if err!=nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult bson.M

	err = organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: org.Email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {

			existingOrg,err := orgVerification.Get(org.Email)
			if err!=nil {
				return nil, err
			}

			if existingOrg.Status!="VERIFIED" {
				return nil, errors.New("Organization not verified")
			}

			org.Name = existingOrg.Name
			org.Password = password
			org.Address = existingOrg.Address
			org.Contact = existingOrg.Contact
			org.Email = existingOrg.Email
			org.Photo = existingOrg.Photo

			result, err := organizationCollection.InsertOne(ctx, org)
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

	//update password
	filter := bson.D{{Key: "email", Value: org.Email}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: password},
		}},
	}

	result, err := organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	
	return result.UpsertedID, nil
}

func CheckPwd(email string,password string) (*Organization, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult Organization

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		if(err == mongo.ErrNoDocuments) {
			return nil, errors.New("No such organization found")
		}
		return nil, err
	}

	check:=VerifyPassword(password,findResult.Password)

	findResult.Password = ""

	if check {
		return &findResult, nil
	}

	return nil, errors.New("Password is incorrect")
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

func Get(email string) (Organization,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)
	if err != nil {
		
		if err == mongo.ErrNoDocuments {
			return Organization{}, errors.New("No such organization found")
		}
		return Organization{}, err
	}
	
	return findResult, nil
}

func GetEvents(email string) ([]events.Event, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult Organization

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		if(err == mongo.ErrNoDocuments) {
			return nil, errors.New("No such organization found")
		}
		return nil, err
	}

	return findResult.Events, nil
}

func AddEvent(email string, event events.Event) (interface{},error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	event.OrgEmail = email

	id , err := events.AddEvent(&event)
	if err!=nil {
		return nil,err
	}


	event.Id = id.(primitive.ObjectID).Hex()
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "events", Value: event},
		}},
	}

	_, err = organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil,err
	}

	return id,nil
}

func AddUserToEvent(user events.UserInfo,event events.Event) (error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	option := options.Update()

	filter := bson.D{
		{Key: "email", Value: event.OrgEmail},
		{Key: "events._id", Value: event.Id},
	}
	option.SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "elem._id", Value: event.Id}},
		},
	})
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "events.$[elem].Volunteers", Value: user},
		}},
	}

	_, err := organizationCollection.UpdateOne(ctx, filter, update, option)
	if err != nil {
		return err
	}

	return nil

}

func DeleteEvent(email string, eventId string) (interface{},error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "events", Value: bson.D{
				{Key: "eventId", Value: eventId},
			}},
		}},
	}

	result, err := organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil,err
	}

	if result.MatchedCount == 0 {
		return nil,errors.New("No such organization found")
	}

	err = events.DeleteEvent(eventId)
	if err!=nil {
		return nil,err
	}

	return result.UpsertedID,nil
}

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) (string,error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }

    return string(bytes), nil
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true

    if err != nil {
        check = false
    }

    return check
}