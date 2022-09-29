package organization

import (
	"context"
	"donutBackend/db"
	"donutBackend/models/orgVerificationList"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var organizationCollection = new(mongo.Collection)

func init() {
	organizationCollection = db.Get().Collection("organizations")
}

func Insert(org *Organization) (interface{}, error) {

	password,err := HashPassword(org.Password)
	if err!=nil {
		return nil, err
	}

	existingOrg,err := orgVerification.Get(org.Email)
	if err!=nil {
		return nil, err
	}

	if existingOrg.Verified=="false" {
		return nil, errors.New("Organization not verified")
	}

	org.Name = existingOrg.Name
	org.Password = password
	org.Address = existingOrg.Address
	org.Phone = existingOrg.Phone
	org.Email = existingOrg.Email
	org.Photo = existingOrg.Photo

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

	id := findResult["_id"]
	stringId := id.(primitive.ObjectID).Hex()
	return stringId, nil
}

func Get(email string,password string) (*Organization, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult Organization

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		return nil, err
	}

	check:=VerifyPassword(password,findResult.Password)

	findResult.Password = ""

	if check {
		return &findResult, nil
	}

	return nil, errors.New("Password is incorrect")
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
