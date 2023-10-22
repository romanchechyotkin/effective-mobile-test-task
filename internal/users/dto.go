package users

type UserRequestDto struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	SecondName string `json:"secondName,omitempty"`
}

type UserResponseDto struct {
	LastName    string `json:"lastName"`
	FirstName   string `json:"firstName"`
	SecondName  string `json:"secondName,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type AgeRequestDto struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderRequestDto struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type NationalityRequestDto struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}
