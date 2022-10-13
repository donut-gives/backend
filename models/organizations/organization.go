package organization

import (
	"donutBackend/models/events"
)

type Organization struct {

	Id 	 		string `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    	string `bson:"email,omitempty" json:"email,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Status 		string `bson:"status,omitempty" json:"status,omitempty"`
	Name     	string `bson:"name,omitempty" json:"name,omitempty"`
	Photo    	string `bson:"photo,omitempty" json:"photo,omitempty"`
	Location 	string `bson:"location,omitempty" json:"location,omitempty"`
	Address  	string `bson:"address,omitempty" json:"address,omitempty"`
	Website  	string `bson:"website,omitempty" json:"website,omitempty"`
	Contact    	string `bson:"contact,omitempty" json:"contact,omitempty"`
	Events 		[]events.Event `bson:"events,omitempty" json:"events,omitempty"`
	
}

