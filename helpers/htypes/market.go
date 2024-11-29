package htypes

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kgjoner/cornucopia/helpers/i18n"
	"github.com/kgjoner/cornucopia/helpers/normalizederr"
)

type Market string

type marketValues struct {
	BRAZIL Market
	USA    Market
}

func MarketByTimezone(timezone string) (Market, error) {
	if timezone == "" {
		// TODO: add error code
		return "", normalizederr.NewRequestError("no timezone specified", "")
	}

	data, _ := ioutil.ReadFile("./pkg/domains/shop/shoptyp/assets/marketByTimezone.json")

	var marketsMap map[string]string
	json.Unmarshal(data, &marketsMap)

	market, exists := marketsMap[timezone]
	if !exists {
		// TODO: add error code
		return "", normalizederr.NewRequestError("invalid timezone", "")
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
		// TODO: add error code
		return "", normalizederr.NewRequestError("invalid currency", "")
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

var MarketValues = Market.Enumerate("").(marketValues)
