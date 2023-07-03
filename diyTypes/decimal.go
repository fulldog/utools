package diyTypes

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

type DiyDecimal decimal.Decimal

// MarshalJSON implements the json.Marshaller interface.
func (d DiyDecimal) MarshalJSON() ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(decimal.Decimal(d).InexactFloat64())
}

func (d DiyDecimal) ToDecimal() decimal.Decimal {
	return decimal.Decimal(d)
}
