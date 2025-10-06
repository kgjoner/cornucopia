package htypes

import (
	"encoding/json"
	"fmt"

	"github.com/kgjoner/cornucopia/v2/helpers/apperr"
	"github.com/kgjoner/cornucopia/v2/helpers/validator"
)

type Price map[Currency]PriceValues

func NewPrice() Price {
	return make(Price)
}

func ParsePrice(jsonData string) (Price, error) {
	var price Price
	json.Unmarshal([]byte(jsonData), &price)

	err := validator.Validate(price)
	if err != nil {
		return price, err
	}

	return price, nil
}

func (p Price) UpsertCurrency(currency Currency, fullPrice int) error {
	err := validator.Validate(currency)
	if err != nil {
		return err
	}

	p[currency], err = newPriceValues(fullPrice)
	if err != nil {
		return err
	}

	return nil
}

func (p Price) DeleteCurrency(currency Currency) error {
	err := validator.Validate(currency)
	if err != nil {
		return err
	}

	delete(p, currency)
	return nil
}

func (p Price) SetDiscount(discount int, currency Currency) error {
	values, err := p.Values(currency)
	if err != nil {
		return err
	}

	err = values.setDiscount(discount)
	if err != nil {
		return err
	}

	p[currency] = values
	return nil
}

func (p Price) SetDiscountPercentage(discountPercentage int, currency Currency) error {
	values, err := p.Values(currency)
	if err != nil {
		return err
	}

	err = values.setDiscountPercentage(discountPercentage)
	if err != nil {
		return err
	}

	p[currency] = values
	return nil
}

func (p Price) RemoveDiscount(currency Currency) error {
	values, err := p.Values(currency)
	if err != nil {
		return err
	}

	values.removeDiscount()
	p[currency] = values
	return nil
}

func (p Price) Values(currency Currency) (PriceValues, error) {
	err := validator.Validate(currency)
	if err != nil {
		return PriceValues{}, err
	}

	values, exists := p[currency]
	if !exists {
		return PriceValues{}, apperr.NewValidationError(fmt.Sprintf("%v is not set for this price yet. Add it first.", currency))
	}

	return values, nil
}

func (p Price) String() string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}

type Currency string

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
)

func (c Currency) Enumerate() any {
	return []Currency{
		BRL,
		USD,
	}
}

type PriceValues struct {
	FullPrice int `json:"fullPrice" validate:"required,min=100"`
	SalePrice int `json:"salePrice" validate:"required,min=100"`
}

func newPriceValues(fullPrice int) (PriceValues, error) {
	values := PriceValues{
		fullPrice,
		fullPrice,
	}

	err := validator.Validate(values)
	if err != nil {
		return PriceValues{}, err
	}

	return values, nil
}

func (p PriceValues) IsValid() error {
	errs := make(map[string]string)

	if p.SalePrice > p.FullPrice {
		errs["SalePrice"] = "must not be higher than full price"
	}

	if len(errs) != 0 {
		return apperr.NewMapError(errs)
	}

	return nil
}

func (p *PriceValues) setDiscount(discount int) error {
	p.SalePrice = p.FullPrice - discount

	err := validator.Validate(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *PriceValues) setDiscountPercentage(discountPercentage int) error {
	p.SalePrice = p.FullPrice * (100 - discountPercentage) / 100

	err := validator.Validate(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *PriceValues) removeDiscount() {
	p.SalePrice = p.FullPrice
}
