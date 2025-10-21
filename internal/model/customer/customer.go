package customer

type Customer struct {
	Document string
	Name     string
	Birthday string

	Favorites []string //list of document of favorite providers
}
