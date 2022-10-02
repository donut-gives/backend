package orgVerification

type Organization struct {
	Id       string `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
	Verified string `bson:"verified,omitempty" json:"verified,omitempty"`
	Name     string `bson:"name,omitempty" json:"name,omitempty"`
	Photo    string `bson:"photo,omitempty" json:"photo,omitempty"`
	Address  string `bson:"address,omitempty" json:"address,omitempty"`
	Phone    string `bson:"phone,omitempty" json:"phone,omitempty"`
}
