package validate

import (
	"unicode"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func Registration(req *pb.RegisterUserRequest) error {
	if req.GetEmail() == "" {
		return ErrEmailEmpty
	}
	if req.GetPassword() == "" {
		return ErrPasswordEmpty
	}
	if req.GetName() == "" {
		return ErrUsernameEmpty
	}
	return nil
}

func Login(req *pb.LoginRequest) error {
	if req.GetEmail() == "" {
		return ErrEmailEmpty
	}
	if req.GetPassword() == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func Password(password string) error {
	if password == "" {
		return ErrPasswordEmpty
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
		return ErrPasswordWrongUpperSymbol
	}

	if !hasDigit {
		return ErrPasswordWrongDigitSymbol
	}

	return nil
}
