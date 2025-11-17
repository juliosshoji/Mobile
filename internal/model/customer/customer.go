package customer

type Customer struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`

	Favorites []string `json:"favorites"` //list of document of favorite providers
}
