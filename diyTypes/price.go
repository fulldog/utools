package diyTypes

import (
	"github.com/dustin/go-humanize"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"regexp"
	"strconv"
	"strings"
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

var AmountExt = &amountExt{}

type amountExt struct{}

func (a amountExt) ConvertStyle(f float64, scale int32) string {
	if f <= 0 {
		return "0"
	}
	if f < 1000 {
		return decimal.NewFromFloat(f).Round(scale).String()
	}
	if f < 10000 {
		return decimal.NewFromFloat(f/1000).Round(scale).String() + "K"
	}
	return decimal.NewFromFloat(f/10000).Round(scale).String() + "W"
}

// GetIncrease 计算涨幅
func (a amountExt) GetIncrease(currentData, lastData float64) string {
	if lastData <= 0 {
		return ""
	}
	s := humanize.Commaf(decimal.NewFromFloat((currentData/lastData - 1) * 100).Round(2).InexactFloat64())
	sr := strings.Split(s, ".")
	if len(sr) == 1 {
		s += ".00"
	} else if len(sr) == 2 {
		if len(sr[1]) == 1 {
			s += "0"
		}
	}
	return s
}
func (a amountExt) Float2ChinaCny(num float64) string {
	money := num
	if num < 0 {
		money = -num
	}
	strnum := strconv.FormatFloat(money*100, 'f', 0, 64)
	sliceUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	s := sliceUnit[len(sliceUnit)-len(strnum):]
	upperDigitUnit := map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
	str := ""
	for k, v := range strnum[:] {
		str = str + upperDigitUnit[string(v)] + s[k]
	}
	reg, _ := regexp.Compile(`零角零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, _ = regexp.Compile(`零角`)
	str = reg.ReplaceAllString(str, "零")

	reg, _ = regexp.Compile(`零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, _ = regexp.Compile(`零[仟佰拾]`)
	str = reg.ReplaceAllString(str, "零")

	reg, _ = regexp.Compile(`零{2,}`)
	str = reg.ReplaceAllString(str, "零")

	reg, _ = regexp.Compile(`零亿`)
	str = reg.ReplaceAllString(str, "亿")

	reg, _ = regexp.Compile(`零万`)
	str = reg.ReplaceAllString(str, "万")

	reg, _ = regexp.Compile(`零*元`)
	str = reg.ReplaceAllString(str, "元")

	reg, _ = regexp.Compile(`亿零{0, 3}万`)
	str = reg.ReplaceAllString(str, "^元")

	reg, _ = regexp.Compile(`零元`)
	str = reg.ReplaceAllString(str, "零")
	if num < 0 {
		str = "负" + str
	}
	return str
}
