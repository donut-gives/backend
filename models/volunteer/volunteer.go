package volunteer

type MediaMeta struct {
	AllowedFileTypes []string `json:"allowed_file_types"`
	MaxFileSize      int      `json:"max_file_size"`
}

type MediaValue struct {
	Url string `json:"url"`
}

type StringTextMeta struct {
	MaxChar  int    `json:"max_char"`
	MaxWords int    `json:"max_words"`
	MinChar  int    `json:"min_char"`
	MinWords int    `json:"min_words"`
	Prefix   string `json:"prefix"`
	Suffix   string `json:"suffix"`
	Default  int    `json:"default"`
}

type NumberTextMeta struct {
	UpperLimit int    `json:"upper_limit"`
	LowerLimit int    `json:"lower_limit"`
	Prefix     string `json:"prefix"`
	Suffix     string `json:"suffix"`
	Default    int    `json:"default"`
}

type TextMeta struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type TextValue struct {
	Data string `json:"data"`
}

type IdString struct {
	Id     string `json:"id"`
	String string `json:"string"`
}

type DropdownMeta struct {
	Options []IdString `json:"options"`
	Default string     `json:"default"`
}

type DropdownValue struct {
	Option string `json:"option"`
}

type RangeMeta struct {
	LeastValue int `json:"leastValue"`
	MaxValue   int `json:"max_value"`
}

type RangeValue struct {
	LowerValue int `json:"lower_value"`
	UpperValue int `json:"upper_value"`
}

type ChoiceMeta struct {
	Options []IdString `json:"options"`
	Default []string   `json:"default"`
	Multi   bool       `json:"multi"`
}

type ChoiceValue struct {
	Options []string `json:"options"`
}

type UrlMeta struct {
	BaseUrl string `json:"base_url"`
}

type UrlValue struct {
	Url string `json:"url"`
}

// FormValue don't convert type field to int rather take a separate array in the body of request to
// indicate the type of each field to use ENUMs and set Type string on the basis of this
type FormValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// FormField don't bind into this directly from body rather use the Type Array to create FormField
// during request and bind meta separately and then add this FormField to Opportunity
type FormField struct {
	Type     string      `json:"type"`
	Required bool        `json:"required"` //try string if false not inserting in database
	Title    string      `json:"title"`
	Meta     interface{} `json:"meta"`
}

type Submission struct {
	Email      string      `bson:"email,omitempty" json:"email,omitempty"`
	Name       string      `bson:"name,omitempty" json:"name,omitempty"`
	Photo      string      `bson:"picture,omitempty" json:"picture,omitempty"`
	FormFields []FormValue `bson:"formFields,omitempty" json:"formFields,omitempty"`
}

type Opportunity struct {
	Id          string       `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string       `bson:"title,omitempty" json:"title,omitempty"`
	Description string       `bson:"description,omitempty" json:"description,omitempty"`
	Location    string       `bson:"location,omitempty" json:"location,omitempty"`
	FormFields  []FormField  `bson:"formFields,omitempty" json:"formFields,omitempty"`
	Submissions []Submission `bson:"users,omitempty" json:"users,omitempty"`
}
