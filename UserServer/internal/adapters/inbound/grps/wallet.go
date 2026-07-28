package grps

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grps/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grps/validate"
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
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	if err := validate.Deposit(req); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	amount, err := mapper.ToProtoMoney(req.Amount)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	success, balances, err := s.wallet.Deposit(ctx, req.UserId, req.Asset, amount)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.DepositResponse{
		Success: success,
		Balance: mapper.ToProtoBalance(balances),
	}, nil
}

func (s *ServerApi) GetBalances(
	ctx context.Context,
	req *pb.EmptyRequest,
) (
	*pb.UserBalancesInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	balancesResponse, err := s.wallet.GetBalances(ctx)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}
	return &pb.UserBalancesInfoResponse{
		Balances: mapper.ToProtoBalances(balancesResponse),
	}, nil
}
