package organization

import (
	"donutBackend/models/events"
	. "donutBackend/utils/location"
)



type Message struct {
	Body string `bson:"body,omitempty" json:"body,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Designation string `bson:"designation,omitempty" json:"designation,omitempty"`
}

type Story struct {
	Paragraphs []string `bson:"paragraphs,omitempty" json:"paragraphs,omitempty"`
	MediaFile string `bson:"mediaFile,omitempty" json:"mediaFile,omitempty"`
	MediaType string `bson:"mediaType,omitempty" json:"mediaType,omitempty"`
}

type References struct {
	ArticleURL string `bson:"articleURL,omitempty" json:"articleURL,omitempty"`
	Title string `bson:"articleTitle,omitempty" json:"articleTitle,omitempty"`
	MediaName string `bson:"mediaName,omitempty" json:"mediaName,omitempty"`
	ImageURL string `bson:"imageURL,omitempty" json:"imageURL,omitempty"`
}

type Impact struct {
	Quantity int `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Unit string `bson:"unit,omitempty" json:"unit,omitempty"`
}

type Statistics struct {
	Story Story `bson:"story,omitempty" json:"story,omitempty"`
	Financials string `bson:"financials,omitempty" json:"financials,omitempty"`
	EmployeeCount int `bson:"employeeCount,omitempty" json:"employeeCount,omitempty"`
	References []References `bson:"refrences,omitempty" json:"refrences,omitempty"`
	Cause string `bson:"cause,omitempty" json:"cause,omitempty"`
	Donations int `bson:"donations,omitempty" json:"donations,omitempty"`
	Impact []Impact `bson:"impact,omitempty" json:"impact,omitempty"`
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
	DonateLink  string `bson:"donateLink,omitempty" json:"donateLink,omitempty"`
}

type OrganizationProfile struct {
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Name     	string `bson:"name,omitempty" json:"name,omitempty"`
	DonutName   string `bson:"donutName,omitempty" json:"donutName,omitempty"`
	Photo    	string `bson:"photo,omitempty" json:"photo,omitempty"`
	Location 	string `bson:"location,omitempty" json:"location,omitempty"`
	Tags 		[]string `bson:"tags,omitempty" json:"tags,omitempty"`
    Coordinates Point `json:"coordinates" bson:"coordinates"`
	Stats 		Statistics `bson:"stats,omitempty" json:"stats,omitempty"`
	DonateLink  string `bson:"donateLink,omitempty" json:"donateLink,omitempty"`
}

