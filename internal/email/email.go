package email

import (
	"errors"
	"regexp"
)

// Validate checks if the email is valid.
func Validate(email string) error {
	emailRegexp, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		return err
	}

	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}

	return nil
}
