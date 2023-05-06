package organization

import (
	"github.com/donut-gives/backend/models/volunteer"
	. "github.com/donut-gives/backend/utils/location"
)

type Message struct {
	Body        string `bson:"body,omitempty" json:"body,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	Designation string `bson:"designation,omitempty" json:"designation,omitempty"`
}

type Story struct {
	Paragraphs []string `bson:"paragraphs,omitempty" json:"paragraphs,omitempty"`
	MediaFile  string   `bson:"media_file,omitempty" json:"media_file,omitempty"`
	MediaType  string   `bson:"media_type,omitempty" json:"media_type,omitempty"`
}

type References struct {
	ArticleURL string `bson:"article_url,omitempty" json:"article_url,omitempty"`
	Title      string `bson:"article_title,omitempty" json:"article_title,omitempty"`
	MediaName  string `bson:"media_name,omitempty" json:"media_name,omitempty"`
	ImageURL   string `bson:"image_url,omitempty" json:"image_url,omitempty"`
}

type Impact struct {
	Quantity int    `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Unit     string `bson:"unit,omitempty" json:"unit,omitempty"`
}

type Statistics struct {
	Story         Story        `bson:"story,omitempty" json:"story,omitempty"`
	EmployeeCount int          `bson:"employee_count,omitempty" json:"employee_count,omitempty"`
	References    []References `bson:"references,omitempty" json:"references,omitempty"`
	Donations     int          `bson:"donations,omitempty" json:"donations,omitempty"`
	Impact        []Impact     `bson:"impact,omitempty" json:"impact,omitempty"`
	Messages      []Message    `bson:"messages,omitempty" json:"messages,omitempty"`
}

type Organization struct {
	Id            string                  `bson:"_id,omitempty" json:"_id,omitempty"`
	Email         string                  `bson:"email,omitempty" json:"email,omitempty"`
	Password      string                  `bson:"password,omitempty" json:"password,omitempty"`
	Description   string                  `bson:"description,omitempty" json:"description,omitempty"`
	YearFounded   int                     `bson:"year_founded,omitempty" json:"year_founded,omitempty"`
	Username      string                  `bson:"username,omitempty" json:"username,omitempty"`
	Name          string                  `bson:"name,omitempty" json:"name,omitempty"`
	Photo         string                  `bson:"photo,omitempty" json:"photo,omitempty"`
	Location      string                  `bson:"location,omitempty" json:"location,omitempty"`
	Tags          []string                `bson:"tags,omitempty" json:"tags,omitempty"`
	Opportunities []volunteer.Opportunity `bson:"opportunities,omitempty" json:"opportunities,omitempty"`
	Stats         Statistics              `bson:"stats,omitempty" json:"stats,omitempty"`
	DonateLink    string                  `bson:"donate_link,omitempty" json:"donate_link,omitempty"`
}

type Profile struct {
	Description string     `bson:"description,omitempty" json:"description,omitempty"`
	Name        string     `bson:"name,omitempty" json:"name,omitempty"`
	Username    string     `bson:"username,omitempty" json:"username,omitempty"`
	Photo       string     `bson:"photo,omitempty" json:"photo,omitempty"`
	Location    string     `bson:"location,omitempty" json:"location,omitempty"`
	Tags        []string   `bson:"tags,omitempty" json:"tags,omitempty"`
	Coordinates Point      `json:"coordinates" bson:"coordinates"`
	Stats       Statistics `bson:"stats,omitempty" json:"stats,omitempty"`
	DonateLink  string     `bson:"donate_link,omitempty" json:"donate_link,omitempty"`
}
