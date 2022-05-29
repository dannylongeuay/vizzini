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
		search, err := NewSearch(tt.fen, UCI_DEFAULT_DEPTH, 0)
		if err != nil {
			t.Error(err)
		}
		actual := search.Negamax(1, math.MinInt+1, math.MaxInt)
		if actual != tt.expected {
			t.Errorf("\n%v != %v\n\n%v", actual, tt.expected, search.ToString())
		}
	}
}
