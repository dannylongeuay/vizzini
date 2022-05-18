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
			t.Errorf("incorrect knight attacks at %v", COORD_MAP[tt.coord])
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
			t.Errorf("incorrect king attacks at %v", COORD_MAP[tt.coord])
		}
	}
}

func TestBitboardBishopAttacks(t *testing.T) {
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
		actual := BISHOP_ATTACKS[tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect bishop attacks at %v", COORD_MAP[tt.coord])
		}
	}
}
func TestBitboardRookAttacks(t *testing.T) {
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
		actual := ROOK_ATTACKS[tt.coord]
		var expected Bitboard
		for _, attackCoord := range tt.attacks {
			expected.SetBit(attackCoord)
		}
		if !IsBitboardEqual(t, actual, expected) {
			t.Errorf("incorrect rook attacks at %v", COORD_MAP[tt.coord])
		}
	}
}