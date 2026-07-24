package user

import (
	"unicode"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserIDEmpty   = status.Error(codes.InvalidArgument, "user id empty")
	ErrEmailEmpty    = status.Error(codes.InvalidArgument, "email empty")
	ErrUsernameEmpty = status.Error(codes.InvalidArgument, "username empty")
	ErrAssetEmpty    = status.Error(codes.InvalidArgument, "asset empty")
	ErrAmountEmpty   = status.Error(codes.InvalidArgument, "amount empty")

	ErrPasswordEmpty            = status.Error(codes.InvalidArgument, "password empty")
	ErrPasswordWrongUpperSymbol = status.Error(codes.InvalidArgument, "password wrong, upper symbol not found")
	ErrPasswordWrongDigitSymbol = status.Error(codes.InvalidArgument, "password wrong, digit not found")
)

func ValidateUserId(id string) error {
	if id == "" {
		return ErrUserIDEmpty
	}
	return nil
}

func ValidateDepositRequest(req *pb.DepositRequest) error {
	if req.GetUserId() == "" {
		return ErrUserIDEmpty
	}
	if req.GetAsset() == "" {
		return ErrAssetEmpty
	}
	if req.GetAmount().Currency == "" {
		return ErrAmountEmpty
	}

	return nil
}

func ValidateRegistration(req *pb.RegisterUserRequest) error {
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

func ValidateLogin(req *pb.LoginRequest) error {
	if req.GetEmail() == "" {
		return ErrEmailEmpty
	}
	if req.GetPassword() == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func ValidatePassword(password string) error {
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
