package validate

import (
	"ITK_Code/m/v2/internal/core/errors"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func Deposit(req *pb.DepositRequest) error {
	if req.GetUserId() == "" {
		return errors.ErrUserIDEmpty
	}
	if req.GetAsset() == "" {
		return errors.ErrAssetEmpty
	}
	if req.GetAmount().Currency == "" {
		return errors.ErrAmountEmpty
	}

	return nil
}
