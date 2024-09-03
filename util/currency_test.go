package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedCurrency(t *testing.T) {
	validCurrency := IsSupportedCurrency("EUR")
	require.True(t, validCurrency)

	invalidCurrency := IsSupportedCurrency("ABC")
	require.False(t, invalidCurrency)
}
