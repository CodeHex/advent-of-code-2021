package bits

import (
	"strconv"
)

type BitField struct {
	Value  uint
	Length int
	str    string
}

func NewBitField(bin string) (BitField, error) {
	str := bin
	length := len(bin)
	val, err := strconv.ParseUint(bin, 2, 32)
	if err != nil {
		return BitField{}, err
	}
	return BitField{Value: uint(val), Length: length, str: str}, nil
}

func (b BitField) String() string {
	return b.str
}

func (b BitField) Get(pos int) bool {
	index := b.Length - pos - 1
	return b.Value&(1<<index) != 0
}

func (b BitField) Invert() BitField {
	inverted := ""
	for _, c := range b.str {
		if c == '0' {
			inverted = inverted + "1"
		} else {
			inverted = inverted + "0"
		}
	}
	result, err := NewBitField(inverted)
	if err != nil {
		panic(err)
	}
	return result
}
