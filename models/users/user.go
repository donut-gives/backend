package users

//import "go.mongodb.org/mongo-driver/bson/primitive"

type GoogleUser struct {
	Email     string `bson:"email,omitempty" json:"email,omitempty"`
	Verified  bool   `bson:"email_verified,omitempty" json:"email_verified,omitempty"`
	Name      string `bson:"name,omitempty" json:"name,omitempty"`
	Photo     string `bson:"picture,omitempty" json:"picture,omitempty"`
	FirstName string `bson:"given_name,omitempty" json:"given_name,omitempty"`
	LastName  string `bson:"family_name,omitempty" json:"family_name,omitempty"`
}