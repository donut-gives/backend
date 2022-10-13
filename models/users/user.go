package users

//import "go.mongodb.org/mongo-driver/bson/primitive"
import (
	"donutBackend/models/events"
)

type Transaction struct {
	MerchantId string  `json:"merchantId"`
	OrderId    string  `json:"Id" bson:"Id"`
	Timestamp  string  `json:"timestamp" bson:"timestamp"`
	Amount     float64 `json:"amount" bson:"amount"`
	Type       string  `json:"type" bson:"type"`
	Mode       string  `json:"modeOfPayment" bson:"modeOfPayment"`
	Status     string  `json:"status" bson:"status"`
}

type GoogleUser struct {
	Email        string        `bson:"email,omitempty" json:"email,omitempty"`
	Verified     bool          `bson:"emailVerified,omitempty" json:"emailVerified,omitempty"`
	Name         string        `bson:"name,omitempty" json:"name,omitempty"`
	Photo        string        `bson:"picture,omitempty" json:"picture,omitempty"`
	FirstName    string        `bson:"fitstName,omitempty" json:"firstName,omitempty"`
	LastName     string        `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Wallet       float64       `bson:"wallet,omitempty" json:"wallet,omitempty"`
	Transactions []Transaction `bson:"transactions,omitempty" json:"transactions,omitempty"`
	Events       []events.Event       `bson:"events,omitempty" json:"events,omitempty"`
}
