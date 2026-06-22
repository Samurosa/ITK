package user

import (
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverApi) Registration(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (
	*pb.RegisterUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateRegistration(req); err != nil {
		return nil, err
	}

	id, tokens, createdAt, err := s.user.Registration(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "failed to register user")
	}

	return &pb.RegisterUserResponse{
		Id:        id,
		Tokens:    ToProtoTokens(tokens),
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *serverApi) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (
	*pb.TokenPairResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if err := ValidateLogin(req); err != nil {
		return nil, err
	}

	tokens, err := s.user.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to authorize user")
	}

	return ToProtoTokens(tokens), nil
}

func (s *serverApi) Logout(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (
	*pb.LogoutResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	success, loggedOutAt, err := s.user.Logout(ctx, req.RefreshToken, req.DeviceId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to unauthenticated user")
	}

	return &pb.LogoutResponse{
		Success:     success,
		LoggedOutAt: timestamppb.New(loggedOutAt),
	}, nil
}

func (s *serverApi) RefreshToken(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (
	*pb.TokenPairResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	tokens, err := s.user.RefreshToken(ctx, req.RefreshToken, req.DeviceId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}
	return ToProtoTokens(tokens), nil
}
