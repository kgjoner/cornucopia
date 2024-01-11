package htypes

type Document struct {
	Number string
	Kind   string `validate:"oneof=CPF CNPJ PASSPORT"`
}

func (d Document) IsZero() bool {
	return d.Number == "" && d.Kind == ""
}
