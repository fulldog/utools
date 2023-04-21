package diyTypes

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

type Li2yuan float64

func (l *Li2yuan) UnmarshalJSON(data []byte) error {
	var t float64
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, &t)
	if err == nil {
		t = decimal.NewFromFloat(t * 0.001).Round(3).InexactFloat64()
		*l = Li2yuan(t)
	}
	return nil
}
func (l *Li2yuan) ToFloat64() float64 {
	return float64(*l)
}

type Fen2yuan float64

func (l *Fen2yuan) UnmarshalJSON(data []byte) error {
	var t float64
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, &t)
	if err == nil {
		t = decimal.NewFromFloat(t * 0.01).Round(3).InexactFloat64()
		*l = Fen2yuan(t)
	}
	return nil
}
func (l *Fen2yuan) ToFloat64() float64 {
	return float64(*l)
}
