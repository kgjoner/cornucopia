package prim

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/kgjoner/cornucopia/v3/apperr"
	"github.com/kgjoner/cornucopia/v3/validator"
)

type Market string

const (
	MarketBrazil Market = "brazil"
	MarketUSA    Market = "usa"
)

func (m Market) Enumerate() any {
	return []Market{
		MarketBrazil,
		MarketUSA,
	}
}

func MarketByTimezone(timezone string) (Market, error) {
	if timezone == "" {
		return "", apperr.NewValidationError("no timezone specified")
	}

	market, exists := marketsByTimezone[timezone]
	if !exists {
		return "", apperr.NewValidationError("invalid timezone")
	}

	return Market(market), nil
}

func MarketByCurrency(currency Currency) (Market, error) {
	marketByCurrency := map[Currency]Market{
		BRL: MarketBrazil,
		USD: MarketUSA,
	}

	market, exists := marketByCurrency[currency]
	if !exists {
		return "", apperr.NewValidationError("invalid currency")
	}

	return market, nil
}

func (m Market) Language() Locale {
	languageByMarket := map[Market]Locale{
		MarketBrazil: Portuguese,
		MarketUSA:    English,
	}

	return languageByMarket[m]
}

func (m Market) Currency() Currency {
	currencyByMarket := map[Market]Currency{
		MarketBrazil: BRL,
		MarketUSA:    USD,
	}

	return currencyByMarket[m]
}

func (m Market) IsZero() bool {
	return m == ""
}

func (m *Market) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	*m = Market(strings.ToLower(str))
	return validator.Validate(*m)
}

func (m *Market) UnmarshalText(text []byte) error {
	*m = Market(strings.ToLower(string(text)))
	return validator.Validate(*m)
}

/* ================================================================================
	INIT
================================================================================ */

//go:embed assets/marketByTimezone.json
var marketsByTimezoneJSON []byte

var marketsByTimezone map[string]string

func init() {
	marketsByTimezone = make(map[string]string)
	if err := json.Unmarshal(marketsByTimezoneJSON, &marketsByTimezone); err != nil {
		panic("failed to load markets: " + err.Error())
	}
}
