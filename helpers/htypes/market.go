package htypes

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/kgjoner/cornucopia/helpers/i18n"
	"github.com/kgjoner/cornucopia/helpers/normalizederr"
	"github.com/kgjoner/cornucopia/helpers/validator"
)

type Market string

type marketValues struct {
	BRAZIL Market
	USA    Market
}

func MarketByTimezone(timezone string) (Market, error) {
	if timezone == "" {
		return "", normalizederr.NewValidationError("no timezone specified")
	}

	market, exists := marketsByTimezone[timezone]
	if !exists {
		return "", normalizederr.NewValidationError("invalid timezone")
	}

	return Market(market), nil
}

func MarketByCurrency(currency Currency) (Market, error) {
	marketByCurrency := map[Currency]Market{
		BRL: MarketValues.BRAZIL,
		USD: MarketValues.USA,
	}

	market, exists := marketByCurrency[currency]
	if !exists {
		return "", normalizederr.NewValidationError("invalid currency")
	}

	return market, nil
}

func (m Market) Enumerate() any {
	return marketValues{
		"brazil",
		"usa",
	}
}

func (m Market) Language() i18n.Language {
	languageByMarket := map[Market]i18n.Language{
		MarketValues.BRAZIL: i18n.LanguageValues.PT_BR,
		MarketValues.USA:    i18n.LanguageValues.EN_US,
	}

	return languageByMarket[m]
}

func (m Market) Currency() Currency {
	currencyByMarket := map[Market]Currency{
		MarketValues.BRAZIL: BRL,
		MarketValues.USA:    USD,
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

var MarketValues = Market.Enumerate("").(marketValues)

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
