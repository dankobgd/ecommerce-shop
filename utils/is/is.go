package is

import "encoding/json"

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

// ValidJSON checks if string is valid json
func ValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
