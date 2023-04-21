package utools

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"github.com/foursking/bstring"
	"github.com/fulldog/utools/ecb"
	"github.com/thinkeridea/go-extend/exstrings"
)

func PriceEncode(price string, token string, key string, bidid string) string {
	price = fmt.Sprintf("%-8s", price)
	iv := bidid[len(bidid)-16:]
	pad := HmacSha1(iv, token)[:8]
	enc_price := StrByXOR(pad, price)
	sig := HmacSha1(price+iv, key)[:4]
	b64 := bstring.Base64EncodeString(iv + enc_price + sig)
	return exstrings.Replace(exstrings.Replace(b64, "+", "-", -1), "/", "_", -1)
}

func PriceDecode(b64 string, token string) string {
	b64 = exstrings.Replace(exstrings.Replace(b64, "-", "+", -1), "_", "/", -1)
	str, err := bstring.Base64DecodeString(b64)
	if err != nil || len(str) < 28 {
		return ""
	}
	return StrByXOR(str[16:len(str)-4], HmacSha1(str[:16], token)[:8])
}

func AESEncrypt(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	ecb := ecb.NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	des := make([]byte, len(content))
	ecb.CryptBlocks(des, content)
	return des
}

// AesDecrypt ECB PKCS5 解密
func AesDecrypt(crypted, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	blockMode := ecb.NewECBDecrypter(block)
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
