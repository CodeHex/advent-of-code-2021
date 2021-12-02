package convert

import "strconv"

type Convertable interface {
	string | int
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

func FuncFor[T Convertable]() func(string) T {
	val := *new(T)
	switch (interface{})(val).(type) {
	case string:
		return stringConvert[T]
	case int:
		return intConvert[T]
	default:
		panic("unsupported converter")
	}
}
