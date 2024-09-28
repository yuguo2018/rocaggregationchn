package core

import (
	"errors"
	"regexp"
)

var printable7BitAscii = regexp.MustCompile("^[A-Za-z0-9!\"#$%&'()*+,\\-./:;<=>?@[\\]^_`{|}~ ]+$")

// ValidatePasswordFormat returns an error if the password is too short, or consists of characters
// outside the range of the printable 7bit ascii set
func ValidatePasswordFormat(password string) error {
	if len(password) < 10 {
		return errors.New("password too short (<10 characters)")
	}
	if !printable7BitAscii.MatchString(password) {
		return errors.New("password contains invalid characters - only 7bit printable ascii allowed")
	}
	return nil
}
