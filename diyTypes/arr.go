package diyTypes

import jsoniter "github.com/json-iterator/go"

// DistinctList 非并发安全
type DistinctList[T comparable] struct {
	arr         []T
	arrDist     []T
	distinctMap map[T]struct{}
}

// DiyListAny 任意类型的list
type DiyListAny[T any] []T

// SumFloat64 统计浮点数
func (lst DiyListAny[T]) SumFloat64(fn func(x T) float64) (f float64) {
	for i := 0; i < len(lst); i++ {
		f += fn(lst[i])
	}
	return f
}

// SumInt 统计整数
func (lst DiyListAny[T]) SumInt(fn func(x T) int) (f int) {
	for i := 0; i < len(lst); i++ {
		f += fn(lst[i])
	}
	return f
}

// FindOne 查询一个
func (lst DiyListAny[T]) FindOne(fn func(x T) bool) (xy *T) {
	for i := 0; i < len(lst); i++ {
		if fn(lst[i]) {
			xy = &lst[i]
			return
		}
	}
	return
}

// FindAny 查找任意值
func (lst DiyListAny[T]) FindAny(fn func(x T) (any, bool)) (any, bool) {
	for i := 0; i < len(lst); i++ {
		if x, y := fn(lst[i]); y {
			return x, true
		}
	}
	return nil, false
}

// HasOne 是否含有元素
func (lst DiyListAny[T]) HasOne(fn func(x T) bool) bool {
	for i := 0; i < len(lst); i++ {
		if fn(lst[i]) {
			return true
		}
	}
	return false
}

// FindList 返回查询的集合
func (lst DiyListAny[T]) FindList(fn func(x T) bool) []T {
	var rt = make([]T, 0, len(lst))
	for i := 0; i < len(lst); i++ {
		if fn(lst[i]) {
			rt = append(rt, lst[i])
		}
	}
	return rt
}

// DistinctListString 请使用可比较类型
func (lst DiyListAny[T]) DistinctListString(fn func(x T) string) []T {
	var rt = make([]T, 0, len(lst))
	var filter = make(map[string]struct{}, len(lst))
	for i := 0; i < len(lst); i++ {
		s := fn(lst[i])
		if _, ok := filter[s]; !ok {
			filter[s] = struct{}{}
			rt = append(rt, lst[i])
		}
	}
	return rt
}

func (lst DiyListAny[T]) DistinctListInt(fn func(x T) int) []T {
	var rt = make([]T, 0, len(lst))
	var filter = make(map[int]struct{}, len(lst))
	for i := 0; i < len(lst); i++ {
		s := fn(lst[i])
		if _, ok := filter[s]; !ok {
			filter[s] = struct{}{}
			rt = append(rt, lst[i])
		}
	}
	return rt
}

func (lst DiyListAny[T]) ToMapString(fn func(x T) string) map[string]T {
	var filter = make(map[string]T, len(lst))
	for i := 0; i < len(lst); i++ {
		filter[fn(lst[i])] = lst[i]
	}
	return filter
}
func (lst DiyListAny[T]) ToMapInt(fn func(x T) int) map[int]T {
	var filter = make(map[int]T, len(lst))
	for i := 0; i < len(lst); i++ {
		filter[fn(lst[i])] = lst[i]
	}
	return filter
}

func (lst DiyListAny[T]) Cap() int {
	return cap(lst)
}
func (lst DiyListAny[T]) Len() int {
	return len(lst)
}
func (lst DiyListAny[T]) Slice(st, et uint) []T {
	if int(et) > lst.Len() {
		et = uint(lst.Len())
	}
	if st > et {
		st = et
	}
	return lst[st:et]
}

func (l *DistinctList[T]) distinct(x T) {
	if l.distinctMap == nil {
		l.distinctMap = make(map[T]struct{}, 10)
	}
	l.distinctMap[x] = struct{}{}
}
func (l *DistinctList[T]) Adds(x ...T) *DistinctList[T] {
	for _, t := range x {
		l.Add(t)
	}
	return l
}
func (l *DistinctList[T]) Add(x T) *DistinctList[T] {
	l.distinct(x)
	l.arr = append(l.arr, x)
	l.arrDist = nil
	return l
}

func (l *DistinctList[T]) Len() int {
	return len(l.arr)
}

func (l *DistinctList[T]) Distinct() []T {
	if l.arrDist != nil {
		return l.arrDist
	}
	if l.Len() > 0 {
		for k, _ := range l.distinctMap {
			l.arrDist = append(l.arrDist, k)
		}
		return l.arrDist
	}
	return nil
}

func (l *DistinctList[T]) ReValue(x []T) *DistinctList[T] {
	l.arr = nil
	l.distinctMap = nil
	l.Adds(x...)
	return l
}

// Except 排除arr元素
func (l *DistinctList[T]) Except(arr ...T) *DistinctList[T] {
	if l.distinctMap == nil {
		l.distinctMap = make(map[T]struct{}, 10)
	}
	for _, t := range arr {
		if _, ok := l.distinctMap[t]; ok {
			delete(l.distinctMap, t)
		}
	}
	l.arr = nil
	for k, _ := range l.distinctMap {
		l.Add(k)
	}
	return l
}

func (l *DistinctList[T]) ToList() []T {
	return l.arr
}

func (l *DistinctList[T]) UnmarshalJSON(data []byte) error {
	var arr []T
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	l.Adds(arr...)
	return nil
}

func (l *DistinctList[T]) ToJsonString(unique bool) (s string) {
	if unique {
		s, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(l.Distinct())
	} else {
		s, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(l.arr)
	}
	return
}
