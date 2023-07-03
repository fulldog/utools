package rmb

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/foursking/bstring"
	"github.com/fulldog/utools"
	"github.com/shopspring/decimal"
	"github.com/thinkeridea/go-extend/exstrings"
	"regexp"
	"strconv"
	"strings"
)

//func ------------------------------------------------->>>>>>>>>>

// PriceEncode 价格加密
func PriceEncode(price string, token string, key string, bidid string) string {
	price = fmt.Sprintf("%-8s", price)
	iv := bidid[len(bidid)-16:]
	pad := utools.HmacSha1(iv, token)[:8]
	enc_price := utools.StrByXOR(pad, price)
	sig := utools.HmacSha1(price+iv, key)[:4]
	b64 := bstring.Base64EncodeString(iv + enc_price + sig)
	return exstrings.Replace(exstrings.Replace(b64, "+", "-", -1), "/", "_", -1)
}

func PriceDecode(b64 string, token string) string {
	b64 = exstrings.Replace(exstrings.Replace(b64, "-", "+", -1), "_", "/", -1)
	str, err := bstring.Base64DecodeString(b64)
	if err != nil || len(str) < 28 {
		return ""
	}
	return utools.StrByXOR(str[16:len(str)-4], utools.HmacSha1(str[:16], token)[:8])
}

// AESEncrypt xx
func AESEncrypt(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	ecbx := utools.NewECBEncrypted(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	des := make([]byte, len(content))
	ecbx.CryptBlocks(des, content)
	return des
}

// AesDecrypt ECB PKCS5 解密
func AesDecrypt(crypted, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockMode := utools.NewECBDecrypted(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// ConvertStyle 转换格式 ru:1000=1K 10000=1w
func ConvertStyle(f float64, scale int32) string {
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
func GetIncrease(currentData, lastData float64) string {
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

// ToChinaCny 人民币大写
func ToChinaCny(num float64) string {
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
