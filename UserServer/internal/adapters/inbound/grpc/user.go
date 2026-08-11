package grpc

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/validate"
	requestContext "ITK_Code/m/v2/internal/core/context"
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerApi) GetUser(ctx context.Context,
	_ *pb.Empty,
) (
	*pb.UserInfoResponse,
	error,
) {
	log := s.log.Named("GetUser")

	userID, err := requestContext.UserID(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return &pb.UserInfoResponse{}, mapper.ToGRPC(err)
	}
	log.Debug("user id from context", zap.String("id", userID))

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

func (s *ServerApi) UpdateUserInfo(ctx context.Context,
	req *pb.UpdateUserInfoRequest,
) (
	*pb.Empty,
	error,
) {
	log := s.log.Named("UpdateUserInfo")
	if err := req.Validate(); err != nil {
		log.Error("invalid request", zap.Error(err))
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

	userID, err := requestContext.UserID(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return nil, mapper.ToGRPC(err)
	}
	log.Debug("user id from context", zap.String("id", userID))

	err = s.user.UpdateUserInfo(
		ctx,
		userID,
		name,
		email,
	)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.Empty{}, nil
}

func (s *ServerApi) DeleteUser(ctx context.Context,
	_ *pb.Empty,
) (
	*pb.Empty,
	error,
) {
	log := s.log.Named("DeleteUser")

	userID, err := requestContext.UserID(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return nil, mapper.ToGRPC(err)
	}
	log.Debug("user id from context", zap.String("id", userID))

	err = s.user.DeleteUser(ctx, userID)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.Empty{}, nil
}

func (s *ServerApi) ChangePassword(ctx context.Context,
	req *pb.ChangeUserRequest,
) (
	*pb.Empty,
	error,
) {
	log := s.log.Named("ChangePassword")

	if err := req.Validate(); err != nil {
		log.Error("invalid request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	userID, err := requestContext.UserID(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return nil, mapper.ToGRPC(err)
	}
	log.Debug("user id from context", zap.String("id", userID))

	if err := validate.ComparePasswords(req.GetOldPassword(), req.GetNewPassword()); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Password(req.NewPassword); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	err = s.user.ChangePassword(ctx, userID, req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.Empty{}, nil
}
