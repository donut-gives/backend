package events

type MediaMeta struct {
	AllowedFileTypes []string `json:"allowedFileTypes"`
	MaxFileSize int `json:"maxFileSize"`
}

type MediaValue struct {
	Url string `json:"url"`
}

type StringTextMeta struct {
	MaxChar int `json:"maxChar"`
	MaxWords int `json:"maxWords"`
	MinChar int `json:"minChar"`
	MinWords int `json:"minWords"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
	Default int `json:"default"`
}

type NumberTextMeta struct {
	UpperLimit int `json:"upperLimit"`
	LowerLimit int `json:"lowerLimit"`
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix"`
	Default int `json:"default"`
}

type TextMeta struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

type TextValue struct {
	Data string `json:"data"`
}

type IdString struct {
	Id string `json:"id"`
	String string `json:"string"`
}

type DropdownMeta struct {
	Options []IdString `json:"options"`
	Default string `json:"default"`
}

type DropdownValue struct {
	Option string `json:"option"`
}

type RangeMeta struct {
	LeastValue int `json:"leastValue"`
	MaxValue int `json:"maxValue"`
}

type RangeValue struct {
	LowerValue int `json:"lowerValue"`
	UpperValue int `json:"upperValue"`
}

type ChoiceMeta struct {
	Options []IdString `json:"options"`
	Default []string `json:"default"`
	Multi bool `json:"multi"`
}

type ChoiceValue struct {
	Options []string `json:"options"`
}

type UrlMeta struct {
	BaseUrl string `json:"baseUrl"`
}

type UrlValue struct {
	Url string `json:"url"`
}

//dont convert type field to int rather take a seperate array in the body of request to 
//insdicate the type of each field to use ENUMs and set Type string on the basis of this
type FormValue struct {
	Type string `json:"type"`
	Value interface{} `json:"value"`
}

//dont bind into this directly from bodt rather use the Type Array to create FormFields 
//during request and bind meta seperately and then add this FormFields to the Event 
type FormFields struct {
	Type string `json:"type"`
	Required bool `json:"required"` //try string if false not inserting in database
	Title string `json:"title"`
	Meta interface{} `json:"meta"`
}

type UserInfo struct {
	Email  string `bson:"email,omitempty" json:"email,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Photo  string `bson:"picture,omitempty" json:"picture,omitempty"`
	FormFields []FormValue `bson:"formFields,omitempty" json:"formFields,omitempty"`
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
	FormFields []FormFields `bson:"formFields,omitempty" json:"formFields,omitempty"`
}
