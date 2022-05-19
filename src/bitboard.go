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
	return *bb & COORD_MASK_BITBOARDS[c]
}

func (bb *Bitboard) SetBit(c Coord) {
	*bb |= COORD_MASK_BITBOARDS[c]
}

func (bb *Bitboard) ClearBit(c Coord) {
	*bb &= COORD_CLEAR_BITBOARDS[c]
}

func (bb Bitboard) Count() int {
	var count int
	for bb > 0 {
		count++
		bb &= bb - 1
	}
	return count
}

func (bb Bitboard) LSBIndex() (int, error) {
	if bb == 0 {
		return 0, fmt.Errorf("cannot retrieve LSB Index from bitboard that is zero")
	}
	tmpBB := (bb & -bb) - 1
	count := tmpBB.Count()

	return count, nil
}
