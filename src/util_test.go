package main

import (
	"testing"
)

func TestColorBySquare(t *testing.T) {
	tests := []struct {
		square   Square
		expected Color
	}{
		{BLACK_PAWN, BLACK},
		{WHITE_KING, WHITE},
		{EMPTY, COLOR_NONE},
	}
	for _, tt := range tests {
		actual := colorBySquare(tt.square)
		if actual != tt.expected {
			t.Errorf("color: %v != %v", actual, tt.expected)
		}
	}
}

func TestSquareIndexByCoord(t *testing.T) {
	tests := []struct {
		coord    SquareCoord
		expected SquareIndex
	}{
		{"a8", 21},
		{"h1", 98},
	}
	for _, tt := range tests {
		actual, err := squareIndexByCoord(tt.coord)
		if err != nil {
			t.Error(err)
		}
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
}

func TestSquareIndexByFileRank(t *testing.T) {
	tests := []struct {
		file     File
		rank     Rank
		expected SquareIndex
	}{
		{FILE_A, RANK_EIGHT, 21},
		{FILE_H, RANK_ONE, 98},
	}
	for _, tt := range tests {
		actual := squareIndexByFileRank(tt.file, tt.rank)
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
}

func TestRankBySquareIndex(t *testing.T) {
	tests := []struct {
		squareIndex SquareIndex
		expected    Rank
	}{
		{21, RANK_EIGHT},
		{32, RANK_SEVEN},
		{43, RANK_SIX},
		{54, RANK_FIVE},
		{65, RANK_FOUR},
		{76, RANK_THREE},
		{87, RANK_TWO},
		{98, RANK_ONE},
		{0, RANK_NONE},
	}
	for _, tt := range tests {
		actual := rankBySquareIndex(tt.squareIndex)
		if actual != tt.expected {
			t.Errorf("file: %v != %v", actual, tt.expected)
		}
	}
}

func TestFileBySquareIndex(t *testing.T) {
	tests := []struct {
		squareIndex SquareIndex
		expected    File
	}{
		{21, FILE_A},
		{32, FILE_B},
		{43, FILE_C},
		{54, FILE_D},
		{65, FILE_E},
		{76, FILE_F},
		{87, FILE_G},
		{98, FILE_H},
		{0, FILE_NONE},
	}
	for _, tt := range tests {
		actual := fileBySquareIndex(tt.squareIndex)
		if actual != tt.expected {
			t.Errorf("file: %v != %v", actual, tt.expected)
		}
	}
}
