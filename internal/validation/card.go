package validation

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	errCardExpired           = errors.New("card is expired")
	errInvalidExpirationDate = errors.New("invalid expiration date")
	errInvalidMonthFormat    = errors.New("invalid month format")
	errInvalidYearFormat     = errors.New("invalid year format")
)

func ValidateCardNumber(number string) bool {
	var sum int
	var alternate bool

	// Iterate over the card number string in reverse
	for i := len(number) - 1; i >= 0; i-- {
		n := number[i]

		// Check if the character is a digit
		if !unicode.IsDigit(rune(n)) {
			return false
		}

		digit := int(n - '0')

		// Double every second digit
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	// Valid card numbers will have a sum divisible by 10
	return sum%10 == 0
}

// ValidateExpirationDate checks if the expiration date is in the format MM/YY and not in the past.
func ValidateExpirationDate(expirationDate string) error {
	// Split the date into month and year
	parts := strings.Split(expirationDate, "/")
	if len(parts) != 2 {
		return errInvalidExpirationDate
	}

	if len(parts[0]) != 2 {
		return errInvalidMonthFormat
	}
	if len(parts[1]) != 2 {
		return errInvalidYearFormat
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return errInvalidMonthFormat
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return errInvalidExpirationDate
	}

	// Create the time for the last day of the expiration month
	expirationTime := time.Date(year+2000, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)

	// Check if the expiration time is before the current time
	if expirationTime.Before(time.Now()) {
		return errCardExpired
	}
	return nil
}

// ValidateCVV checks if the CVV is 3 or 4 digits.
func ValidateCVV(cvv string) bool {
	if len(cvv) < 3 || len(cvv) > 4 {
		return false
	}

	for _, r := range cvv {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
