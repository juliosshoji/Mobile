package review

type Review struct {
	Id string `json:"id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Rating uint8 `json:"rating"`

	CustomerId string `json:"customer_id"`
	ProviderId string `json:"provider_id"`
}
