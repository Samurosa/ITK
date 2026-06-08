package user

import (
	"context"

	pb "github.com/Truncklin/exchange-contract/protobuf/gen/go/user"

	"google.golang.org/grpc"
)

type serverApi struct {
	pb.UnimplementedUserServiceServer
}

func NewServer(grps *grpc.Server) {
	pb.RegisterUserServiceServer(grps, &serverApi{})
}

func (s *serverApi) Register(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (
	*pb.RegisterUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &pb.RegisterUserResponse{}, nil
}

func (s *serverApi) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (
	*pb.GetUserResponse,
	error,
) {

	return &pb.GetUserResponse{}, nil
}

func (s *serverApi) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (
	*pb.UpdateUserResponse,
	error,
) {
	return &pb.UpdateUserResponse{}, nil
}

func (s *serverApi) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (
	*pb.DeleteUserResponse,
	error,
) {
	return &pb.DeleteUserResponse{}, nil
}

func (s *serverApi) Deposit(
	ctx context.Context,
	req *pb.DepositRequest,
) (
	*pb.DepositResponse,
	error,
) {
	return &pb.DepositResponse{}, nil
}

func (s *serverApi) Authorization(
	ctx context.Context,
	req *pb.AuthorizationRequest,
) (
	*pb.AuthorizationResponse,
	error,
) {
	return &pb.AuthorizationResponse{}, nil
}
