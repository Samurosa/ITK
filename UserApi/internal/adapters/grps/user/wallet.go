package user

import (
	"ITK_Code/m/v2/internal/core/user/models"
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

	success, balance, err := s.user.Deposit(
		ctx,
		req.UserId,
		req.Asset,
		models.Money{
			Currency: req.Amount.Currency,
			Units:    req.Amount.Units,
			Nanos:    req.Amount.Nanos,
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to deposit funds")
	}
	return &pb.DepositResponse{
		Success: success,
		Balance: ToProtoBalance(balance),
	}, nil
}

func (s *ServerApi) GetBalance(
	ctx context.Context,
	req *pb.UserIDRequest,
) (
	*pb.UserBalanceInfoResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := ValidateUserId(req.UserId); err != nil {
		return nil, status.Error(codes.NotFound, "invalid user id")
	}

	balancesResponse, err := s.user.GetBalance(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get balance")
	}
	return &pb.UserBalanceInfoResponse{
		Balances: ToProtoBalances(balancesResponse),
	}, nil
}
