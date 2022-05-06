package main

import (
	"fmt"
	"testing"
)

func perft(b *board, depth int) int {
	nodes := 0

	moves := b.generateMoves(b.sideToMove)

	if depth == 1 {
		return len(moves)
	}

	for _, m := range moves {
		b.makeMove(m)
		nodes += perft(b, depth-1)
		b.undoMove()
	}

	return nodes
}

func divide(b *board, depth int) map[string]int {
	results := make(map[string]int)

	moves := b.generateMoves(b.sideToMove)

	for _, m := range moves {
		cb := b.copyBoard()
		cb.makeMove(m)
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

func TestPerft(t *testing.T) {
	tests := []struct {
		fen      string
		depth    int
		expected int
	}{
		{STARTING_FEN, 1, 20},
		{STARTING_FEN, 2, 400},
		{STARTING_FEN, 3, 8902},
		{STARTING_FEN, 4, 197281},
		// {STARTING_FEN, 5, 4865609},
		// {STARTING_FEN, 6, 119060324},
	}
	seedKeys(181818)
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actual := perft(b, tt.depth)
		if actual != tt.expected {
			t.Errorf("nodes %v != %v", actual, tt.expected)
		}
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		fen             string
		depth           int
		expectedMoves   int
		expectedResults map[string]int
	}{
		{STARTING_FEN, 3, 20,
			map[string]int{
				"b1a3": 400,
				"b1c3": 440,
				"g1f3": 440,
				"g1h3": 400,
				"a2a4": 420,
				"a2a3": 380,
				"b2b4": 421,
				"b2b3": 420,
				"c2c4": 441,
				"c2c3": 420,
				"d2d4": 560,
				"d2d3": 539,
				"e2e4": 600,
				"e2e3": 599,
				"f2f4": 401,
				"f2f3": 380,
				"g2g4": 421,
				"g2g3": 420,
				"h2h4": 420,
				"h2h3": 380,
			},
		},
		{"rnbqkbnr/pppppppp/8/8/8/2N5/PPPPPPPP/R1BQKBNR b KQkq - 1 1", 2, 20,
			map[string]int{
				"a7a6": 22,
				"a7a5": 22,
				"b7b6": 22,
				"b7b5": 22,
				"c7c6": 22,
				"c7c5": 22,
				"d7d6": 22,
				"d7d5": 22,
				"e7e6": 22,
				"e7e5": 22,
				"f7f6": 22,
				"f7f5": 22,
				"g7g6": 22,
				"g7g5": 22,
				"h7h6": 22,
				"h7h5": 22,
				"b8a6": 22,
				"b8c6": 22,
				"g8f6": 22,
				"g8h6": 22,
			},
		},
		{"rnbqkbnr/1ppppppp/p7/8/8/2N5/PPPPPPPP/R1BQKBNR w KQkq - 0 2", 1, 22,
			map[string]int{
				"a2a3": 1,
				"a2a4": 1,
				"b2b3": 1,
				"b2b4": 1,
				"d2d3": 1,
				"d2d4": 1,
				"e2e3": 1,
				"e2e4": 1,
				"f2f3": 1,
				"f2f4": 1,
				"g2g3": 1,
				"g2g4": 1,
				"h2h3": 1,
				"h2h4": 1,
				"c3b1": 1,
				"c3d5": 1,
				"c3e4": 1,
				"c3a4": 1,
				"c3b5": 1,
				"a1b1": 1,
				"g1h3": 1,
				"g1f3": 1,
			},
		},
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
		{"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1", 3, 20,
			map[string]int{
				"a7a6": 416,
				"a7a5": 460,
				"b7b6": 460,
				"b7b5": 461,
				"c7c6": 460,
				"c7c5": 484,
				"d7d6": 591,
				"d7d5": 612,
				"e7e6": 656,
				"e7e5": 657,
				"f7f6": 416,
				"f7f5": 438,
				"g7g6": 460,
				"g7g5": 461,
				"h7h6": 417,
				"h7h5": 459,
				"b8a6": 438,
				"b8c6": 482,
				"g8f6": 482,
				"g8h6": 438,
			},
		},
		{"rnbqkbnr/pppppp1p/8/6p1/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 0 2", 2, 22,
			map[string]int{
				"a2a3": 21,
				"a2a4": 21,
				"b2b3": 21,
				"b2b4": 21,
				"c2c3": 21,
				"c2c4": 21,
				"d2d3": 21,
				"d2d4": 21,
				"e2e3": 21,
				"e2e4": 21,
				"g2g3": 21,
				"g2g4": 20,
				"h2h3": 21,
				"h2h4": 22,
				"f3g5": 20,
				"f3h4": 22,
				"f3d4": 21,
				"f3g1": 21,
				"f3e5": 20,
				"b1c3": 21,
				"b1a3": 21,
				"h1g1": 21,
			},
		},
		{"rnbqkbnr/pppppp1p/8/6N1/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 2", 1, 20,
			map[string]int{
				"a7a6": 1,
				"a7a5": 1,
				"b7b6": 1,
				"b7b5": 1,
				"c7c6": 1,
				"c7c5": 1,
				"d7d6": 1,
				"d7d5": 1,
				"e7e6": 1,
				"e7e5": 1,
				"f7f6": 1,
				"f7f5": 1,
				"h7h6": 1,
				"h7h5": 1,
				"b8a6": 1,
				"b8c6": 1,
				"f8g7": 1,
				"f8h6": 1,
				"g8f6": 1,
				"g8h6": 1,
			},
		},
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
