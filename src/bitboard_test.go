package main

import "testing"

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
	var bitboard Bitboard
	bitboard.SetBit(A1)
	bitboard.SetBit(A8)
	bitboard.SetBit(H1)
	bitboard.SetBit(H8)
	if bitboard.ToString() != expected {
		t.Errorf("bitboard to string: %v \n\n!=\n %v", bitboard.ToString(), expected)
	}
}
