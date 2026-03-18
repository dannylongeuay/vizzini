package main

import (
	"fmt"
	"math/bits"
)

func (bb *Bitboard) ToString() string {
	var s string
	s += "\n"
	for i, c := range MIRROR_COORDS {
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
	// var count int
	// for bb > 0 {
	// 	count++
	// 	bb &= bb - 1
	// }
	return bits.OnesCount64(uint64(bb))
}

func (bb Bitboard) LSBIndex() Coord {
	if bb == 0 {
		panic(fmt.Errorf("cannot retrieve LSB Index from bitboard that is zero"))
	}

	tmpBB := (bb & -bb) - 1
	count := Coord(tmpBB.Count())

	return count
}

func (bb *Bitboard) PopLSB() Coord {
	lsbIndex := bb.LSBIndex()
	bb.ClearBit(lsbIndex)
	return lsbIndex
}
