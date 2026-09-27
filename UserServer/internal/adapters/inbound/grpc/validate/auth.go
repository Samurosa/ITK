package validate

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"strings"
	"unicode"
)

func ComparePasswords(oldPassword, newPassword string) error {
	if strings.Compare(oldPassword, newPassword) == 0 {
		return corerrors.ErrPasswordsMatch
	}

	return nil
}

func Password(password string) error {
	if password == "" {
		return corerrors.ErrPasswordEmpty
	}

	var hasUpper, hasLower, hasDigit bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper {
		return corerrors.ErrPasswordWrongUpperSymbol
	}

	if !hasLower {
		return corerrors.ErrPasswordWrongLowerSymbol
	}

	if !hasDigit {
		return corerrors.ErrPasswordWrongDigitSymbol
	}

	return nil
}
