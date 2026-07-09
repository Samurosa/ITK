package user

import (
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func (s *ServerApi) Deposit(
	ctx context.Context,
	req *pb.DepositRequest,
) (
	*pb.DepositResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, ToGRPC(err)
	}
	if err := ValidateUserId(req.UserId); err != nil {
		return nil, ToGRPC(err)
	}
	if err := ValidateDepositRequest(req); err != nil {
		return nil, ToGRPC(err)
	}
	amount, err := ToProtoMoney(req.Amount)
	if err != nil {
		return nil, ToGRPC(err)
	}
	success, balances, err := s.wallet.Deposit(ctx, req.UserId, req.Asset, amount)
	if err != nil {
		return nil, ToGRPC(err)
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
		return nil, ToGRPC(err)
	}

	if err := ValidateUserId(req.UserId); err != nil {
		return nil, ToGRPC(err)
	}

	balancesResponse, err := s.wallet.GetBalances(ctx, req.UserId)
	if err != nil {
		return nil, ToGRPC(err)
	}
	return &pb.UserBalancesInfoResponse{
		Balances: ToProtoBalances(balancesResponse),
	}, nil
}
