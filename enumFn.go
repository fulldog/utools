package utools

import "reflect"

func StructKeyIsExist(p interface{}, k string) bool {
	return reflect.ValueOf(p).FieldByName(k).IsValid()
}

func EnumIntIsExist(p interface{}, k int, desc bool) (string, bool) {
	ref := reflect.ValueOf(p)
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Int() == int64(k) {
			if desc {
				return reflect.TypeOf(p).Field(i).Tag.Get("desc"), true
			}
			return "", true
		}
	}
	return "", false
}

func EnumIntGetDesc(p interface{}, k int) string {
	ref := reflect.ValueOf(p)
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Int() == int64(k) {
			return reflect.TypeOf(p).Field(i).Tag.Get("desc")
		}
	}
	return ""
}

func EnumIntGetDescDefault(p interface{}, k int, def string) string {
	if r := EnumIntGetDesc(p, k); r != "" {
		return r
	}
	return def
}
func GetEnumKeyByValue(enum interface{}, v interface{}, k string) (s string) {
	rvf := reflect.ValueOf(enum)
	rtf := reflect.TypeOf(enum)
	for i := 0; i < rtf.NumField(); i++ {
		if reflect.DeepEqual(v, rvf.Field(i).Interface()) {
			s = rtf.Field(i).Tag.Get(k)
			if s == "" {
				s = rtf.Field(i).Name
			}
			break
		}
	}
	return s
}
