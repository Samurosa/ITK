package grpc

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/validate"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserServer) Deposit(
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

	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	amount, err := mapper.WithProtoMoney(req.Amount)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Money(amount); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	balances, err := s.wallet.Deposit(ctx, userID, amount, req.IdempotencyKey)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.DepositResponse{
		Balance: mapper.ToProtoBalance(balances),
	}, nil
}

func (s *UserServer) GetBalances(
	ctx context.Context,
	_ *emptypb.Empty,
) (
	*pb.UserBalancesInfoResponse,
	error,
) {
	userID, ok := sharedContext.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "context not provided")
	}

	balancesResponse, err := s.wallet.GetBalances(ctx, userID)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.UserBalancesInfoResponse{
		Balances: mapper.ToProtoBalances(balancesResponse),
	}, nil
}
