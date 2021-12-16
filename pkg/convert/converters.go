package convert

import (
	"adventofcode2021/pkg/bits"
	"strconv"
)

type Convertable interface {
	string | int | bits.BitField
}

func stringConvert[T Convertable](x string) T {
	return (interface{})(x).(T)
}

func intConvert[T Convertable](x string) T {
	r, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	return (interface{})(r).(T)
}

func bitFieldConvert[T Convertable](x string) T {
	return (interface{})(bits.NewBitField(x)).(T)
}

func FuncFor[T Convertable]() func(string) T {
	val := *new(T)
	switch (interface{})(val).(type) {
	case string:
		return stringConvert[T]
	case int:
		return intConvert[T]
	case bits.BitField:
		return bitFieldConvert[T]
	default:
		panic("unsupported converter")
	}
}

func Apply[T Convertable](in string) T {
	return FuncFor[T]()(in)
}
