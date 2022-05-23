package main

import (
	"testing"
)

func Perft(b *Board, depth int) int {
	var nodes int

	if depth == 0 {
		return 1
	}

	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	b.GenerateMoves(&moves, b.sideToMove)
	for _, m := range moves {
		err := b.MakeMove(m)
		if err == nil {
			nodes += Perft(b, depth-1)
		}
		b.UndoMove()
	}

	return nodes
}

func TestPerft(t *testing.T) {
	tests := []struct {
		fen      string
		depth    int
		expected int
	}{
		{STARTING_FEN, 1, 20},
		{STARTING_FEN, 2, 400},
		{STARTING_FEN, 3, 8902},
		{STARTING_FEN, 4, 197281},
		{STARTING_FEN, 5, 4865609},
		{STARTING_FEN, 6, 119060324},
	}
	SeedKeys(181818)
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actual := Perft(b, tt.depth)
		if actual != tt.expected {
			t.Errorf("nodes %v != %v", actual, tt.expected)
		}
	}
}
