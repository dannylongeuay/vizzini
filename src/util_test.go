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
