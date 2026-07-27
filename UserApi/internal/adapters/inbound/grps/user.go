package grps

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grps/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grps/validate"
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerApi) GetUser(
	ctx context.Context,
	req *pb.EmptyRequest,
) (
	*pb.UserInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	user, err := s.user.GetUser(ctx)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.UserInfoResponse{
		UserId:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      mapper.ToProtoRole(user.Role),
		CreatedAt: timestamppb.New(user.CreateTime),
		UpdatedAt: timestamppb.New(user.UpdateTime),
	}, nil
}

func (s *ServerApi) UpdateUserInfo(
	ctx context.Context,
	req *pb.UpdateUserInfoRequest,
) (
	*pb.UpdateUserInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
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
		name,
		email,
	)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.UpdateUserInfoResponse{
		Updated:       updated,
		UpdatedInfoAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *ServerApi) DeleteUser(
	ctx context.Context,
	req *pb.EmptyRequest,
) (
	*pb.DeleteUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	success, deletedUserAt, err := s.user.DeleteUser(ctx)
	if err != nil {
		return nil, mapper.ToGRPC(err)
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
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	if err := validate.Password(req.NewPassword); err != nil {
		return nil, err
	}

	isSuccess, userPasswordChangedAt, err := s.user.ChangePassword(ctx,
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}
	return &pb.ChangeUserResponse{
		Success:               isSuccess,
		UserPasswordChangedAt: timestamppb.New(userPasswordChangedAt),
	}, nil
}
