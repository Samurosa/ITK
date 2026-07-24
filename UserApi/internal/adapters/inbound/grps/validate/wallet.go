package validate

import pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"

func Deposit(req *pb.DepositRequest) error {
	if req.GetUserId() == "" {
		return ErrUserIDEmpty
	}
	if req.GetAsset() == "" {
		return ErrAssetEmpty
	}
	if req.GetAmount().Currency == "" {
		return ErrAmountEmpty
	}

	return nil
}
