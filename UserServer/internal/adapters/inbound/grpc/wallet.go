package grpc

import (
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserServer) Deposit(
	_ context.Context,
	req *pb.DepositRequest,
) (
	*pb.DepositResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	return nil, status.Error(codes.Unimplemented, "deposit is not implemented")
}

func (s *UserServer) GetBalances(
	_ context.Context,
	_ *emptypb.Empty,
) (
	*pb.UserBalancesInfoResponse,
	error,
) {

	return nil, status.Error(codes.Unimplemented, "balance lookup is not implemented")
}
