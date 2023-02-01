package weblinks

type Link struct {
	Id         string     `bson:"_id,omitempty" json:"_id,omitempty"`
	Name 	 string     `bson:"name,omitempty" json:"name,omitempty"`
	//Url 	 string     `bson:"url,omitempty" json:"url,omitempty"`
	Count 	 int        `bson:"count,omitempty" json:"count,omitempty"`
	Internal string       `bson:"internal,omitempty" json:"internal,omitempty"`
}
