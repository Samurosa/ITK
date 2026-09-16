package validate

import (
	"ITK_Code/m/v2/internal/core/coreErrors"
	"ITK_Code/m/v2/internal/core/spot/models"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func CreateSpot(log *zap.Logger, reqSpot models.CreateSpot) error {
	if strings.Compare(reqSpot.BaseAsset, reqSpot.QuoteAsset) == 0 {
		log.Error("Base Asset cannot be equal to Quote Asset")
		return coreErrors.ErrCompareBaseQuoteAsset
	}

	minOrderSize, err := decimal.NewFromString(reqSpot.MinOrderSize)
	if err != nil {
		log.Error("Invalid MinOrderSize", zap.Error(err))
		return coreErrors.ErrInvalidMinOrderSize
	}

	maxOrderSize, err := decimal.NewFromString(reqSpot.MaxOrderSize)
	if err != nil {
		log.Error("Invalid MaxOrderSize", zap.Error(err))
		return coreErrors.ErrInvalidMaxOrderSize
	}

	if minOrderSize.Exponent()*-1 > reqSpot.QuantityPrecision {
		log.Error("invalid min order size relative to SizePrecision")
		return coreErrors.ErrInvalidOrderSizePrecision
	}

	if maxOrderSize.Exponent()*-1 > reqSpot.QuantityPrecision {
		log.Error("invalid max order size relative to SizePrecision")
		return coreErrors.ErrInvalidOrderSizePrecision
	}

	if minOrderSize.GreaterThan(maxOrderSize) {
		log.Error("MinOrderSize cannot be greater than MaxOrderSize")
		return coreErrors.ErrInvalidMinOrderGreaterMaxOrder
	}
	return nil
}
