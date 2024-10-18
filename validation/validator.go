package validation

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength, MaxLength int) error {
	valueLength := len(value)
	if valueLength < minLength || valueLength > MaxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, MaxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 20); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return errors.New("must contain only lowercase letters, digits or underscore")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return errors.New("is not a valid email address")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 20); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return errors.New("must contain only letters or spaces")
	}

	return nil
}
