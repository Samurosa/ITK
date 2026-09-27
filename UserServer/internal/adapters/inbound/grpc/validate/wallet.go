package validate

import (
	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/wallet"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func Deposit(req *pb.DepositRequest) error {
	if req.GetAmount() == nil {
		return corerrors.ErrAmountEmpty
	}
	if req.GetAmount().Currency == "" {
		return corerrors.ErrAmountEmpty
	}

	return nil
}

func Money(money wallet.Money) error {
	if money.Amount.IsZero() {
		return corerrors.ErrAmountIsZero
	}
	if money.Amount.IsNegative() {
		return corerrors.ErrAmountIsNegative
	}
	return nil
}
