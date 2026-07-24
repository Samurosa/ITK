package validate

import (
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
