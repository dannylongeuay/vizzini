package main

import (
	"math"
	"testing"
)

func TestSearchNegamax(t *testing.T) {
	tests := []struct {
		fen      string
		expected int
	}{
		{
			STARTING_FEN,
			30,
		},
	}
	for _, tt := range tests {
		board, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		search := Search{Board: board}
		actual := search.Negamax(1, math.MinInt+1, math.MaxInt)
		if actual != tt.expected {
			t.Errorf("\n%v != %v\n\n%v", actual, tt.expected, board.ToString())
		}
	}
}
