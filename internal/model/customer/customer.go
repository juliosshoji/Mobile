package customer

type Customer struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`

	Favorites []string `json:"favorites"` //list of document of favorite providers
}
