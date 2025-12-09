package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	IsValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	IsValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(val string, minLen int, maxLen int) error {
	n := len(val)
	if n < minLen || n > maxLen {
		return fmt.Errorf("string length must be between %d and %d", minLen, maxLen)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !IsValidUsername(username) {
		return fmt.Errorf("must only contain lowercase letters, numbers, and underscores")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if err := ValidateString(fullName, 3, 100); err != nil {
		return err
	}
	if !IsValidFullName(fullName) {
		return fmt.Errorf("must only contain letters, numbers, and spaces")
	}
	return nil
}


func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("the email is not valid")
	}
	return nil
}