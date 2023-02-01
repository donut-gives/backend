package pendingemails

type PendingEmail struct {
	Id         string     `bson:"_id,omitempty" json:"_id,omitempty"`
	Email      string     `bson:"email,omitempty" json:"email,omitempty"`
	Body 	 string     `bson:"body,omitempty" json:"body,omitempty"`
	Reason	 string     `bson:"reason,omitempty" json:"reason,omitempty"`
}
