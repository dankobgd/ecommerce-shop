package is

// ValidEmail checks if the email format is valid
func ValidEmail(email string) bool {
	if err := ValidateFormat(email); err != nil {
		return false
	}
	return true
}

// ValidEmailAndMX checks if the email format is valid and MX record exists
func ValidEmailAndMX(email string) bool {
	if err := ValidateHost(email); err != nil {
		return false
	}
	return true
}
