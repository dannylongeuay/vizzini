package main

import (
	"testing"
)

type testSquareChecks struct {
	file   File
	rank   Rank
	square Square
}

func TestNewBoard(t *testing.T) {
	tests := []struct {
		fen          string
		sideToMove   Color
		castleRights Bits
		epIndex      int
		halfMove     int
		fullMove     int
		checks       []testSquareChecks
	}{
		{StartingFEN, White, 15, -1, 0, 1, []testSquareChecks{
			{FILE_A, RANK_8, BlackRook},
			{FILE_B, RANK_8, BlackKnight},
			{FILE_C, RANK_8, BlackBishop},
			{FILE_D, RANK_8, BlackQueen},
			{FILE_E, RANK_8, BlackKing},
			{FILE_F, RANK_8, BlackBishop},
			{FILE_G, RANK_8, BlackKnight},
			{FILE_H, RANK_8, BlackRook},
			{FILE_A, RANK_7, BlackPawn},
			{FILE_B, RANK_7, BlackPawn},
			{FILE_C, RANK_7, BlackPawn},
			{FILE_D, RANK_7, BlackPawn},
			{FILE_E, RANK_7, BlackPawn},
			{FILE_F, RANK_7, BlackPawn},
			{FILE_G, RANK_7, BlackPawn},
			{FILE_H, RANK_7, BlackPawn},
			{FILE_A, RANK_6, Empty},
			{FILE_H, RANK_3, Empty},
			{FILE_A, RANK_2, WhitePawn},
			{FILE_B, RANK_2, WhitePawn},
			{FILE_C, RANK_2, WhitePawn},
			{FILE_D, RANK_2, WhitePawn},
			{FILE_E, RANK_2, WhitePawn},
			{FILE_F, RANK_2, WhitePawn},
			{FILE_G, RANK_2, WhitePawn},
			{FILE_H, RANK_2, WhitePawn},
			{FILE_A, RANK_1, WhiteRook},
			{FILE_B, RANK_1, WhiteKnight},
			{FILE_C, RANK_1, WhiteBishop},
			{FILE_D, RANK_1, WhiteQueen},
			{FILE_E, RANK_1, WhiteKing},
			{FILE_F, RANK_1, WhiteBishop},
			{FILE_G, RANK_1, WhiteKnight},
			{FILE_H, RANK_1, WhiteRook},
			{FILE_None, RANK_None, Invalid},
		},
		},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			Black, 15, 75, 0, 1, []testSquareChecks{
				{FILE_E, RANK_4, WhitePawn},
			}},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			White, 15, 43, 0, 2, []testSquareChecks{
				{FILE_C, RANK_5, BlackPawn},
			}},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			Black, 15, -1, 1, 2, []testSquareChecks{
				{FILE_F, RANK_3, WhiteKnight},
			}},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
			White, 0, -1, 0, 1, []testSquareChecks{}},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w q - 0 1",
			White, 1, -1, 0, 1, []testSquareChecks{}},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1",
			White, 2, -1, 0, 1, []testSquareChecks{}},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Q - 0 1",
			White, 4, -1, 0, 1, []testSquareChecks{}},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w K - 0 1",
			White, 8, -1, 0, 1, []testSquareChecks{}},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		for _, check := range tt.checks {
			squareIndex := squareIndexByFileRank(check.file, check.rank)

			if b.squares[squareIndex] != check.square {
				t.Errorf("square: %v != %v", b.squares[squareIndex], check.square)
			}
		}

		if b.sideToMove != tt.sideToMove {
			t.Errorf("side: %v != %v", b.sideToMove, tt.sideToMove)
		}

		if b.castleRights != tt.castleRights {
			t.Errorf("castling rights: %v != %v", b.castleRights, tt.castleRights)
		}

		if b.epIndex != tt.epIndex {
			t.Errorf("en passant index: %v != %v", b.epIndex, tt.epIndex)
		}

		if b.halfMove != tt.halfMove {
			t.Errorf("half move clock: %v != %v", b.halfMove, tt.halfMove)
		}

		if b.fullMove != tt.fullMove {
			t.Errorf("full move clock: %v != %v", b.fullMove, tt.fullMove)
		}
	}
}

func TestSquareIndexByFileRank(t *testing.T) {
	tests := []struct {
		file     File
		rank     Rank
		expected int
	}{
		{FILE_A, RANK_8, 21},
		{FILE_H, RANK_1, 98},
	}
	for _, tt := range tests {
		actual := squareIndexByFileRank(tt.file, tt.rank)
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
}

func TestSquareByIndexes64(t *testing.T) {
	tests := []struct {
		index    int
		expected int
	}{
		{0, 21},
		{63, 98},
	}
	for _, tt := range tests {

		actual := SquareIndexes64[tt.index]
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
}
