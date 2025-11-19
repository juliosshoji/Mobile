package provider

type Specialty string

var (
	validSpecialties = map[Specialty]int{
		"GARDENER":   1,
		"ELETRICIAN": 2,
		"COOK":       3,
	}
	validContactTypes = map[Contact]int{
		"EMAIL": 1,
		"PHONE": 2,
	}
)

func (ref Specialty) IsValid() bool {
	_, ok := validSpecialties[ref]
	return ok
}

type Contact string

func (ref Contact) IsValid() bool {
	_, ok := validContactTypes[ref]
	return ok
}

type Provider struct {
	Document string `json:"document"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`

	Specialties []Specialty `json:"specialty"`

	ContactType    Contact `json:"contact_type"`
	ContactAddress string  `json:"contact_address"`

	ProfilePhoto string `json:"profile_photo"`
}
