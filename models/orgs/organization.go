package organization

import (
	"donutBackend/models/events"
)

type Point struct {
    Type        string    `json:"type" bson:"type"`
    Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
// NewPoint returns a GeoJSON Point with longitude and latitude.
func NewPoint(long, lat float64) Point {
    return Point{
        "Point",
        []float64{long, lat},
    }
}

type Message struct {
	Body string `bson:"body,omitempty" json:"body,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Designation string `bson:"designation,omitempty" json:"designation,omitempty"`
}

type Statistics struct {
	Story string `bson:"story,omitempty" json:"story,omitempty"`
	Financials string `bson:"financials,omitempty" json:"financials,omitempty"`
	EmployeeCount int `bson:"employeeCount,omitempty" json:"employeeCount,omitempty"`
	References []string `bson:"refrences,omitempty" json:"refrences,omitempty"`
	Cause string `bson:"cause,omitempty" json:"cause,omitempty"`
	Donations int `bson:"donations,omitempty" json:"donations,omitempty"`
	Messages []Message `bson:"messages,omitempty" json:"messages,omitempty"`
}

type Organization struct {
	Id 	 		string `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    	string `bson:"email,omitempty" json:"email,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	YearFounded int `bson:"yearFounded,omitempty" json:"yearFounded,omitempty"`
	DonutName   string `bson:"donutName,omitempty" json:"donutName,omitempty"`
	Name     	string `bson:"name,omitempty" json:"name,omitempty"`
	Photo    	string `bson:"photo,omitempty" json:"photo,omitempty"`
	Location 	string `bson:"location,omitempty" json:"location,omitempty"`
	Tags 		[]string `bson:"tags,omitempty" json:"tags,omitempty"`
    Coordinates Point `json:"coordinates" bson:"coordinates"`
	Events 		[]events.Event `bson:"events,omitempty" json:"events,omitempty"`
	Stats 		Statistics `bson:"stats,omitempty" json:"stats,omitempty"`
}

type OrganizationProfile struct {
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Name     	string `bson:"name,omitempty" json:"name,omitempty"`
	Photo    	string `bson:"photo,omitempty" json:"photo,omitempty"`
	Location 	string `bson:"location,omitempty" json:"location,omitempty"`
	Tags 		[]string `bson:"tags,omitempty" json:"tags,omitempty"`
    Coordinates Point `json:"coordinates" bson:"coordinates"`
	Stats 		Statistics `bson:"stats,omitempty" json:"stats,omitempty"`
	//Website  	string `bson:"website,omitempty" json:"website,omitempty"`
	//Contact    	string `bson:"contact,omitempty" json:"contact,omitempty"`
	//Events 		[]events.Event `bson:"events,omitempty" json:"events,omitempty"`
}

