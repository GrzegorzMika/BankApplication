package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsSupportedCurrency checks if the given currency is supported.
//
// Supported currencies are:
//   - USD
//   - EUR
//   - CAD
//
// Returns true if the given currency is supported, false otherwise.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
