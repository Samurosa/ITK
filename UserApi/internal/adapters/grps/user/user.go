package user

import (
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverApi) GetUser(
	ctx context.Context,
	req *pb.UserIDRequest,
) (
	*pb.UserInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}

	name, email, balances, role, createdAt, updatedAt,
		err := s.user.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get user")
	}

	return &pb.UserInfoResponse{
		Name:      name,
		Email:     email,
		Balances:  ToProtoBalances(balances),
		Role:      ToProtoRole(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *serverApi) UpdateUserInfo(
	ctx context.Context,
	req *pb.UpdateUserInfoRequest,
) (
	*pb.UpdateUserInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}

	name := ""
	email := ""
	if req.Name != nil {
		name = req.GetName()
	}
	if req.Email != nil {
		email = req.GetEmail()
	}

	updated, updatedAt, err := s.user.UpdateUserInfo(
		ctx,
		req.Id,
		name,
		email,
	)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.UpdateUserInfoResponse{
		Updated:       updated,
		UpdatedInfoAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *serverApi) DeleteUser(
	ctx context.Context,
	req *pb.UserIDRequest,
) (
	*pb.DeleteUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}

	success, deletedUserAt, err := s.user.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "failed to delete user")
	}
	return &pb.DeleteUserResponse{
		Success:       success,
		DeletedUserAt: timestamppb.New(deletedUserAt),
	}, nil
}

func (s *serverApi) ChangePassword(
	ctx context.Context,
	req *pb.ChangeUserRequest,
) (
	*pb.ChangeUserResponse,
	error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isSuccess, userPasswordChangedAt := s.user.ChangePassword(ctx, req.Id,
		req.OldPassword,
		req.NewPassword,
	)

	return &pb.ChangeUserResponse{
		Success:               isSuccess,
		UserPasswordChangedAt: timestamppb.New(userPasswordChangedAt),
	}, nil
}
