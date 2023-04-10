package email_sender

type EmailSender struct {
	Id     string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Email  string `bson:"email,omitempty" json:"email,omitempty"`
	Active string `bson:"active,omitempty" json:"active,omitempty"`
	Token  string `bson:"token,omitempty" json:"token,omitempty"`
}
