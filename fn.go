package utools

import (
	"bytes"
	"compress/zlib"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/valyala/fastrand"
	"gopkg.in/gomail.v2"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

// TernaryOperation 三元运算
func TernaryOperation[T comparable](bo bool, a, c T) T {
	if bo {
		return a
	}
	return c
}

// String2Bytes []byte 转string
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Md5(str string) string {
	h := md5.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func TimeCost(fn string) func() {
	s := time.Now()
	return func() {
		fmt.Println(fmt.Sprintf("fun %s cost %d ms", fn, time.Since(s).Milliseconds()))
	}
}

func FloatParse(p float64, size int32) float64 {
	p, _ = decimal.NewFromFloat(p).Round(size).Float64()
	return p
}

func FloatParseFloor(p float64, size int32) float64 {
	p, _ = decimal.NewFromFloat(p).RoundFloor(size).Float64()
	return p
}

func SqlNullDateParse(d sql.NullTime, layout string) string {
	if d.Valid {
		return d.Time.Format(layout)
	}
	return ""
}

func DecodeInterface(in interface{}, out interface{}) {
	err := mapstructure.Decode(in, out)
	if err != nil {
		panic("decode interface error " + err.Error())
	}
}

func FastMtRand(min, max uint32) uint32 {
	return fastrand.Uint32n(max-min+1) + min
}

func Divisor(min, max int) (maxDivisor int) {
	if min > max {
		x := max
		max = min
		min = x
	}
	//用大数对小数取余
	complement := max % min
	//余数不为零，小数作为大数,将余数作为小数，大数对小数递归求余
	if complement != 0 {
		maxDivisor = Divisor(complement, min)
	} else {
		//当余数为零，小数就是最大公约数
		maxDivisor = min
	}
	return maxDivisor
}

// UderscoreToUpperCamelCase 下划线单词转为大写驼峰单词
func UderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

// UderscoreToLowerCamelCase 下划线单词转为小写驼峰单词
func UderscoreToLowerCamelCase(s string) string {
	s = UderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// CalculatePercentage 数据概览计算同比/环比 返回 1,234.111
func CalculatePercentage(t1, t2 float64) string {
	if t2 == 0 {
		return "0.00"
	}
	return humanize.Commaf(decimal.NewFromFloat((t1/t2 - 1) * 100).Round(2).InexactFloat64())
}

// CalculateIncrease  数据概览计算同比/环比  1234.111
func CalculateIncrease(currentData, lastData float64) string {
	if lastData <= 0 {
		return ""
	}
	return decimal.NewFromFloat((currentData/lastData - 1) * 100).Round(2).String()
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandString 获取随机字符串
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func FindOneFromList[T any](t []*T, fn func(x *T) bool) *T {
	for i, _ := range t {
		if fn(t[i]) {
			return t[i]
		}
	}
	return nil
}

// FilterList 过滤数据 根据需要返回
func FilterList[T any, T2 comparable](in []T, fn func(c T) (T2, bool)) []T2 {
	rt := make([]T2, 0, len(in))
	filter := make(map[T2]struct{}, len(in))
	for i := 0; i < len(in); i++ {
		x, ok := fn(in[i])
		if ok {
			if _, y := filter[x]; !y {
				rt = append(rt, x)
				filter[x] = struct{}{}
			}
		}
	}
	return rt
}
func FindStringOrInt[T any, T2 string | int](in []T, fn func(x T) T2) []T2 {
	var rt []T2

	for i := 0; i < len(in); i++ {
		rt = append(rt, fn(in[i]))
	}
	return rt
}

func ObjectGetOrNil[T any, T2 any](obj *T, fn func(x *T) T2) T2 {
	return fn(obj)
}

func FindListFormList[T any](in []T, fn func(c T) bool) []T {
	var t = make([]T, 0, len(in))
	for i := 0; i < len(in); i++ {
		if fn(in[i]) {
			t = append(t, in[i])
		}
	}
	return t
}
func FindAny[T any](t []T, fn func(c T) bool) bool {
	for i, _ := range t {
		if fn(t[i]) {
			return true
		}
	}
	return false
}

// ArrExcept ab差集
func ArrExcept[T comparable](a, b []T) []T {
	var mp = make(map[T]struct{}, len(a))
	for _, t := range b {
		mp[t] = struct{}{}
	}
	var r []T
	for _, v := range a {
		if _, bo := mp[v]; !bo {
			r = append(r, v)
		}
	}
	return r
}
func ArrUnique[T comparable](a []T) []T {
	arr := make([]T, 0, len(a))
	mp := make(map[T]struct{}, len(a))
	for _, k := range a {
		if _, bo := mp[k]; !bo {
			mp[k] = struct{}{}
			arr = append(arr, k)
		}
	}
	return arr
}

// VerifyFormat email verify
func VerifyFormat(verify, pattern string) bool {
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(verify)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type EmailModel struct {
	From        string
	To          []string
	Cc          []string
	Subject     string
	Body        string
	File        []string
	ContentType string
}

func SendEmail(req *EmailModel) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", req.From, "天玑")
	m.SetHeader("To", req.To...)
	m.SetHeader("Subject", req.Subject)
	m.SetBody(req.ContentType, req.Body)
	if len(req.Cc) > 0 {
		m.SetHeader("Cc", req.Cc...)
	}
	for _, s := range req.File {
		m.Attach(s)
	}
	return stmp.DialAndSend(m)
}

func GetClientIP(r *http.Request) string {
	for _, h := range []string{"X-Real-Ip", "X-Forwarded-For"} {
		ips := r.Header.Get(h)
		if ips != "" {
			return strings.TrimSpace(strings.Split(ips, ",")[0])
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

func GetIdentityStr(prefix string) string {
	t := time.Now()
	timeStr := t.Format("20060102150405") // yyyyMMddHHmmss
	hashCodeV := rand.New(rand.NewSource(t.UnixNano())).Int63()
	if hashCodeV < 0 {
		hashCodeV = -hashCodeV
	}
	return prefix + timeStr + fmt.Sprintf("%011d", hashCodeV)
}

func Recover(fu func(r any)) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				fu(r)
			}
		}()
	}
}

func CreateNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func CreateNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

func CreateNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(i),
		Valid: i != 0,
	}
}

func CreateNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   f != 0,
	}
}
func StrByXOR(message string, keywords string) string {
	messageLen := len(message)
	keywordsLen := len(keywords)

	result := ""

	for i := 0; i < messageLen; i++ {
		result += string(message[i] ^ keywords[i%keywordsLen])
	}
	return result
}

