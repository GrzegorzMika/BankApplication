package util

import "testing"

func TestIsSupportedCurrency_ValidCurrencies(t *testing.T) {
	tests := []struct {
		currency string
		expected bool
	}{
		{USD, true},
		{EUR, true},
		{CAD, true},
	}

	for _, test := range tests {
		actual := IsSupportedCurrency(test.currency)
		if actual != test.expected {
			t.Errorf("Expected IsSupportedCurrency(%s) to be %t, but it was %t", test.currency, test.expected, actual)
		}
	}
}

func TestIsSupportedCurrency_InvalidCurrency(t *testing.T) {
	tests := []struct {
		currency string
		expected bool
	}{
		{"INR", false},
		{"JPY", false},
	}

	for _, test := range tests {
		actual := IsSupportedCurrency(test.currency)
		if actual != test.expected {
			t.Errorf("Expected IsSupportedCurrency(%s) to be %t, but it was %t", test.currency, test.expected, actual)
		}
	}
}
