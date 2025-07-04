package enum

import (
	"fmt"
)

type CurrencyProvider int

const (
	BOG       CurrencyProvider = 1
	BASIS     CurrencyProvider = 2
	CREDO     CurrencyProvider = 3
	CRYSTAL   CurrencyProvider = 4
	KURSI     CurrencyProvider = 5
	LIBERTY   CurrencyProvider = 6
	PROCREDIT CurrencyProvider = 7
	RICO      CurrencyProvider = 8
	SWISS     CurrencyProvider = 9
	TBC       CurrencyProvider = 10
	TERA      CurrencyProvider = 11
)

type CurrencyCode int

const (
	USD CurrencyCode = 1
	EUR CurrencyCode = 2
	GBP CurrencyCode = 3
)

var currencyNames = map[CurrencyProvider]string{
	BOG:       "BOG",
	BASIS:     "Basis",
	CREDO:     "Credo",
	CRYSTAL:   "Crystal",
	KURSI:     "Kursi",
	LIBERTY:   "Liberty",
	PROCREDIT: "Pro Credit",
	RICO:      "Rico",
	SWISS:     "Swiss Capital",
	TBC:       "TBC",
	TERA:      "Tera",
}

var currencyCodes = map[CurrencyCode]string{
	USD: "USD",
	EUR: "EUR",
	GBP: "GBP",
}

func (p CurrencyProvider) String() string {
	if s, ok := currencyNames[p]; ok {
		return s
	}
	return fmt.Sprintf("CurrencyProvider(%d)", int(p))
}

func (p CurrencyCode) String() string {
	if s, ok := currencyCodes[p]; ok {
		return s
	}
	return fmt.Sprintf("CurrencyCode(%d)", int(p))
}

func (p CurrencyProvider) Logo() string {
	var imagePath = map[CurrencyProvider]string{
		BOG:       "static/img/logos/currency/bog.png",
		BASIS:     "static/img/logos/currency/basis.png",
		CREDO:     "static/img/logos/currency/credo.png",
		CRYSTAL:   "static/img/logos/currency/crystal.png",
		KURSI:     "static/img/logos/currency/kursi.png",
		LIBERTY:   "static/img/logos/currency/liberty.png",
		PROCREDIT: "static/img/logos/currency/pro-credit.png",
		RICO:      "static/img/logos/currency/rico.png",
		SWISS:     "static/img/logos/currency/swiss-capital.png",
		TBC:       "static/img/logos/currency/tbc.png",
		TERA:      "static/img/logos/currency/tera.png",
	}

	if s, ok := imagePath[p]; ok {
		return s
	}
	return fmt.Sprintf("CurrencyProvider(%d)", int(p))
}
