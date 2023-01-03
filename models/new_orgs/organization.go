package orgVerification

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

type Organization struct {
	Id 	 		string `bson:"_id,omitempty" json:"_id,omitempty"`
	Email    	string `bson:"email,omitempty" json:"email,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Status 		string `bson:"status,omitempty" json:"status,omitempty"`
	Name     	string `bson:"name,omitempty" json:"name,omitempty"`
	Photo    	string `bson:"photo,omitempty" json:"photo,omitempty"`
	Location 	string `bson:"location,omitempty" json:"location,omitempty"`
	Tags 		[]string `bson:"tags,omitempty" json:"tags,omitempty"`
    Coordinates Point `json:"coordinates" bson:"coordinates"`
}