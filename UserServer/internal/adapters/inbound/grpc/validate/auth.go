package validate

import (
	"ITK_Code/m/v2/internal/core/errorcore"
	"unicode"
)

func ComparePasswords(oldPassword, newPassword string) error {
	if oldPassword == newPassword {
		return errorcore.ErrPasswordsMatch
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
		return errorcore.ErrPasswordWrongUpperSymbol
	}

	if !hasLower {
		return errorcore.ErrPasswordWrongLowerSymbol
	}

	if !hasDigit {
		return errorcore.ErrPasswordWrongDigitSymbol
	}

	return nil
}