func StrPadLeft(input string, padLength int, padString string) string {
	output := padString

	for padLength > len(output) {
		output += output
	}

	if len(input) >= padLength {
		return input
	}

	return output[:padLength-len(input)] + input
}

func StrSplit(str string, length int) []string {
	strs := []rune(str)
	c := len(strs)
	var arr []string
	if length < 1 || length >= c {
		return []string{str}
	}
	for i := 0; i < c; i += length {
		arr = append(arr, string(strs[i:i+length]))
	}
	return arr
}

// MicroTime php 毫秒
func MicroTime() int64 {
	return time.Now().UnixNano() / 1000000
}

func Gzuncompress(s string) string {
	var out bytes.Buffer
	in := bytes.NewBufferString(s)
	r, _ := zlib.NewReader(in)
	io.Copy(&out, r)
	return out.String()
}

func MaxInt(ii ...int) int {
	sort.Ints(ii)
	return ii[len(ii)-1]
}
func MinInt(ii ...int) int {
	sort.Ints(ii)
	return ii[0]
}

func HmacSha256(src string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
	//shaStr:=hex.EncodeToString(h.Sum(nil))
	//return base64.StdEncoding.EncodeToString([]byte(shaStr))
}

func HmacSha1(src string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
	//shaStr:=hex.EncodeToString(h.Sum(nil))
	//return base64.StdEncoding.EncodeToString([]byte(shaStr))
}
