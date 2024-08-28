package util

import "slices"

var (
	supportedCurrency = []string{"EUR", "USD", "CAD"}
)

func IsSupportedCurrency(currency string) bool {
	return slices.Contains(supportedCurrency, currency)
}
