package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCardNumber(t *testing.T) {
	tests := []struct {
		number string
		valid  bool
	}{
		{"4111111111111111", true},     // Valid Visa
		{"5500000000000004", true},     // Valid MasterCard
		{"340000000000009", true},      // Valid American Express
		{"1234567890123456", false},    // Invalid card number
		{"4111-1111-1111-1111", false}, // Invalid format (non-digit characters)
	}

	for _, tt := range tests {
		got := ValidateCardNumber(tt.number)
		assert.Equalf(t, tt.valid, got, tt.number)
	}
}

func TestValidateExpirationDate(t *testing.T) {
	tests := []struct {
		date string
		err  error
	}{
		{"09/27", nil},                      // Valid future date
		{"12/99", nil},                      // Valid far future date
		{"01/20", errCardExpired},           // Expired date
		{"13/25", errInvalidMonthFormat},    // Invalid month
		{"00/25", errInvalidMonthFormat},    // Invalid month
		{"09/2", errInvalidYearFormat},      // Invalid year format
		{"09/2024", errInvalidYearFormat},   // Invalid format
		{"09-24", errInvalidExpirationDate}, // Invalid separator
	}

	for _, tt := range tests {
		got := ValidateExpirationDate(tt.date)
		assert.Equalf(t, tt.err, got, tt.date)
	}
}

func TestValidateCVV(t *testing.T) {
	tests := []struct {
		cvv   string
		valid bool
	}{
		{"123", true},    // Valid 3-digit CVV
		{"1234", true},   // Valid 4-digit CVV
		{"12", false},    // Too short
		{"12345", false}, // Too long
		{"12a", false},   // Non-digit characters
		{"", false},      // Empty string
	}

	for _, tt := range tests {
		got := ValidateCVV(tt.cvv)
		assert.Equalf(t, tt.valid, got, tt.cvv)
	}
}
