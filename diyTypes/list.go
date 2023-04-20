package diyTypes

import "github.com/pkg/errors"

var errNotFound error = errors.New("not found")

// distinct 非并发安全
type distinctList[T ~int | ~string] struct {
	arr         []T
	distinctMap map[T]struct{}
	arrDist     []T
	//lock        sync.Mutex
}

// DiyListAny 任意类型的list
type DiyListAny[T any] []T

func (lst DiyListAny[T]) FindOne(fn func(x T) bool) (xy T, err error) {
	for i := 0; i < len(lst); i++ {
		if fn(lst[i]) {
			xy = lst[i]
			return
		}
	}
	err = errNotFound
	return
}
func (lst DiyListAny[T]) FindAny(fn func(x T) (any, bool)) (any, bool) {
	for i := 0; i < len(lst); i++ {
		if x, y := fn(lst[i]); y {
			return x, true
		}
	}
	return nil, false
}
func (lst DiyListAny[T]) HasOne(fn func(x T) bool) bool {
	for i := 0; i < len(lst); i++ {
		if fn(lst[i]) {
			return true
		}
	}
	return false
}

func (lst DiyListAny[T]) FindList(fn func(x T) bool) []T {
	var rt = make([]T, len(lst))
	for i := 0; i < len(lst); i++ {
		if fn(lst[i]) {
			rt = append(rt, lst[i])
		}
	}
	return rt
}

// DistinctListString 请使用可比较类型
func (lst DiyListAny[T]) DistinctListString(fn func(x T) string) []T {
	var rt = make([]T, len(lst))
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
	var rt = make([]T, len(lst))
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

func NewDistinctList[T ~int | ~string]() *distinctList[T] {
	return &distinctList[T]{}
}
func (l *distinctList[T]) distinct(x T) {
	//l.lock.Lock()
	//defer l.lock.Unlock()
	if l.distinctMap == nil {
		l.distinctMap = make(map[T]struct{}, 10)
	}
	l.distinctMap[x] = struct{}{}
}
func (l *distinctList[T]) Adds(x []T) *distinctList[T] {
	for _, t := range x {
		l.Add(t)
	}
	return l
}
func (l *distinctList[T]) Add(x T) *distinctList[T] {
	l.distinct(x)
	l.arr = append(l.arr, x)
	l.arrDist = nil
	return l
}

func (l *distinctList[T]) Len() int {
	return len(l.arr)
}

func (l *distinctList[T]) Distinct() []T {
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

func (l *distinctList[T]) ReValue(x []T) *distinctList[T] {
	l.arr = nil
	l.distinctMap = nil
	l.Adds(x)
	return l
}

// Except 排除arr元素
func (l *distinctList[T]) Except(arr []T) *distinctList[T] {
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

func (l *distinctList[T]) ToList() []T {
	return l.arr
}
