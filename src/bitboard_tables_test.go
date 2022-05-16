package main

import "testing"

func TestBitboardPawnAttacks(t *testing.T) {
	tests := []struct {
		coord   Coord
		side    Color
		attacks []Coord
	}{
		{
			E4,
			WHITE,
			[]Coord{
				D5,
				F5,
			},
		},
		{
			E4,
			BLACK,
			[]Coord{
				D3,
				F3,
			},
		},
		{
			A4,
			WHITE,
			[]Coord{
				B5,
			},
		},
		{
			A4,
			BLACK,
			[]Coord{
				B3,
			},
		},
		{
			H4,
			WHITE,
			[]Coord{
				G5,
			},
		},
		{
			H4,
			BLACK,
			[]Coord{
				G3,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		actual := PAWN_ATTACKS[tt.side][tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect pawn attacks at %v for side %v", COORD_MAP[tt.coord], tt.side)
		}
	}
}