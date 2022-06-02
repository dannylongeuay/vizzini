package main

import (
	"testing"
)

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
			t.Errorf("incorrect pawn attacks at %v for side %v", COORD_STRINGS[tt.coord], tt.side)
		}
	}
}

func TestBitboardKnightAttacks(t *testing.T) {
	tests := []struct {
		coord   Coord
		attacks []Coord
	}{
		{
			E4,
			[]Coord{
				C5,
				C3,
				D6,
				D2,
				F6,
				F2,
				G5,
				G3,
			},
		},
		{
			A1,
			[]Coord{
				C2,
				B3,
			},
		},
		{
			G7,
			[]Coord{
				E8,
				E6,
				F5,
				H5,
			},
		},
		{
			H3,
			[]Coord{
				G5,
				G1,
				F4,
				F2,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		actual := KNIGHT_ATTACKS[tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect knight attacks at %v", COORD_STRINGS[tt.coord])
		}
	}
}

func TestBitboardKingAttacks(t *testing.T) {
	tests := []struct {
		coord   Coord
		attacks []Coord
	}{
		{
			A4,
			[]Coord{
				A5,
				A3,
				B5,
				B4,
				B3,
			},
		},
		{
			B5,
			[]Coord{
				A6,
				A5,
				A4,
				B6,
				B4,
				C6,
				C5,
				C4,
			},
		},
		{
			H8,
			[]Coord{
				H7,
				G8,
				G7,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		actual := KING_ATTACKS[tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect king attacks at %v", COORD_STRINGS[tt.coord])
		}
	}
}

func TestBitboardBishopMasks(t *testing.T) {
	tests := []struct {
		coord   Coord
		attacks []Coord
	}{
		{
			E4,
			[]Coord{
				F5,
				G6,
				D5,
				C6,
				B7,
				D3,
				C2,
				F3,
				G2,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		actual := BISHOP_MASKS[tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect bishop masks at %v", COORD_STRINGS[tt.coord])
		}
	}
}

func TestBitboardBishopAttackGeneration(t *testing.T) {
	tests := []struct {
		coord    Coord
		blockers []Coord
		attacks  []Coord
	}{
		{
			E4,
			[]Coord{
				G6,
				B1,
				C6,
			},
			[]Coord{
				F5,
				G6,
				D5,
				C6,
				D3,
				C2,
				B1,
				F3,
				G2,
				H1,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		var blockers Bitboard
		for _, c := range tt.blockers {
			blockers.SetBit(c)
		}
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		actual := GenerateBishopAttacksBitboard(int(tt.coord), blockers)
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect bishop attacks at %v", COORD_STRINGS[tt.coord])
		}
	}
}

func TestBitboardRookMasks(t *testing.T) {
	tests := []struct {
		coord   Coord
		attacks []Coord
	}{
		{
			E4,
			[]Coord{
				E5,
				E6,
				E7,
				E3,
				E2,
				B4,
				C4,
				D4,
				F4,
				G4,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		actual := ROOK_MASKS[tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect rook masks at %v", COORD_STRINGS[tt.coord])
		}
	}
}

func TestBitboardRookAttackGeneration(t *testing.T) {
	tests := []struct {
		coord    Coord
		blockers []Coord
		attacks  []Coord
	}{
		{
			E4,
			[]Coord{
				E6,
				A4,
				F4,
			},
			[]Coord{
				E5,
				E6,
				E3,
				E2,
				E1,
				A4,
				B4,
				C4,
				D4,
				F4,
			},
		},
	}
	for _, tt := range tests {
		var blockers Bitboard
		for _, c := range tt.blockers {
			blockers.SetBit(c)
		}
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		actual := GenerateRookAttacksBitboard(int(tt.coord), blockers)
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect rook attacks at %v", COORD_STRINGS[tt.coord])
		}
	}
}

func TestBitboardBishopAttacks(t *testing.T) {
	tests := []struct {
		coord    Coord
		blockers []Coord
		attacks  []Coord
	}{
		{
			E4,
			[]Coord{
				G6,
				B1,
				C6,
			},
			[]Coord{
				F5,
				G6,
				D5,
				C6,
				D3,
				C2,
				B1,
				F3,
				G2,
				H1,
			},
		},
		{
			H1,
			[]Coord{
				G2,
				A3,
			},
			[]Coord{
				G2,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		var blockers Bitboard
		for _, c := range tt.blockers {
			blockers.SetBit(c)
		}
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		actual := BishopAttacks(tt.coord, blockers)
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect bishop attacks at %v", COORD_STRINGS[tt.coord])
		}
	}

}

func TestBitboardRookAttacks(t *testing.T) {
	tests := []struct {
		coord    Coord
		blockers []Coord
		attacks  []Coord
	}{
		{
			E4,
			[]Coord{
				E6,
				A4,
				F4,
			},
			[]Coord{
				E5,
				E6,
				E3,
				E2,
				E1,
				A4,
				B4,
				C4,
				D4,
				F4,
			},
		},
		{
			A1,
			[]Coord{
				A3,
				B1,
				H8,
			},
			[]Coord{
				A2,
				A3,
				B1,
			},
		},
	}
	InitBitboards()
	for _, tt := range tests {
		var blockers Bitboard
		for _, c := range tt.blockers {
			blockers.SetBit(c)
		}
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		actual := RookAttacks(tt.coord, blockers)
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect rook attacks at %v", COORD_STRINGS[tt.coord])
		}
	}

}

// func TestBitboardFindMagicNumbers(t *testing.T) {
// 	InitBitboards()
// 	err := FindMagicNumbers()
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
