package user

import (
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerApi) Deposit(
	ctx context.Context,
	req *pb.DepositRequest,
) (
	*pb.DepositResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := ValidateUserId(req.UserId); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := ValidateDepositRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	amount, err := ToProtoMoney(req.Amount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	success, balances, err := s.wallet.Deposit(ctx, req.UserId, req.Asset, amount)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to deposit funds")
	}

	return &pb.DepositResponse{
		Success: success,
		Balance: ToProtoBalance(balances),
	}, nil
}

func (s *ServerApi) GetBalances(
	ctx context.Context,
	req *pb.UserIDRequest,
) (
	*pb.UserBalancesInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := ValidateUserId(req.UserId); err != nil {
		return nil, status.Error(codes.NotFound, "invalid user id")
	}

	balancesResponse, err := s.wallet.GetBalances(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get wallet")
	}
	return &pb.UserBalancesInfoResponse{
		Balances: ToProtoBalances(balancesResponse),
	}, nil
}
