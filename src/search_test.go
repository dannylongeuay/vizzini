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

func TestSearchRepetition(t *testing.T) {
	tests := []struct {
		fen      string
		mus      []MoveUnpacked
		expected bool
	}{
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
			},
			false,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET},
			},
			true,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET},
			},
			true,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET},
			},
			false,
		},
	}
	for _, tt := range tests {
		search, err := NewSearch(tt.fen, UCI_DEFAULT_DEPTH, 0)
		if err != nil {
			t.Error(err)
		}
		sMoves := "moves"
		for _, mu := range tt.mus {
			move := NewMoveFromMoveUnpacked(mu)
			sMoves += " " + move.ToUCIString()
			search.Board.MakeMove(move)
		}
		actual := search.Repetition()
		if actual != tt.expected {
			t.Errorf("repetition result: %v != %v for %v", actual, tt.expected, sMoves)
		}
	}

}
