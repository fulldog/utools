package diyTypes

import (
	"bytes"
	"reflect"
	"strings"
	"unsafe"
)

type DiyString string
type DiyByte []byte

func NewDiyString(s string) *DiyString {
	r := DiyString(s)
	return &r
}
func (receiver *DiyString) Replace(old, new string) *DiyString {
	*receiver = DiyString(strings.Replace(string(*receiver), old, new, 1))
	return receiver
}

func (receiver *DiyString) ReplaceAll(old, new string) *DiyString {
	*receiver = DiyString(strings.ReplaceAll(string(*receiver), old, new))
	return receiver
}

func (receiver *DiyString) ToString() string {
	return string(*receiver)
}

func (receiver *DiyString) ToByte() []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(receiver))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func NewDiyByte(b []byte) *DiyByte {
	r := DiyByte(b)
	return &r
}

func (d *DiyByte) Replace(old, new []byte) *DiyByte {
	*d = bytes.Replace(*d, old, new, 1)
	return d
}
func (d *DiyByte) ReplaceAll(old, new []byte) *DiyByte {
	*d = bytes.ReplaceAll(*d, old, new)
	return d
}
func (d *DiyByte) ToByte() []byte {
	return *d
}
func (d *DiyByte) ToString() string {
	return *(*string)(unsafe.Pointer(d))
}
