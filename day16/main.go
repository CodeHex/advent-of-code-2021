package main

import (
	"adventofcode2021/pkg/bits"
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/slices"
	"fmt"
	"strings"
)

func main() {
	transmission := fileparser.ReadSingles[string]("day16/input.txt")[0]
	hexToBinary := map[rune]string{
		'0': "0000",
		'1': "0001",
		'2': "0010",
		'3': "0011",
		'4': "0100",
		'5': "0101",
		'6': "0110",
		'7': "0111",
		'8': "1000",
		'9': "1001",
		'A': "1010",
		'B': "1011",
		'C': "1100",
		'D': "1101",
		'E': "1110",
		'F': "1111",
	}

	// convert hex to bits string
	parts := slices.Map([]rune(transmission), func(u rune) string { return hexToBinary[u] })
	binStr := strings.Join(parts, "")

	p := ParsePacket(binStr)

	fmt.Println("[Part 1] Sum of version:", VersionCodeSum(p))
	fmt.Println("[Part 2] Transmission calculation:", Calculate(p))
}

type Packet struct {
	Version      uint64
	TypeID       uint64
	LengthTypeID uint64
	Length       uint64
	Literal      uint64
	BitLength    uint64
	SubPackets   []Packet
}

func ParsePacket(binStr string) Packet {
	index := 0
	p := Packet{}

	// Read version (3 bits)
	versionStr := binStr[index : index+3]
	versBits := bits.NewBitField(versionStr)
	p.Version = versBits.Value
	index += 3

	// Read type ID (3 bits)
	typeIDStr := binStr[index : index+3]
	typeIDBits := bits.NewBitField(typeIDStr)
	p.TypeID = typeIDBits.Value
	index += 3

	// Literal value
	if typeIDBits.Value == 4 {
		literalBin := ""
		for {
			// Read first bit prefix
			firstBit := binStr[index]
			index++
			literalBin += binStr[index : index+4]
			index += 4
			if firstBit == '0' {
				break
			}
		}
		literalBits := bits.NewBitField(literalBin)
		p.Literal = literalBits.Value
	} else {
		// Operator value
		lengthTypeIDBits := bits.NewBitField(string(binStr[index]))
		p.LengthTypeID = lengthTypeIDBits.Value
		index++

		// Length in bits
		if p.LengthTypeID == 0 {
			lengthValStr := binStr[index : index+15]
			lengthValBits := bits.NewBitField(lengthValStr)
			p.Length = lengthValBits.Value
			index += 15
			totalBitLength := uint64(0)
			for totalBitLength < p.Length {
				packet := ParsePacket(binStr[index:])
				p.SubPackets = append(p.SubPackets, packet)
				index += int(packet.BitLength)
				totalBitLength = slices.SumWeighted(p.SubPackets, func(t Packet) uint64 { return t.BitLength })
			}

		} else {
			// Length in packets
			lengthValStr := binStr[index : index+11]
			lengthValBits := bits.NewBitField(lengthValStr)
			p.Length = lengthValBits.Value
			index += 11
			for len(p.SubPackets) < int(p.Length) {
				packet := ParsePacket(binStr[index:])
				p.SubPackets = append(p.SubPackets, packet)
				index += int(packet.BitLength)
			}
		}
	}
	p.BitLength = uint64(index)
	return p
}

func VersionCodeSum(p Packet) uint64 {
	result := p.Version
	for _, sub := range p.SubPackets {
		result += VersionCodeSum(sub)
	}
	return result
}

func Calculate(p Packet) int {
	vals := []int{}
	for _, s := range p.SubPackets {
		vals = append(vals, Calculate(s))
	}

	switch p.TypeID {
	case 0: // Sum
		return slices.Sum(vals)
	case 1: // Product
		result := 1
		for _, val := range vals {
			result *= val
		}
		return result
	case 2: // Minimum
		return slices.Min(vals)
	case 3: // Maximum
		return slices.Max(vals)
	case 4: // Literal
		return int(p.Literal)
	case 5: // Greater
		if len(vals) != 2 {
			panic("greater doesn't have 2 args")
		}
		if vals[0] > vals[1] {
			return 1
		}
		return 0
	case 6: // Less
		if len(vals) != 2 {
			panic("less doesn't have 2 args")
		}
		if vals[0] < vals[1] {
			return 1
		}
		return 0

	case 7: // Equal
		if len(vals) != 2 {
			panic("less doesn't have 2 args")
		}
		if vals[0] == vals[1] {
			return 1
		}
		return 0
	}
	panic("unrecgonized type")
}
