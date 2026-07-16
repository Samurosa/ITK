package user

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"
	"errors"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoBalances(
	balance []wallet.Balance,
) []*pb.Balance {
	result := make(
		[]*pb.Balance,
		0,
		len(balance),
	)

	for _, b := range balance {
		result = append(result,
			&pb.Balance{
				Asset:     b.Asset,
				Available: b.Available.String(),
				Locked:    b.Locked.String(),
			},
		)
	}
	return result
}

func ToProtoBalance(balance wallet.Balance) *pb.Balance {
	return &pb.Balance{
		Asset:     balance.Asset,
		Available: balance.Available.String(),
		Locked:    balance.Locked.String(),
	}
}

func ToProtoRole(role user.Role) pb.Role {
	switch role {
	case user.UserRole:
		return pb.Role_ROLE_USER
	case user.GuestRole:
		return pb.Role_ROLE_GUEST
	case user.PremiumRole:
		return pb.Role_ROLE_PREMIUM
	case user.AdminRole:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

func ToProtoTokens(tokens auth.TokensModel) *pb.TokenPairResponse {
	return &pb.TokenPairResponse{
		AccessToken:     tokens.AccessToken,
		RefreshToken:    tokens.RefreshToken,
		AccessExpiresAt: timestamppb.New(tokens.AccessExpiresAt),
		//AccessCreatesAt:  timestamppb.New(tokens.AccessCreatedAt),
		RefreshExpiresAt: timestamppb.New(tokens.RefreshExpiresAt),
		//RefreshCreateAt:  timestamppb.New(tokens.RefreshCreatedAt),
	}
}

func ToProtoMoney(protoMoney *shared.Money) (wallet.Money, error) {
	if protoMoney == nil {
		return wallet.Money{}, errors.New("money is nil")
	}

	amount, err := decimal.NewFromString(protoMoney.Amount)
	if err != nil {
		return wallet.Money{}, err
	}

	return wallet.Money{
		Currency: protoMoney.Currency,
		Amount:   amount,
	}, nil
}

func ToGRPC(err error) error {
	switch {
	case errors.Is(err, wallet.ErrBalanceNotFound):
		return status.Error(codes.NotFound, "balance not found")
	case errors.Is(err, wallet.ErrCreateNewBalance):
		return status.Error(codes.Internal, "failed to create balance")
	case errors.Is(err, wallet.ErrSaveBalance):
		return status.Error(codes.Internal, "failed to save balance")

	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, user.ErrComparePassword):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, user.ErrUpdateUser):
		return status.Error(codes.Internal, "failed to update user")

	case errors.Is(err, auth.ErrRefreshExpired):
		return status.Error(codes.Unauthenticated, "refresh token expired")
	case errors.Is(err, auth.ErrGenerateToken):
		return status.Error(codes.Internal, "failed to generate token")
	case errors.Is(err, auth.ErrInvalidToken):
		return status.Error(codes.InvalidArgument, "invalid token")
	case errors.Is(err, auth.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session expired")
	case errors.Is(err, auth.ErrSessionNotFound):
		return status.Error(codes.Unauthenticated, "session not found")
	case errors.Is(err, auth.ErrInvalidContext):
		return status.Error(codes.Internal, "invalid context")
	case errors.Is(err, auth.Unauthorized):
		return status.Error(codes.Unauthenticated, "incorrect login or password")
	case errors.Is(err, auth.ErrNoAccess):
		return status.Error(codes.PermissionDenied, "no access")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
