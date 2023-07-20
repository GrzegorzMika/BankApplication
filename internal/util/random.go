package util

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ"
)

var currencies = []string{"USD", "EUR", "GBP", "JPY", "AUD", "CAD", "CHF"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFloat returns a random float between min and max.
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomString returns a random string of length length.
func RandomString(length int) string {
	var sb strings.Builder
	k := len(letters)
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rand.Intn(k)])
	}
	return sb.String()
}

// RandomOwner returns a random owner name.
func RandomOwner() string {
	caser := cases.Title(language.English)
	return caser.String(RandomString(10))
}

// RandomMoney returns a random money amount.
func RandomMoney() float64 {
	return math.Round(RandomFloat(0, 10000)*100) / 100
}

// RandomCurrency returns a random currency.
func RandomCurrency() string {
	return currencies[rand.Intn(len(currencies))]
}
