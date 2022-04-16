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
