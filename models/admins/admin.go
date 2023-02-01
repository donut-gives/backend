package admin

type Admin struct {
	Email    string `json:"email"`
	Priviledge []string `json:"priviledge"`
}