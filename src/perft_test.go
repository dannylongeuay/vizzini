package main

import (
	"testing"
)

type DepthCount struct {
	depth int
	count int
}

func Perft(b *Board, depth int) int {
	var nodes int

	if depth == 0 {
		return 1
	}

	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	b.GenerateMoves(&moves, b.sideToMove, false)
	for _, m := range moves {
		err := b.MakeMove(m)
		if err == nil {
			nodes += Perft(b, depth-1)
		}
		b.UndoMove()
	}

	return nodes
}

func TestPerft(t *testing.T) {
	tests := []struct {
		fen         string
		depthCounts []DepthCount
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			[]DepthCount{{1, 20}, {2, 400}, {3, 8902}, {4, 197281}, {5, 4865609}, {6, 119060324}},
		},
		{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			[]DepthCount{{1, 48}, {2, 2039}, {3, 97862}, {4, 4085603}, {5, 193690690}},
		},
		{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1",
			[]DepthCount{{1, 24}, {2, 496}, {3, 9483}, {4, 182838}, {5, 3605103}, {6, 71179139}},
		},
	}
	InitHashKeys(181818)
	for _, tt := range tests {
		for _, dc := range tt.depthCounts {
			if testing.Short() && dc.count >= 100000 {
				continue
			}
			b, err := NewBoard(tt.fen)
			if err != nil {
				t.Error(err)
			}
			actual := Perft(b, dc.depth)
			if actual != dc.count {
				t.Errorf("nodes %v != %v", actual, dc.count)
			}
		}
	}
}
