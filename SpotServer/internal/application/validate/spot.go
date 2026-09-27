package validate

import (
	"ITK_Code/m/v2/internal/core/coreErrors"
	"ITK_Code/m/v2/internal/core/spot/models"
	"strings"

	"github.com/shopspring/decimal"
)

func CreateSpot(reqSpot models.CreateSpot) error {
	if strings.Compare(reqSpot.BaseAsset, reqSpot.QuoteAsset) == 0 {
		return coreErrors.ErrCompareBaseQuoteAsset
	}

	minOrderSize, err := decimal.NewFromString(reqSpot.MinOrderSize)
	if err != nil {
		return coreErrors.ErrInvalidMinOrderSize
	}

	maxOrderSize, err := decimal.NewFromString(reqSpot.MaxOrderSize)
	if err != nil {
		return coreErrors.ErrInvalidMaxOrderSize
	}
	if !minOrderSize.IsPositive() {
		return coreErrors.ErrInvalidMinOrderSize
	}
	if !maxOrderSize.IsPositive() {
		return coreErrors.ErrInvalidMaxOrderSize
	}

	if decimalPlaces(minOrderSize) > reqSpot.QuantityPrecision {
		return coreErrors.ErrInvalidOrderSizePrecision
	}

	if decimalPlaces(maxOrderSize) > reqSpot.QuantityPrecision {
		return coreErrors.ErrInvalidOrderSizePrecision
	}

	if minOrderSize.GreaterThan(maxOrderSize) {
		return coreErrors.ErrInvalidMinOrderGreaterMaxOrder
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
