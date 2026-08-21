package validate

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func CreateSpot(log *zap.Logger, reqSpot dto.CreateSpot) error {
	if strings.Compare(reqSpot.BaseAsset, reqSpot.QuoteAsset) == 0 {
		log.Error("Base Asset cannot be equal to Quote Asset")
		return errors.ErrCompareBaseQuoteAsset
	}

	minOrderSize, err := decimal.NewFromString(reqSpot.MinOrderSize)
	if err != nil {
		log.Error("Invalid MinOrderSize", zap.Error(err))
		return errors.ErrInvalidMinOrderSize
	}

	maxOrderSize, err := decimal.NewFromString(reqSpot.MaxOrderSize)
	if err != nil {
		log.Error("Invalid MaxOrderSize", zap.Error(err))
		return errors.ErrInvalidMaxOrderSize
	}

	if minOrderSize.Exponent()*-1 > reqSpot.QuantityPrecision {
		log.Error("invalid min order size relative to SizePrecision")
		return errors.ErrInvalidOrderSizePrecision
	}

	if maxOrderSize.Exponent()*-1 > reqSpot.QuantityPrecision {
		log.Error("invalid max order size relative to SizePrecision")
		return errors.ErrInvalidOrderSizePrecision
	}

	//	if reqSpot.MinOrderSize > reqSpot.MaxOrderSize {
	//		log.Error("MinOrderSize cannot be greater than MaxOrderSize")
	//		return errors.ErrInvalidMinOrderGreaterMaxOrder
	//	}
	return nil
}
