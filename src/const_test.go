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
		{IntPow(2, 4) - 1, 0, MOVE_KIND_MASK},
		{IntPow(2, 4) - 1, MOVE_DST_SQUARE_SHIFT, MOVE_DST_SQUARE_MASK},
		{IntPow(2, 4) - 1, MOVE_ORIGIN_SQUARE_SHIFT, MOVE_ORIGIN_SQUARE_MASK},
		{IntPow(2, 6) - 1, MOVE_DST_COORD_SHIFT, MOVE_DST_COORD_MASK},
		{IntPow(2, 6) - 1, MOVE_ORIGIN_COORD_SHIFT, MOVE_ORIGIN_COORD_MASK},
	}
	for _, tt := range tests {
		actual := Move(tt.bits << tt.shift)
		if actual != tt.expected {
			t.Errorf("incorrect move mask: %v != %v", actual, tt.expected)
		}
	}
}

func TestUndoMasks(t *testing.T) {
	tests := []struct {
		bits     int
		shift    int
		expected Undo
	}{
		{IntPow(2, 32) - 1, 0, UNDO_MOVE_MASK},
		{IntPow(2, 4) - 1, UNDO_CLEAR_SQUARE_SHIFT, UNDO_CLEAR_SQUARE_MASK},
		{IntPow(2, 6) - 1, UNDO_HALF_MOVE_SHIFT, UNDO_HALF_MOVE_MASK},
		{IntPow(2, 4) - 1, UNDO_CASTLE_RIGHTS_SHIFT, UNDO_CASTLE_RIGHTS_MASK},
		{IntPow(2, 6) - 1, UNDO_EP_COORD_SHIFT, UNDO_EP_COORD_MASK},
	}
	for _, tt := range tests {
		actual := Undo(tt.bits << tt.shift)
		if actual != tt.expected {
			t.Errorf("incorrect move mask: %v != %v", actual, tt.expected)
		}
	}
}
