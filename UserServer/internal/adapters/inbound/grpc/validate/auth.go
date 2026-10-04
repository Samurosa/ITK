package validate

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"unicode"
)

func ComparePasswords(oldPassword, newPassword string) error {
	if oldPassword == newPassword {
		return corerrors.ErrPasswordsMatch
	}

	return nil
}

func Password(password string) error {
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
