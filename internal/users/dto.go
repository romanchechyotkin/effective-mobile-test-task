package users

type User struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName,omitempty"`
}
