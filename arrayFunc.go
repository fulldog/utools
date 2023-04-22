package utools

import (
	"math/rand"
	"time"
)

func ArrayUnshift[T any](s *[]T, elements ...T) {
	*s = append(elements, *s...)
}

func ArrayPop[T any](arr *[]T) T {
	x := (*arr)[len(*arr)-1]
	*arr = (*arr)[0 : len(*arr)-1]
	return x
}

func ArrayShift[T any](s *[]T) T {
	x := (*s)[0]
	*s = (*s)[1:]
	return x
}

func ArrayMerge[T any](ss ...[]T) []T {
	s := make([]T, 0, len(ss[0])*len(ss))
	for i := 0; i < len(ss); i++ {
		s = append(s, ss[i]...)
	}
	return s
}

func ShuffleString[T any](slice []T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func ArrayKeys[T comparable](m map[T]any) []T {
	arr := make([]T, 0, len(m))
	for i, _ := range m {
		arr = append(arr, i)
	}
	return arr
}

func ArraySlice[T any](slic []T, str int, end int) []T {
	le := len(slic)
	if end > le {
		end = le
	}
	if str < 0 {
		str = 0
	}
	return slic[str:end]
}

func InArray[T comparable](t []T, check T) bool {
	for i, _ := range t {
		if t[i] == check {
			return true
		}
	}
	return false
}
