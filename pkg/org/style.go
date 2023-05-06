package org

type BaseColor struct {
	BackgroundColor string `bson:"background_color,omitempty" json:"background_color,omitempty"`
	SurfaceColor    string `bson:"surface_color,omitempty" json:"surface_color,omitempty"`
}

type PrimaryColor struct {
	Main       string `bson:"main,omitempty" json:"main,omitempty"`
	Light1     string `bson:"light1,omitempty" json:"light1,omitempty"`
	Light2     string `bson:"light2,omitempty" json:"light2,omitempty"`
	UltraLight string `bson:"ultra_light,omitempty" json:"ultra_light,omitempty"`
}

type SecondaryColor struct {
	Main    string `bson:"main,omitempty" json:"main,omitempty"`
	Variant string `bson:"variant,omitempty" json:"variant,omitempty"`
}

type NeutralColor struct {
	Main    string `bson:"main,omitempty" json:"main,omitempty"`
	Variant string `bson:"variant,omitempty" json:"variant,omitempty"`
}

type TextColor struct {
	Main     string `bson:"main,omitempty" json:"main,omitempty"`
	Inactive string `bson:"inactive,omitempty" json:"inactive,omitempty"`
	Inverted string `bson:"inverted,omitempty" json:"inverted,omitempty"`
}

type ShadowColor struct {
	Light string `bson:"light,omitempty" json:"light,omitempty"`
	Medium string `bson:"medium,omitempty" json:"medium,omitempty"`
	Dark string `bson:"dark,omitempty" json:"dark,omitempty"`
}

type ErrorColor struct {
	Error string `bson:"error,omitempty" json:"error,omitempty"`
}

type Style struct {
	BaseColor BaseColor `bson:"base_color,omitempty" json:"base_color,omitempty"`
	PrimaryColor PrimaryColor `bson:"primary_color,omitempty" json:"primary_color,omitempty"`
	SecondaryColor SecondaryColor `bson:"secondary_color,omitempty" json:"secondary_color,omitempty"`
	NeutralColor NeutralColor `bson:"neutral_color,omitempty" json:"neutral_color,omitempty"`
	TextColor TextColor `bson:"text_color,omitempty" json:"text_color,omitempty"`
	ShadowColor ShadowColor `bson:"shadow_color,omitempty" json:"shadow_color,omitempty"`
	ErrorColor ErrorColor `bson:"error_color,omitempty" json:"error_color,omitempty"`
}
