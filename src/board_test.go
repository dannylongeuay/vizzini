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
		checks       []testSquareChecks
	}{
		{StartingFEN, White, 15, []testSquareChecks{
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
		},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		for _, check := range tt.checks {
			square := squareByFileRank(check.file, check.rank)

			if b.squares[square] != check.square {
				t.Errorf("square: %v != %v", b.squares[square], check.square)
			}
		}

		if b.sideToMove != tt.sideToMove {
			t.Errorf("side: %v != %v", b.sideToMove, tt.sideToMove)
		}

		if b.castleRights != tt.castleRights {
			t.Errorf("castling rights: %v != %v", b.castleRights, tt.castleRights)
		}
	}
}

func TestSquareByFileRank(t *testing.T) {
	tests := []struct {
		file     File
		rank     Rank
		expected int
	}{
		{FILE_A, RANK_8, 21},
		{FILE_H, RANK_1, 98},
	}
	for _, tt := range tests {
		actual := squareByFileRank(tt.file, tt.rank)
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
}

func TestSquareByIndex(t *testing.T) {
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
