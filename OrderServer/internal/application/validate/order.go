package validate

import "github.com/shopspring/decimal"

// NUMERIC(30,18) has room for twelve digits before the decimal point.
var maxStoredAmount = decimal.New(1, 12)

func Price(price decimal.Decimal, pricePrecision int32) bool {
	if price.GreaterThanOrEqual(maxStoredAmount) {
		return true
	}
	if decimalPlaces(price) > pricePrecision {
		return true
	}
	return false
}

func Quantity(quantity decimal.Decimal, quantityPrecision int32, minOrderSize decimal.Decimal, maxOrderSize decimal.Decimal) bool {
	if quantity.GreaterThanOrEqual(maxStoredAmount) {
		return true
	}
	if quantity.LessThan(minOrderSize) {
		return true
	}
	if quantity.GreaterThan(maxOrderSize) {
		return true
	}
	if decimalPlaces(quantity) > quantityPrecision {
		return true
	}
	return false
}

func decimalPlaces(value decimal.Decimal) int32 {
	exponent := value.Exponent()
	if exponent >= 0 {
		return 0
	}
	return -exponent
}
