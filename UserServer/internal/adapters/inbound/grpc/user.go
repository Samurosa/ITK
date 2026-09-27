package grpc

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/validate"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UserServer) GetUser(ctx context.Context,
	_ *emptypb.Empty,
) (
	*pb.UserInfoResponse,
	error,
) {
	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	user, err := s.user.GetUser(ctx, userID)
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

func (s *UserServer) UpdateUserInfo(ctx context.Context,
	req *pb.UpdateUserInfoRequest,
) (
	*emptypb.Empty,
	error,
) {
	log := s.log.Named("UpdateUserInfo")
	if err := req.Validate(); err != nil {
		log.Error("invalid request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	name := ""
	if req.Name != nil {
		name = req.GetName()
	}

	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	if err := s.user.UpdateUserInfo(
		ctx,
		userID,
		name,
	); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context,
	_ *emptypb.Empty,
) (
	*emptypb.Empty,
	error,
) {
	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}
	if err := s.user.DeleteUser(ctx, userID); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *UserServer) ChangePassword(ctx context.Context,
	req *pb.ChangePasswordRequest,
) (
	*emptypb.Empty,
	error,
) {
	log := s.log.Named("ChangePassword")

	if err := req.Validate(); err != nil {
		log.Error("invalid request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}
	if err := validate.ComparePasswords(req.GetOldPassword(), req.GetNewPassword()); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Password(req.NewPassword); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := s.user.ChangePassword(ctx,
		userID,
		req.GetOldPassword(),
		req.GetNewPassword(),
	); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}
