package bits

import (
	"bytes"
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
	var inverted bytes.Buffer
	for _, c := range b.str {
		if c == '0' {
			inverted.WriteString("1")
		} else {
			inverted.WriteString("0")
		}
	}
	result, err := NewBitField(inverted.String())
	if err != nil {
		panic(err)
	}
	return result
}
