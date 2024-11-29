package htypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/kgjoner/cornucopia/helpers/validator"
	"github.com/kgjoner/cornucopia/utils/sanitizer"
	"github.com/kgjoner/cornucopia/utils/sliceman"
)

type Document string

func ParseDocument(str string) (*Document, error) {
	d := Document(str)
	return &d, d.IsValid()
}

func (d *Document) IsValid() error {
	if d.IsZero() {
		return nil
	}

	formatted, err := d.Format()
	if err != nil {
		return err
	}

	switch formatted.Kind {
	case "cpf":
		return validateCpf(formatted.Number)
	case "cnpj":
		return validateCnpj(formatted.Number)
	case "passport":
		return validatePassport(formatted.Number)
	default:
		return normalizederr.NewValidationError("not accepted document kind")
	}
}

func (d Document) IsZero() bool {
	return d == ""
}

func (d Document) String() string {
	return string(d)
}

type FormattedDocument struct {
	Kind   string `validate:"required,oneof=cpf cnpj passport"`
	Number string `validate:"required"`
}

func (d *Document) Format() (*FormattedDocument, error) {
	parts := strings.Split(d.String(), ":")
	if len(parts) != 2 {
		parts = make([]string, 2)
		parts[1] = sanitizer.Digit(d.String())
		switch len(parts[1]) {
		case 11:
			parts[0] = "cpf"
		case 14:
			parts[0] = "cnpj"
		default:
			return nil, normalizederr.NewValidationError("unrecognizable document: try to inform its type in the form {type}:{number}")
		}

		*d = Document(fmt.Sprintf("%v:%v", parts[0], parts[1]))

	} else if parts[0] != "passport" {
		parts[1] = sanitizer.Digit(parts[1])
		*d = Document(fmt.Sprintf("%v:%v", parts[0], parts[1]))
	}

	res := &FormattedDocument{
		Kind:   parts[0],
		Number: parts[1],
	}
	return res, validator.Validate(res)
}

func validateCpf(cpf string) error {
	if len(cpf) != 11 {
		return normalizederr.NewValidationError("invalid CPF")
	}

	// Exclude invalid numbers
	invalidNumbers := []string{
		"00000000000",
		"11111111111",
		"22222222222",
		"33333333333",
		"44444444444",
		"55555555555",
		"66666666666",
		"77777777777",
		"88888888888",
		"99999999999",
	}
	if sliceman.IndexOf(invalidNumbers, cpf) != -1 {
		return normalizederr.NewValidationError("invalid CPF")
	}

	// First verification digit
	add := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		add += (10 - i) * digit
	}

	rev := 11 - (add % 11)
	if rev == 10 || rev == 11 {
		rev = 0
	}

	verDig, _ := strconv.Atoi(string(cpf[9]))
	if rev != verDig {
		return normalizederr.NewValidationError("invalid CPF")
	}

	// Second verification digit
	add = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		add += (11 - i) * digit
	}

	rev = 11 - (add % 11)
	if rev == 10 || rev == 11 {
		rev = 0
	}

	verDig, _ = strconv.Atoi(string(cpf[10]))
	if rev != verDig {
		return normalizederr.NewValidationError("invalid CPF")
	}

	return nil
}

func validateCnpj(cnpj string) error {
	if len(cnpj) != 14 {
		return normalizederr.NewValidationError("invalid CNPJ")
	}

	// Exclude invalid numbers
	invalidNumbers := []string{
		"00000000000000",
		"11111111111111",
		"22222222222222",
		"33333333333333",
		"44444444444444",
		"55555555555555",
		"66666666666666",
		"77777777777777",
		"88888888888888",
		"99999999999999",
	}
	if sliceman.IndexOf(invalidNumbers, cnpj) != -1 {
		return normalizederr.NewValidationError("invalid CNPJ")
	}

	// First verification digit
	add := 0
	for i := 0; i < 12; i++ {
		multiplier := 13 - i
		if i < 4 {
			multiplier = 5 - i
		}

		digit, _ := strconv.Atoi(string(cnpj[i]))
		add += multiplier * digit
	}

	rev := 11 - (add % 11)
	if rev == 10 || rev == 11 {
		rev = 0
	}

	verDig, _ := strconv.Atoi(string(cnpj[12]))
	if rev != verDig {
		return normalizederr.NewValidationError("invalid CNPJ")
	}

	// Second verification digit
	add = 0
	for i := 0; i < 13; i++ {
		multiplier := 14 - i
		if i < 5 {
			multiplier = 6 - i
		}

		digit, _ := strconv.Atoi(string(cnpj[i]))
		add += multiplier * digit
	}

	rev = 11 - (add % 11)
	if rev == 10 || rev == 11 {
		rev = 0
	}

	verDig, _ = strconv.Atoi(string(cnpj[13]))
	if rev != verDig {
		return normalizederr.NewValidationError("invalid CNPJ")
	}

	return nil
}

func validatePassport(_ string) error {
	return nil
}
