package val

import (
	"fmt"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFirstName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
	isValidLastname = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int ) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength,maxLength)
	}
	return nil
}

func ValidateFirstname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFirstName(value) {
		 return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidateLastname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidLastname(value) {
		 return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		 return fmt.Errorf("must contain only letters, digits, or underscrore")
	}
	return nil
}

func ValidatePassword(value string) error{
	return ValidateString(value, 8, 100)
}