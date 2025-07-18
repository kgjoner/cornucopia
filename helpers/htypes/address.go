package htypes

type Address struct {
	Line1        string  `json:"line1" validate:"required"`
	Number       string  `json:"number" validate:"required"`
	Line2        string  `json:"line2,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	City         string  `json:"city" validate:"required"`
	State        string  `json:"state" validate:"required"`
	Country      Country `json:"country" validate:"required"`
	ZipCode      ZipCode `json:"zipCode" validate:"required"`
}

// IsZero returns true if the Address contains only zero/empty values
func (a Address) IsZero() bool {
	return a.Line1 == "" &&
		a.Number == "" &&
		a.Line2 == "" &&
		a.Neighborhood == "" &&
		a.City == "" &&
		a.State == "" &&
		a.Country.IsZero() &&
		a.ZipCode.IsZero()
}

type ZipCode string

func (z ZipCode) IsZero() bool {
	return z == ""
}
