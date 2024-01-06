package htypes

type Address struct {
	Line1        string `validate:"required"`
	Number       string `validate:"required"`
	Line2        string
	Neighborhood string
	City         string  `validate:"required"`
	State        string  `validate:"required"`
	Country      Country `validate:"required"`
	ZipCode      ZipCode `validate:"required"`
}

type ZipCode string
