package grps

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grps/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grps/validate"
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ServerApi) Registration(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (
	*pb.RegisterUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Registration(req); err != nil {
		return nil, err
	}

	if err := validate.Password(req.GetPassword()); err != nil {
		return nil, err
	}

	id, createdAt, err := s.auth.Registration(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.RegisterUserResponse{
		UserId:    id,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *ServerApi) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (
	*pb.TokenPairResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Login(req); err != nil {
		return nil, err
	}

	tokens, err := s.auth.Login(ctx, req.Email, req.Password, req.DeviceId)
	if err != nil {
		s.log.Error("failed to login", zap.Error(err))
		return nil, mapper.ToGRPC(err)
	}

	return mapper.ToProtoTokens(tokens), nil
}

func (s *ServerApi) Logout(
	ctx context.Context,
	req *pb.LogoutRequest,
) (
	*pb.LogoutResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	success, loggedOutAt, err := s.auth.Logout(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.LogoutResponse{
		Success:     success,
		LoggedOutAt: timestamppb.New(loggedOutAt),
	}, nil
}

func (s *ServerApi) LogoutAllDevices(
	ctx context.Context,
	req *pb.LogoutAllRequest,
) (
	*pb.LogoutResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	success, loggedOutAt, err := s.auth.LogoutAllDevices(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.LogoutResponse{
		Success:     success,
		LoggedOutAt: timestamppb.New(loggedOutAt),
	}, nil
}

func (s *ServerApi) RefreshToken(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (
	*pb.TokenPairResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	tokens, err := s.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}
	return mapper.ToProtoTokens(tokens), nil
}
