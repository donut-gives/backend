package organization

import (
	"context"
	"donutbackend/db"
	org_verification "donutbackend/models/new_orgs"
	"donutbackend/models/volunteer"
	"strings"
	///"donutBackend/models/users"
	"errors"
	"time"

	. "donutbackend/logger"
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
		Keys: bson.D{
			{Key: "email", Value: 1},
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := organizationCollection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		Logger.Errorf("Error creating email index in organizations collection: %v", err)
	}
}

func SetPassword(org *Organization) (interface{}, error) {

	password, err := HashPassword(org.Password)
	if err != nil {
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

			existingOrg, err := org_verification.Get(org.Email)
			if err != nil {
				return nil, err
			}

			if existingOrg.Status != "VERIFIED" {
				return nil, errors.New("Organization not verified")
			}

			lowerName := strings.ToLower(existingOrg.Name)
			org.Username = strings.Split(lowerName, " ")[0]
			org.Name = existingOrg.Name
			org.Password = password
			org.Email = existingOrg.Email
			org.Photo = existingOrg.Photo
			org.Tags = existingOrg.Tags
			org.Location = existingOrg.Location
			org.Description = existingOrg.Description

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

func CheckPwd(email string, password string) (*Organization, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opts := options.FindOne()
	var findResult Organization

	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "email", Value: email}},
		opts,
	).Decode(&findResult)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No such organization found")
		}
		return nil, err
	}

	check := VerifyPassword(password, findResult.Password)

	findResult.Password = ""

	if check {
		return &findResult, nil
	}

	return nil, errors.New("Password is incorrect")
}

func Find(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{
			{Key: "_id", Value: id},
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

func GetOrg(username string) (Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Profile
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "username", Value: username}},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Profile{}, errors.New("No such organization found")
		}
		return Profile{}, err
	}

	return findResult, nil
}

func GetOpportunities(username string) ([]volunteer.Opportunity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "username", Value: username}},
		opts,
	).Decode(&findResult)
	print(findResult.Username)
	print(findResult.Opportunities)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//return empty array
			return []volunteer.Opportunity{}, errors.New("no such organization found")
		}
		return []volunteer.Opportunity{}, err
	}
	return findResult.Opportunities, nil
}

func GetOpportunity(username string, opportunityId string) (volunteer.Opportunity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{
			{Key: "username", Value: username},
			{Key: "volunteer", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "_id", Value: opportunityId}},
				},
			}},
		},
		opts,
	).Decode(&findResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//return empty array
			return volunteer.Opportunity{}, errors.New("no such organization found")
		}
		return volunteer.Opportunity{}, err
	}
	return findResult.Opportunities[0], nil
}

func GetStats(donutName string) (interface{}, error) {

	//get stats
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "donutName", Value: donutName}},
		opts,
	).Decode(&findResult)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No such organization found")
		}
		return nil, err
	}

	return findResult.Stats, nil
}

func GetMessages(org string) ([]Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "donutName", Value: org}},
		opts,
	).Decode(&findResult)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return []Message{}, errors.New("No such organization found")
		}
		return []Message{}, err
	}

	return findResult.Stats.Messages, nil
}

func GetEmployees(org string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "donutName", Value: org}},
		opts,
	).Decode(&findResult)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return 0, errors.New("No such organization found")
		}
		return 0, err
	}

	return findResult.Stats.EmployeeCount, nil
}

func GetReferences(org string) ([]References, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "donutName", Value: org}},
		opts,
	).Decode(&findResult)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No such organization found")
		}
		return nil, err
	}

	return findResult.Stats.References, nil
}

func GetStory(org string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := options.FindOne()
	var findResult Organization
	err := organizationCollection.FindOne(
		ctx,
		bson.D{{Key: "donutName", Value: org}},
		opts,
	).Decode(&findResult)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return "", errors.New("No such organization found")
		}
		return "", err
	}

	return findResult.Stats.Story, nil
}

func AddOpportunity(id string, event volunteer.Opportunity) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opId, err := volunteer.AddVolunteerOpportunity(&event)
	if err != nil {
		return nil, err
	}

	event.Id = opId.(primitive.ObjectID).Hex()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "volunteer", Value: event},
		}},
	}

	_, err = organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func AddSubmission(orgId string, submission volunteer.Submission, opportunity volunteer.Opportunity) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	option := options.Update()

	filter := bson.D{
		{Key: "_id", Value: orgId},
		{Key: "volunteer._id", Value: opportunity.Id},
	}
	option.SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "elem._id", Value: opportunity.Id}},
		},
	})
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "volunteer.$[elem].Submissions", Value: submission},
		}},
	}

	_, err := organizationCollection.UpdateOne(ctx, filter, update, option)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOpportunity(username string, opportunityId string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "volunteer", Value: bson.D{
				{Key: "_id", Value: opportunityId},
			}},
		}},
	}

	result, err := organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("No such organization found")
	}

	err = volunteer.DeleteEvent(opportunityId)
	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func UpdateOrgProfile(username string, profile Profile) (interface{}, error) {

	profile.Username = username

	prim, err := toDoc(profile)
	primitiveMap := prim.Map()

	update := bson.D{
		{Key: "$set", Value: primitiveMap},
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	filter := bson.D{{Key: "username", Value: username}}
	// update := bson.D{
	// 	{Key: "$set", Value: bson.D{
	// 		{Key: "profile", Value: profile},
	// 		{Key: "profile", Value: profile},
	// 	}},
	// }

	result, err := organizationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("No such organization found")
	}

	return result.UpsertedID, nil
}

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) (string, error) {
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
