package events

type Event struct {
	Id 	 		string `bson:"_id,omitempty" json:"_id,omitempty"`
	OrgId 		string `bson:"orgId,omitempty" json:"orgId,omitempty"`
	Name     	string `bson:"name,omitempty" json:"name,omitempty"`
	Photo    	string `bson:"photo,omitempty" json:"photo,omitempty"`
	Location 	string `bson:"location,omitempty" json:"location,omitempty"`
	Contact		string `bson:"contact,omitempty" json:"contact,omitempty"`
	Users 		[]string `bson:"users,omitempty" json:"users,omitempty"`
}