package events

type UserInfo struct {
	Email  string `bson:"email,omitempty" json:"email,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Photo  string `bson:"picture,omitempty" json:"picture,omitempty"`
	Resume string `bson:"resume,omitempty" json:"resume,omitempty"`
}

type Event struct {
	Id         string     `bson:"_id,omitempty" json:"_id,omitempty"`
	OrgEmail   string     `bson:"orgEmail,omitempty" json:"orgEmail,omitempty"`
	Name       string     `bson:"name,omitempty" json:"name,omitempty"`
	Photo      string     `bson:"photo,omitempty" json:"photo,omitempty"`
	Location   string     `bson:"location,omitempty" json:"location,omitempty"`
	Contact    string     `bson:"contact,omitempty" json:"contact,omitempty"`
	Volunteers []UserInfo `bson:"users,omitempty" json:"users,omitempty"`
	OrgLink	string     `bson:"orgLink,omitempty" json:"orgLink,omitempty"`
}
