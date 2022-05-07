package main

import (
	"fmt"
	"testing"
)

func divide(b *board, depth int) map[string]int {
	results := make(map[string]int)

	moves := b.generateMoves(b.sideToMove)

	for _, m := range moves {
		cb := b.copyBoard()
		err := cb.makeMove(m)
		if err != nil {
			continue
		}
		var nodes int
		if depth > 1 {
			nodes = perft(&cb, depth-1)
		} else {
			nodes = 1
		}
		coord := fmt.Sprint(coordBySquareIndex(m.origin), coordBySquareIndex(m.target))
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
		{STARTING_FEN, 4, 20,
			map[string]int{
				"a2a3": 8457,
				"a2a4": 9329,
				"b2b3": 9345,
				"b2b4": 9332,
				"c2c3": 9272,
				"c2c4": 9744,
				"d2d3": 11959,
				"d2d4": 12435,
				"e2e3": 13134,
				"e2e4": 13160,
				"f2f3": 8457,
				"f2f4": 8929,
				"g2g3": 9345,
				"g2g4": 9328,
				"h2h3": 8457,
				"h2h4": 9329,
				"b1c3": 9755,
				"b1a3": 8885,
				"g1h3": 8881,
				"g1f3": 9748,
			},
		},
		{"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1", 3, 20,
			map[string]int{
				"a7a6": 529,
				"a7a5": 585,
				"b7b6": 585,
				"b7b5": 586,
				"c7c6": 586,
				"c7c5": 661,
				"d7d6": 752,
				"d7d5": 728,
				"e7e6": 835,
				"e7e5": 891,
				"f7f6": 530,
				"f7f5": 557,
				"g7g6": 584,
				"g7g5": 567,
				"h7h6": 533,
				"h7h5": 587,
				"b8a6": 557,
				"b8c6": 613,
				"g8f6": 613,
				"g8h6": 556,
			},
		},
		{"rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2", 2, 29,
			map[string]int{
				"e1d2": 31,
				"d4d5": 29,
				"d4e5": 29,
				"a2a3": 31,
				"a2a4": 31,
				"b2b3": 31,
				"b2b4": 30,
				"c2c3": 31,
				"c2c4": 31,
				"e2e3": 31,
				"e2e4": 30,
				"f2f3": 31,
				"f2f4": 32,
				"g2g3": 31,
				"g2g4": 31,
				"h2h3": 31,
				"h2h4": 31,
				"b1c3": 31,
				"b1d2": 31,
				"b1a3": 31,
				"c1d2": 31,
				"c1e3": 31,
				"c1f4": 32,
				"c1g5": 28,
				"c1h6": 30,
				"d1d2": 31,
				"d1d3": 31,
				"g1h3": 31,
				"g1f3": 31,
			},
		},
		{"rnbqkbnr/pppp1ppp/8/4p1B1/3P4/8/PPP1PPPP/RN1QKBNR b KQkq - 1 2", 1, 28,
			map[string]int{
				"a7a6": 1,
				"a7a5": 1,
				"b7b6": 1,
				"b7b5": 1,
				"c7c6": 1,
				"c7c5": 1,
				"d7d6": 1,
				"d7d5": 1,
				"f7f6": 1,
				"f7f5": 1,
				"g7g6": 1,
				"h7h6": 1,
				"h7h5": 1,
				"e5e4": 1,
				"e5d4": 1,
				"b8a6": 1,
				"b8c6": 1,
				"d8e7": 1,
				"d8f6": 1,
				"d8g5": 1,
				"f8e7": 1,
				"f8d6": 1,
				"f8c5": 1,
				"f8b4": 1,
				"f8a3": 1,
				"g8f6": 1,
				"g8e7": 1,
				"g8h6": 1,
			},
		},
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
	seedKeys(181818)
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
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
