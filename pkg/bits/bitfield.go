package bits

import (
	"bytes"
	"strconv"
)

type BitField struct {
	Value  uint64
	Length int
	str    string
}

func NewBitField(bin string) BitField {
	str := bin
	length := len(bin)
	val, err := strconv.ParseUint(bin, 2, 64)
	if err != nil {
		panic(err)
	}
	return BitField{Value: uint64(val), Length: length, str: str}
}

func NewBitFieldForVal(val uint64, length int) BitField {
	b := BitField{Value: val, Length: length}
	str := ""
	for i := 0; i < length; i++ {
		if b.Get(i) {
			str += "1"
		} else {
			str += "0"
		}
	}
	b.str = str
	return b
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
	result := NewBitField(inverted.String())
	return result
}
