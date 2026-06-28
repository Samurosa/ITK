package user

import (
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateUserId(id string) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Id is required")
	}
	return nil
}

func ValidateDepositRequest(req *pb.DepositRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "UserId is required")
	}
	if req.GetAsset() == "" {
		return status.Error(codes.InvalidArgument, "Asset is required")
	}
	if req.GetAmount().Currency == "" {
		return status.Error(codes.InvalidArgument, "Amount is required")
	}

	return nil
}

func ValidateRegistration(req *pb.RegisterUserRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password is required")
	}
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "Password is required")
	}
	return nil
}

func ValidateLogin(req *pb.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password is required")
	}
	return nil
}
