package main

import (
	"fmt"
	"testing"
)

func divide(b *Board, depth int) map[string]int {
	results := make(map[string]int)

	moves := b.GenerateMoves(b.sideToMove)

	for _, m := range moves {
		cb := b.CopyBoard()
		err := cb.MakeMove(m)
		if err != nil {
			continue
		}
		var nodes int
		if depth > 1 {
			nodes = perft(&cb, depth-1)
		} else {
			nodes = 1
		}
		coord := fmt.Sprint("a2", "a3")
		results[coord] = nodes
	}

	return results
}

func TestDivide(t *testing.T) {
	tests := []struct {
		fen             string
		depth           int
		expectedMoves   int
		expectedResults map[string]int
	}{
		// {STARTING_FEN, 4, 20,
		// 	map[string]int{
		// 		"a2a3": 8457,
		// 		"a2a4": 9329,
		// 		"b2b3": 9345,
		// 		"b2b4": 9332,
		// 		"c2c3": 9272,
		// 		"c2c4": 9744,
		// 		"d2d3": 11959,
		// 		"d2d4": 12435,
		// 		"e2e3": 13134,
		// 		"e2e4": 13160,
		// 		"f2f3": 8457,
		// 		"f2f4": 8929,
		// 		"g2g3": 9345,
		// 		"g2g4": 9328,
		// 		"h2h3": 8457,
		// 		"h2h4": 9329,
		// 		"b1c3": 9755,
		// 		"b1a3": 8885,
		// 		"g1h3": 8881,
		// 		"g1f3": 9748,
		// 	},
		// },
		// {STARTING_FEN, 5, 20,
		// 	map[string]int{
		// 		"a2a3": 181046,
		// 		"a2a4": 217832,
		// 		"b2b3": 215255,
		// 		"b2b4": 216145,
		// 		"c2c3": 222861,
		// 		"c2c4": 240082,
		// 		"d2d3": 328511,
		// 		"d2d4": 361790,
		// 		"e2e3": 402988,
		// 		"e2e4": 405385,
		// 		"f2f3": 178889,
		// 		"f2f4": 198473,
		// 		"g2g3": 217210,
		// 		"g2g4": 214048,
		// 		"h2h3": 181044,
		// 		"h2h4": 218829,
		// 		"b1c3": 234656,
		// 		"b1a3": 198572,
		// 		"g1h3": 198502,
		// 		"g1f3": 233491,
		// 	},
		// },
	}
	SeedKeys(181818)
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actualResults := divide(b, tt.depth)
		if len(actualResults) != tt.expectedMoves {
			t.Errorf("expected moves %v != %v", tt.expectedMoves, len(actualResults))
		}
		for coord, expectedNodes := range tt.expectedResults {
			actualNodes, ok := actualResults[coord]
			if !ok {
				t.Errorf("coord %v not found", coord)
			}
			if actualNodes != expectedNodes {
				t.Errorf("coord %v nodes: %v != %v", coord, actualNodes, expectedNodes)
			}
		}
		for coord := range actualResults {
			_, ok := tt.expectedResults[coord]
			if !ok {
				t.Errorf("coord %v found but is not a valid move", coord)
			}
		}
	}
}
