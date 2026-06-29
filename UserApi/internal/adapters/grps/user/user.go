package user

import (
	"ITK_Code/m/v2/internal/core/user/models"
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerApi) GetUser(
	ctx context.Context,
	req *pb.UserIDRequest,
) (
	*pb.UserInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.UserId); err != nil {
		return nil, err
	}

	user, err := s.user.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get user")
	}

	return &pb.UserInfoResponse{
		UserId:    req.UserId,
		Name:      user.Name,
		Email:     user.Email,
		Balances:  ToProtoBalances(user.Balances),
		Role:      ToProtoRole(user.Role),
		CreatedAt: timestamppb.New(user.CreateTime),
		UpdatedAt: timestamppb.New(user.UpdateTime),
	}, nil
}

func (s *ServerApi) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServerApi) UpdateUserInfo(
	ctx context.Context,
	req *pb.UpdateUserInfoRequest,
) (
	*pb.UpdateUserInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.UserId); err != nil {
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
		req.UserId,
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

func (s *ServerApi) DeleteUser(
	ctx context.Context,
	req *pb.UserIDRequest,
) (
	*pb.DeleteUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.UserId); err != nil {
		return nil, err
	}

	success, deletedUserAt, err := s.user.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "failed to delete user")
	}
	return &pb.DeleteUserResponse{
		Success:       success,
		DeletedUserAt: timestamppb.New(deletedUserAt),
	}, nil
}

func (s *ServerApi) ChangePassword(
	ctx context.Context,
	req *pb.ChangeUserRequest,
) (
	*pb.ChangeUserResponse,
	error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isSuccess, userPasswordChangedAt := s.user.ChangePassword(ctx, req.UserId,
		req.OldPassword,
		req.NewPassword,
	)

	return &pb.ChangeUserResponse{
		Success:               isSuccess,
		UserPasswordChangedAt: timestamppb.New(userPasswordChangedAt),
	}, nil
}
