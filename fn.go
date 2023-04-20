package utools

func TernaryOperation[T int | string | int8 | int64 | int32 | int16 | float64](bo bool, a, c T) T {
	if bo {
		return a
	}
	return c
}
