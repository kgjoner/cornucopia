package htypes

type Document struct {
	Number string
	Kind   string `validate:"oneof=CPF CNPJ PASSPORT"`
}
