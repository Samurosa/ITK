package grpc

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/validate"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
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
		logging.FromContext(ctx).Named("GetUser").Error("authenticated user context is missing")
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
	log := logging.FromContext(ctx).Named("UpdateUserInfo")
	if err := req.Validate(); err != nil {
		log.Debug("invalid user update request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	name := req.GetName()

	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		log.Error("authenticated user context is missing")
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
		logging.FromContext(ctx).Named("DeleteUser").Error("authenticated user context is missing")
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
	log := logging.FromContext(ctx).Named("ChangePassword")

	if err := req.Validate(); err != nil {
		log.Debug("invalid password change request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		log.Error("authenticated user context is missing")
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}
	if err := validate.ComparePasswords(req.GetOldPassword(), req.GetNewPassword()); err != nil {
		log.Debug("password change rejected: new password matches current password")
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Password(req.NewPassword); err != nil {
		log.Debug("new password rejected by policy", zap.Error(err))
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
