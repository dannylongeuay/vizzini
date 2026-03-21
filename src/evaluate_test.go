package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		fen      string
		expected int
	}{
		{
			STARTING_FEN,
			0,
		},
		{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			-30,
		},
		{
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			30,
		},
		{
			"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			0,
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			0,
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
			-15,
		},
		{
			"r1bq1rk1/ppppbppp/2n2n2/1B2p3/4P3/3P1N2/PPP2PPP/RNBQ1RK1 w - - 1 6",
			-10,
		},
	}
	for _, tt := range tests {
		board, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actual := board.Evaluate()

		if actual != tt.expected {
			t.Errorf("score: %v != %v\n\n%v", actual, tt.expected, board.ToString())
		}
	}
}

func TestEvaluateMaterialImbalance(t *testing.T) {
	tests := []struct {
		name      string
		fen       string
		minScore  int
		maxScore  int
	}{
		{
			name:     "white up a knight",
			fen:      "4k3/8/8/8/8/8/8/4KN2 w - - 0 1",
			minScore: 250,
			maxScore: 350,
		},
		{
			name:     "white up a queen",
			fen:      "4k3/8/8/8/8/8/8/3QK3 w - - 0 1",
			minScore: 850,
			maxScore: 950,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := NewBoard(tt.fen)
			if err != nil {
				t.Fatal(err)
			}
			score := board.Evaluate()
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("score %v not in expected range [%v, %v]", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestEvaluateSymmetry(t *testing.T) {
	// Color-flipped positions should produce negated scores.
	tests := []struct {
		name    string
		fenW    string // White to move
		fenB    string // Exact color mirror, Black to move
	}{
		{
			// White up a knight, White to move → positive.
			// Mirror: Black up a knight, White to move → negative.
			name: "knight advantage",
			fenW: "4k3/8/8/8/8/8/8/4KN2 w - - 0 1",
			fenB: "4kn2/8/8/8/8/8/8/4K3 w - - 0 1",
		},
		{
			name: "rook and pawns",
			fenW: "4k3/pppp4/8/8/8/8/PPPP4/4K2R w - - 0 1",
			fenB: "4k2r/pppp4/8/8/8/8/PPPP4/4K3 w - - 0 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boardW, err := NewBoard(tt.fenW)
			if err != nil {
				t.Fatal(err)
			}
			boardB, err := NewBoard(tt.fenB)
			if err != nil {
				t.Fatal(err)
			}
			scoreW := boardW.Evaluate()
			scoreB := boardB.Evaluate()
			if scoreW != -scoreB {
				t.Errorf("Evaluate(W)=%v, Evaluate(B)=%v, expected Evaluate(W) == -Evaluate(B)", scoreW, scoreB)
			}
		})
	}
}
