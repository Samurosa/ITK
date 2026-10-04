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

func (s *UserServer) Registration(ctx context.Context,
	req *pb.RegisterUserRequest,
) (
	*pb.RegisterUserResponse,
	error,
) {
	log := logging.FromContext(ctx).Named("Registration")

	if err := req.Validate(); err != nil {
		log.Debug("invalid registration request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	if err := validate.Password(req.GetPassword()); err != nil {
		log.Debug("registration password rejected by policy", zap.Error(err))
		return nil, mapper.ToGRPC(err)
	}

	requestContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		log.Error("registration request context is missing")
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	id, createdAt, err := s.auth.Registration(ctx,
		req.Email,
		req.Password,
		req.Name,
		requestContext.Metadata.ClientIP,
		requestContext.Metadata.DeviceID,
	)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.RegisterUserResponse{
		UserId:    id,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *UserServer) Login(ctx context.Context,
	req *pb.LoginRequest,
) (
	*pb.TokenPairResponse,
	error,
) {
	log := logging.FromContext(ctx).Named("Login")

	if err := req.Validate(); err != nil {
		log.Debug("invalid login request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	requestContext, ok := sharedContext.GetRequestContext(ctx)
	if !ok {
		log.Error("login request context is missing")
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	tokens, err := s.auth.Login(ctx,
		req.Email,
		req.Password,
		requestContext.Metadata.ClientIP,
		requestContext.Metadata.DeviceID,
	)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return mapper.ToProtoTokens(tokens), nil
}

func (s *UserServer) Logout(ctx context.Context,
	req *pb.LogoutRequest,
) (
	*emptypb.Empty,
	error,
) {
	log := logging.FromContext(ctx).Named("Logout")

	if err := req.Validate(); err != nil {
		log.Debug("invalid logout request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	jti, ok := sharedContext.JTI(ctx)
	if !ok {
		log.Error("authenticated session context is missing")
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	err := s.auth.Logout(ctx, jti, req.RefreshToken)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *UserServer) LogoutAllDevices(ctx context.Context,
	_ *emptypb.Empty,
) (
	*emptypb.Empty,
	error,
) {
	jti, ok := sharedContext.JTI(ctx)
	if !ok {
		logging.FromContext(ctx).Named("LogoutAllDevices").Error("authenticated session context is missing")
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	err := s.auth.LogoutAllDevices(ctx, jti)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *UserServer) RefreshToken(ctx context.Context,
	req *pb.RefreshTokenRequest,
) (
	*pb.TokenPairResponse,
	error,
) {
	log := logging.FromContext(ctx).Named("RefreshToken")

	if err := req.Validate(); err != nil {
		log.Debug("invalid refresh token request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	tokens, err := s.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}
	return mapper.ToProtoTokens(tokens), nil
}
