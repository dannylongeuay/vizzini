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
		coord := fmt.Sprint(coordBySquareIndex(m.origin), coordBySquareIndex(m.target))
		nodes := perft(&cb, depth-1)
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
	}
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

// b1a3 	 1 	 20 	 400 	 8885 	 198572 	 4856835
// b1c3 	 1 	 20 	 440 	 9755 	 234656 	 5708064
// g1f3 	 1 	 20 	 440 	 9748 	 233491 	 5723523
// g1h3 	 1 	 20 	 400 	 8881 	 198502 	 4877234
// a2a4 	 1 	 20 	 420 	 9329 	 217832 	 5363555
// a2a3 	 1 	 20 	 380 	 8457 	 181046 	 4463267
// b2b4 	 1 	 20 	 421 	 9332 	 216145 	 5293555
// b2b3 	 1 	 20 	 420 	 9345 	 215255 	 5310358
// c2c4 	 1 	 20 	 441 	 9744 	 240082 	 5866666
// c2c3 	 1 	 20 	 420 	 9272 	 222861 	 5417640
// d2d4 	 1 	 20 	 560 	 12435 	 361790 	 8879566
// d2d3 	 1 	 20 	 539 	 11959 	 328511 	 8073082
// e2e4 	 1 	 20 	 600 	 13160 	 405385 	 9771632
// e2e3 	 1 	 20 	 599 	 13134 	 402988 	 9726018
// f2f4 	 1 	 20 	 401 	 8929 	 198473 	 4890429
// f2f3 	 1 	 20 	 380 	 8457 	 178889 	 4404141
// g2g4 	 1 	 20 	 421 	 9328 	 214048 	 5239875
// g2g3 	 1 	 20 	 420 	 9345 	 217210 	 5346260
// h2h4 	 1 	 20 	 420 	 9329 	 218829 	 5385554
// h2h3 	 1 	 20  	 380 	 8457 	 181044 	 4463070

func TestDivide(t *testing.T) {
	tests := []struct {
		fen      string
		depth    int
		expected map[string]int
	}{
		{STARTING_FEN, 3,
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
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		results := divide(b, tt.depth)
		for coord, nodes := range tt.expected {
			n, ok := results[coord]
			if !ok {
				t.Errorf("coord %v not found", coord)
			}
			if n != nodes {
				t.Errorf("coord %v nodes: %v != %v", coord, n, nodes)
			}
		}
	}
}
