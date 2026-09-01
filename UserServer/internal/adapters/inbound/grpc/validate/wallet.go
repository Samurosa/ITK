package validate

import (
	"ITK_Code/m/v2/internal/core/coreErrors"
	"ITK_Code/m/v2/internal/core/wallet"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func Deposit(req *pb.DepositRequest) error {
	if req.GetUserId() == "" {
		return coreErrors.ErrUserIDEmpty
	}
	if req.GetAsset() == "" {
		return coreErrors.ErrAssetEmpty
	}
	if req.GetAmount() == nil {
		return coreErrors.ErrAmountEmpty
	}
	if req.GetAmount().Currency == "" {
		return coreErrors.ErrAmountEmpty
	}
	if req.GetAsset() != req.GetAmount().Currency {
		return coreErrors.ErrInvalidAsset
	}

	return nil
}

func Money(money wallet.Money) error {
	if money.Amount.IsZero() {
		return coreErrors.ErrAmountIsZero
	}
	if money.Amount.IsNegative() {
		return coreErrors.ErrAmountIsNegative
	}
	return nil
}
