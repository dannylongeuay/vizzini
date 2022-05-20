package main

import (
	"testing"
)

func TestColor(t *testing.T) {
	white := BLACK ^ 1
	if white != WHITE {
		t.Errorf("white: %v != %v", white, WHITE)
	}

	black := WHITE ^ 1
	if black != BLACK {
		t.Errorf("black: %v != %v", black, BLACK)
	}
}

func TestMoveMasks(t *testing.T) {
	tests := []struct {
		bits     int
		shift    int
		expected Move
	}{
		{MOVE_KINDS, 0, MOVE_KIND_MASK},
		{MOVE_KINDS, MOVE_DST_SQUARE_SHIFT, MOVE_DST_SQUARE_MASK},
		{MOVE_KINDS, MOVE_ORIGIN_SQUARE_SHIFT, MOVE_ORIGIN_SQUARE_MASK},
		{BOARD_SQUARES - 1, MOVE_DST_COORD_SHIFT, MOVE_DST_COORD_MASK},
		{BOARD_SQUARES - 1, MOVE_ORIGIN_COORD_SHIFT, MOVE_ORIGIN_COORD_MASK},
	}
	for _, tt := range tests {
		actual := Move(tt.bits << tt.shift)
		if actual != tt.expected {
			t.Errorf("incorrect move mask: %v != %v", actual, tt.expected)
		}
	}
}
