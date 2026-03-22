package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		minScore int
		maxScore int
	}{
		{
			name:     "starting position is zero",
			fen:      STARTING_FEN,
			minScore: 0,
			maxScore: 0,
		},
		{
			name:     "white e4 first move",
			fen:      "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			minScore: -68,
			maxScore: -48,
		},
		{
			name:     "sicilian defense",
			fen:      "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			minScore: 45,
			maxScore: 65,
		},
		{
			name:     "symmetric e4 e5",
			fen:      "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			minScore: -10,
			maxScore: 10,
		},
		{
			name:     "ruy lopez before castling",
			fen:      "r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			minScore: -17,
			maxScore: 3,
		},
		{
			name:     "ruy lopez after castling",
			fen:      "r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
			minScore: -27,
			maxScore: -7,
		},
		{
			name:     "ruy lopez middlegame",
			fen:      "r1bq1rk1/ppppbppp/2n2n2/1B2p3/4P3/3P1N2/PPP2PPP/RNBQ1RK1 w - - 1 6",
			minScore: -6,
			maxScore: 14,
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
			minScore: 900,
			maxScore: 1000,
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

func TestEvaluateMobility(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		minScore int
		maxScore int
	}{
		{
			// Open position: active pieces should score higher than cramped ones.
			name:     "active knight vs trapped knight",
			fen:      "4k3/8/8/4N3/8/8/8/4K3 w - - 0 1",
			minScore: 290,
			maxScore: 340,
		},
		{
			// Bishop pair bonus should be awarded.
			name:     "white bishop pair vs single bishop",
			fen:      "4k3/8/8/8/8/8/8/2B1KB2 w - - 0 1",
			minScore: 650,
			maxScore: 800,
		},
		{
			// Starting position is symmetric: full eval should cancel out.
			name:     "starting position eval is zero",
			fen:      STARTING_FEN,
			minScore: 0,
			maxScore: 0,
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

func TestEvaluateKingSafety(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		minScore int
		maxScore int
	}{
		{
			// Castled king with intact pawn shield should score better than broken shield.
			// White castled kingside with full shield vs black castled with missing g-pawn.
			name:     "intact shield vs broken shield",
			fen:      "r1bq1rk1/pppp1p1p/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 5",
			minScore: 490,
			maxScore: 540,
		},
		{
			// Multiple black attackers aimed at white king zone.
			name:     "attackers near white king",
			fen:      "r2q1rk1/ppp2ppp/2n5/3pp3/2B1n3/4PN2/PPP1QPPP/R1B2RK1 b - - 0 9",
			minScore: -230,
			maxScore: -180,
		},
		{
			// Open file near king penalty: white king on g1, no pawns on g-file.
			name:     "open file near king",
			fen:      "r1bqk2r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1P1P/RNBQ1RK1 w kq - 0 5",
			minScore: 220,
			maxScore: 270,
		},
		{
			// Starting position should still evaluate to 0 (symmetric).
			name:     "starting position is zero",
			fen:      STARTING_FEN,
			minScore: 0,
			maxScore: 0,
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
