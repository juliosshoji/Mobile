package customer

type ServicesDone struct {
	ProviderDocument string `json:"provider_document"`
	ServiceDate      string `json:"service_date"`
	ReviewID         string `json:"review_id"`
}

type Customer struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`

	Favorites []string `json:"favorites"`

	Services []ServicesDone `json:"services_done"`
}
