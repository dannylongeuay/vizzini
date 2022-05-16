package main

import "testing"

func IsBitboardEqual(t *testing.T, a Bitboard, b Bitboard) bool {
	if a != b {
		t.Errorf("\nactual:%v\n\n!=\n\nexpected:%v", a.ToString(), b.ToString())
		return false
	}
	return true
}

func TestBitboardToString(t *testing.T) {
	expected := `
1 0 0 0 0 0 0 1 
0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 
1 0 0 0 0 0 0 1 

9295429630892703873`
	InitBitboards()
	var bitboard Bitboard
	bitboard.SetBit(A1)
	bitboard.SetBit(A8)
	bitboard.SetBit(H1)
	bitboard.SetBit(H8)
	if bitboard.ToString() != expected {
		t.Errorf("bitboard to string: %v\n\n!=\n%v", bitboard.ToString(), expected)
	}
}

func TestBitboardHasBit(t *testing.T) {
	tests := []struct {
		coord    Coord
		bb       Bitboard
		expected bool
	}{
		{
			A1,
			Bitboard(1),
			true,
		},
		{
			A1,
			Bitboard(0),
			false,
		},
		{
			B1,
			Bitboard(2),
			true,
		},
		{
			A2,
			Bitboard(256),
			true,
		},
		{
			A3,
			Bitboard(256),
			false,
		},
		{
			B2,
			Bitboard(512),
			true,
		},
	}
	InitBitboards()
	for _, tt := range tests {
		actual := tt.bb.HasBit(tt.coord)
		if actual != tt.expected {
			t.Errorf("coord %v is not present in bitboard: %v", COORD_MAP[tt.coord], tt.bb.ToString())
		}
	}
}

func TestBitboardSetBit(t *testing.T) {
	tests := []struct {
		coords []Coord
	}{
		{[]Coord{A1, B1}},
		{[]Coord{A1, A2, A3}},
		{[]Coord{H6, H7, H8}},
	}
	InitBitboards()
	for _, tt := range tests {
		var bb Bitboard
		var expected Bitboard
		for _, c := range tt.coords {
			bb.SetBit(c)
			tmpBB := Bitboard(IntPow(2, int(c)))
			expected |= tmpBB
		}
		if !IsBitboardEqual(t, bb, expected) {
			t.Errorf("incorrect set bits: %v", tt.coords)
		}
	}
}

func TestBitboardClearBit(t *testing.T) {
	tests := []struct {
		coord    Coord
		bb       Bitboard
		expected Bitboard
	}{
		{
			A1,
			Bitboard(1),
			Bitboard(0),
		},
		{
			B1,
			Bitboard(1),
			Bitboard(1),
		},
		{
			B1,
			Bitboard(3),
			Bitboard(1),
		},
	}
	for _, tt := range tests {
		tt.bb.ClearBit(tt.coord)
		if !IsBitboardEqual(t, tt.bb, tt.expected) {
			t.Errorf("incorrect clear bit at %v", COORD_MAP[tt.coord])
		}
	}
}
