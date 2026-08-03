package validate

import (
	"ITK_Code/m/v2/internal/core/errors"
	"ITK_Code/m/v2/internal/core/wallet"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func Deposit(req *pb.DepositRequest) error {
	if req.GetUserId() == "" {
		return errors.ErrUserIDEmpty
	}
	if req.GetAsset() == "" {
		return errors.ErrAssetEmpty
	}
	if req.GetAmount() == nil {
		return errors.ErrAmountEmpty
	}
	if req.GetAmount().Currency == "" {
		return errors.ErrAmountEmpty
	}
	if req.GetAsset() != req.GetAmount().Currency {
		return errors.ErrInvalidAsset
	}

	return nil
}

func Money(money wallet.Money) error {
	if money.Amount.IsZero() {
		return errors.ErrAmountIsZero
	}
	if money.Amount.IsNegative() {
		return errors.ErrAmountIsNegative
	}
	return nil
}
