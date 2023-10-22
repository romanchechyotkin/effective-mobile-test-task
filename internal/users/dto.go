package users

type UserRequestDto struct {
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name,omitempty"`
}

type UserResponseDto struct {
	ID          string `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type UpdateUserDto struct {
	LastName    string `json:"last_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	SecondName  string `json:"second_name,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nationality string `json:"nationality,omitempty"`
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
