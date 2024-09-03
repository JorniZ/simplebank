package util

import "slices"

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

var (
	supportedCurrency = []string{USD, EUR, CAD}
)

func IsSupportedCurrency(currency string) bool {
	return slices.Contains(supportedCurrency, currency)
}
