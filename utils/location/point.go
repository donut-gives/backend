package location

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