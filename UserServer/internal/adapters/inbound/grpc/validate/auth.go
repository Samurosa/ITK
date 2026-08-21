package validate

import (
	"ITK_Code/m/v2/internal/core/errors"
	"strings"
	"unicode"
)

func ComparePasswords(oldPassword, newPassword string) error {
	if strings.Compare(oldPassword, newPassword) == 0 {
		return errors.ErrPasswordsMatch
	}

	return nil
}

func Password(password string) error {
	if password == "" {
		return errors.ErrPasswordEmpty
	}

	var hasUpper, hasDigit bool

	for _, char := range password {

		switch {
		case unicode.IsUpper(char):
			hasUpper = true

		case unicode.IsLower(char):
			hasDigit = true
		}
	}
	if !hasUpper {
		return errors.ErrPasswordWrongUpperSymbol
	}

	if !hasDigit {
		return errors.ErrPasswordWrongDigitSymbol
	}

	return nil
}
