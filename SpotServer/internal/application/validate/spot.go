package validate

import (
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/spot/models"

	"github.com/shopspring/decimal"
)

var maxStoredAmount = decimal.New(1, 12)

func CreateSpot(reqSpot models.CreateSpot) error {
	if reqSpot.BaseAsset == reqSpot.QuoteAsset {
		return errorcore.ErrCompareBaseQuoteAsset
	}

	minOrderSize, err := decimal.NewFromString(reqSpot.MinOrderSize)
	if err != nil {
		return errorcore.ErrInvalidMinOrderSize
	}

	maxOrderSize, err := decimal.NewFromString(reqSpot.MaxOrderSize)
	if err != nil {
		return errorcore.ErrInvalidMaxOrderSize
	}
	if !minOrderSize.IsPositive() || minOrderSize.GreaterThanOrEqual(maxStoredAmount) {
		return errorcore.ErrInvalidMinOrderSize
	}
	if !maxOrderSize.IsPositive() || maxOrderSize.GreaterThanOrEqual(maxStoredAmount) {
		return errorcore.ErrInvalidMaxOrderSize
	}

	if decimalPlaces(minOrderSize) > reqSpot.QuantityPrecision {
		return errorcore.ErrInvalidOrderSizePrecision
	}

	if decimalPlaces(maxOrderSize) > reqSpot.QuantityPrecision {
		return errorcore.ErrInvalidOrderSizePrecision
	}

	if minOrderSize.GreaterThan(maxOrderSize) {
		return errorcore.ErrInvalidMinOrderGreaterMaxOrder
	}
	return nil
}

func decimalPlaces(value decimal.Decimal) int32 {
	exponent := value.Exponent()
	if exponent >= 0 {
		return 0
	}
	return -exponent
}
