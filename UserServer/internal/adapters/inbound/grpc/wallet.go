package grpc

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/validate"
	requestContext "ITK_Code/m/v2/internal/core/context"
	"ITK_Code/m/v2/internal/core/errors"
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"go.uber.org/zap"
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

	amount, err := mapper.WithProtoMoney(req.Amount)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	if err := validate.Money(amount); err != nil {
		return nil, mapper.ToGRPC(err)
	}

	balances, err := s.wallet.Deposit(ctx, req.UserId, req.Asset, amount, req.IdempotencyKey)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.DepositResponse{
		Balance: mapper.ToProtoBalance(balances),
	}, nil
}

func (s *ServerApi) GetBalances(
	ctx context.Context,
	_ *pb.Empty,
) (
	*pb.UserBalancesInfoResponse,
	error,
) {
	log := s.log.Named("GetBalances")

	id, err := requestContext.UserID(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return nil, errors.ErrInvalidContext
	}
	log.Debug("user id from context", zap.String("id", id))

	balancesResponse, err := s.wallet.GetBalances(ctx, id)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.UserBalancesInfoResponse{
		Balances: mapper.ToProtoBalances(balancesResponse),
	}, nil
}
