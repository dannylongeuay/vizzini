package main

import (
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
		search, err := NewSearch(tt.fen, DEFAULT_MAX_DEPTH, DEFAULT_MAX_NODES)
		if err != nil {
			t.Error(err)
		}
		actual := search.Negamax(1, MIN_SCORE, MAX_SCORE)
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
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH, 0},
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET, 0},
			},
			false,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET, 0},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET, 0},
			},
			true,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH, 0},
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET, 0},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET, 0},
			},
			true,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET, 0},
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH, 0},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET, 0},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET, 0},
			},
			false,
		},
	}
	for _, tt := range tests {
		search, err := NewSearch(tt.fen, DEFAULT_MAX_DEPTH, DEFAULT_MAX_NODES)
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

func TestSearchMateInX(t *testing.T) {
	tests := []struct {
		fen         string
		movesToMate int
		expected    string
	}{
		{
			"1r5k/p6p/6p1/2pQPP2/2P3P1/3R4/q7/1N1K4 b - - 4 41",
			1,
			"pv b8b1",
		},
		{
			"3r4/2pq1p2/3k4/p5Q1/3P4/P5P1/7P/6K1 w - - 7 35",
			2,
			"pv g5c5 d6e6 c5e5",
		},
		{
			"7r/p7/2p5/1p3pk1/2pq4/6K1/PP4P1/3RQR2 b - - 2 26",
			3,
			"pv d4h4 g3f3 h4f4 f3e2 h8e8",
		},
	}
	for _, tt := range tests {
		if testing.Short() && tt.movesToMate > 2 {
			continue
		}
		search, err := NewSearch(tt.fen, tt.movesToMate*2, DEFAULT_MAX_NODES)
		search.quiet = true
		if err != nil {
			t.Error(err)
		}
		search.IterativeDeepening()
		actual := search.GetPvLineString()
		if actual != tt.expected {
			t.Errorf("pv line: %v != %v", actual, tt.expected)
		}

	}
}
