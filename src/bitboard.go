package main

import (
	"fmt"
)

func (bb *Bitboard) ToString() string {
	var s string
	s += "\n"
	for i, c := range PRINT_MAP {
		if i != 0 && i%8 == 0 {
			s += "\n"
		}
		var mask Bitboard
		mask.SetBit(c)
		intersection := *bb & mask
		shifted := intersection >> c
		s += fmt.Sprintf("%v ", shifted)
	}
	s += fmt.Sprintf("\n\n%v", *bb)
	return s
}

func (bb *Bitboard) HasBit(c Coord) bool {
	return bb.GetBit(c) > 0
}

func (bb *Bitboard) GetBit(c Coord) Bitboard {
	return *bb & Bitboard(1<<c)
}

func (bb *Bitboard) SetBit(c Coord) {
	*bb |= Bitboard(1 << c)
}

func (bb *Bitboard) ClearBit(c Coord) {
	if bb.HasBit(c) {
		*bb ^= Bitboard(1 << c)
	}
}
