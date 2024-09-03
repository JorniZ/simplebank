package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	randInt := RandomInt(20, 25)
	require.GreaterOrEqual(t, randInt, int64(20))
	require.LessOrEqual(t, randInt, int64(25))
}

func TestRandomString(t *testing.T) {
	randString := RandomString(5)
	require.Len(t, randString, 5)

	randString = RandomString(50)
	require.Len(t, randString, 50)
}

func TestRandomOwner(t *testing.T) {
	owner := RandomOwner()
	require.NotEmpty(t, owner)
}

func TestRandomMoney(t *testing.T) {
	money := RandomMoney()
	require.NotEmpty(t, money)
	require.LessOrEqual(t, money, int64(1000))
}

func TestRandomCurrency(t *testing.T) {
	currency := RandomCurrency()
	require.NotEmpty(t, currency)
	require.Contains(t, supportedCurrency, currency)
}

func TestRandomEmail(t *testing.T) {
	email := RandomEmail()
	require.NotEmpty(t, email)
	require.Contains(t, email, "@email.com")
}
