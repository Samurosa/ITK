package user

import (
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerApi) GetUser(
	ctx context.Context,
	req *pb.EmptyRequest,
) (
	*pb.UserInfoResponse,
	error,
) {
	_ = req
	user, err := s.user.GetUser(ctx)
	if err != nil {
		return nil, ToGRPC(err)
	}

	return &pb.UserInfoResponse{
		UserId:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      ToProtoRole(user.Role),
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
		return nil, ToGRPC(err)
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
		return nil, ToGRPC(err)
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
		return nil, ToGRPC(err)
	}

	success, deletedUserAt, err := s.user.DeleteUser(ctx)
	if err != nil {
		return nil, ToGRPC(err)
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
		return nil, ToGRPC(err)
	}

	if err := ValidatePassword(req.NewPassword); err != nil {
		return nil, err
	}

	isSuccess, userPasswordChangedAt, err := s.user.ChangePassword(ctx,
		req.OldPassword,
		req.NewPassword,
	)
	if err != nil {
		return nil, ToGRPC(err)
	}
	return &pb.ChangeUserResponse{
		Success:               isSuccess,
		UserPasswordChangedAt: timestamppb.New(userPasswordChangedAt),
	}, nil
}
