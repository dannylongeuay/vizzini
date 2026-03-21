package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestSearchNegamax(t *testing.T) {
	tests := []struct {
		fen      string
		expected int
	}{
		{
			STARTING_FEN,
			30,
		},
	}
	for _, tt := range tests {
		search, err := NewSearch(tt.fen, DEFAULT_MAX_DEPTH, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		search.evalNoise = 0
		actual := search.Negamax(1, MIN_SCORE, MAX_SCORE)
		if actual != tt.expected {
			t.Errorf("\n%v != %v\n\n%v", actual, tt.expected, search.ToString())
		}
	}
}

func TestSearchRepetition(t *testing.T) {
	tests := []struct {
		fen      string
		mus      []MoveUnpacked
		expected bool
	}{
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
			},
			false,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET},
			},
			true,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET},
			},
			true,
		},
		{
			STARTING_FEN,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, EMPTY, QUIET},
				{B8, C6, BLACK_KNIGHT, EMPTY, QUIET},
				{E2, E4, WHITE_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, EMPTY, DOUBLE_PAWN_PUSH},
				{F3, G1, WHITE_KNIGHT, EMPTY, QUIET},
				{C6, B8, BLACK_KNIGHT, EMPTY, QUIET},
			},
			false,
		},
	}
	for _, tt := range tests {
		search, err := NewSearch(tt.fen, DEFAULT_MAX_DEPTH, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		sMoves := "moves"
		for _, mu := range tt.mus {
			move := NewMoveFromMoveUnpacked(mu)
			sMoves += " " + move.ToUCIString()
			search.Board.MakeMove(move)
		}
		actual := search.Repetition()
		if actual != tt.expected {
			t.Errorf("repetition result: %v != %v for %v", actual, tt.expected, sMoves)
		}
	}

}

func TestSearchMateInX(t *testing.T) {
	tests := []struct {
		fen         string
		movesToMate int
		expected    string
	}{
		{
			"1r5k/p6p/6p1/2pQPP2/2P3P1/3R4/q7/1N1K4 b - - 4 41",
			1,
			"pv b8b1",
		},
		{
			"3r4/2pq1p2/3k4/p5Q1/3P4/P5P1/7P/6K1 w - - 7 35",
			2,
			"pv g5c5 d6e6 c5e5",
		},
		{
			"7r/p7/2p5/1p3pk1/2pq4/6K1/PP4P1/3RQR2 b - - 2 26",
			3,
			"pv d4h4 g3f3 h4f4 f3e2 h8e8",
		},
		// {
		// 	"2Q5/p6p/4pk1p/1b1p4/5q2/1N3B1P/PPPN4/4K2n b - - 16 30",
		// 	4,
		// 	"pv f4e3 e1d1 h1f2 d1c1 e3e1 f3d1 e1d1",
		// },
	}
	for _, tt := range tests {
		if testing.Short() && tt.movesToMate > 3 {
			continue
		}
		search, err := NewSearch(tt.fen, tt.movesToMate*2, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		search.quiet = true
		search.evalNoise = 0
		search.temperature = 0
		search.IterativeDeepening()
		actual := search.GetPvLineString()
		if actual != tt.expected {
			t.Errorf("pv line: %v != %v", actual, tt.expected)
		}

	}
}

// Priority 1: Terminal nodes and draw detection

func TestSearchFiftyMoveRule(t *testing.T) {
	// KR vs K — White is winning. halfMove=99 should give a large score, halfMove=100 should be DRAW.
	fen99 := "4k3/8/8/8/8/8/8/4K2R w - - 99 60"
	fen100 := "4k3/8/8/8/8/8/8/4K2R w - - 100 60"

	search99, err := NewSearch(fen99, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search99.evalNoise = 0
	score99 := search99.Negamax(1, MIN_SCORE, MAX_SCORE)

	search100, err := NewSearch(fen100, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search100.evalNoise = 0
	score100 := search100.Negamax(1, MIN_SCORE, MAX_SCORE)

	if score100 != DRAW {
		t.Errorf("halfMove=100 should return DRAW (0), got %v", score100)
	}
	if score99 <= 100 {
		t.Errorf("halfMove=99 should return large positive score, got %v", score99)
	}
	if score99-score100 < 100 {
		t.Errorf("score difference should be dramatic: score99=%v score100=%v", score99, score100)
	}
}

func TestSearchStalemate(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		wantDraw bool
	}{
		{
			name:     "stalemate",
			fen:      "k7/2Q5/1K6/8/8/8/8/8 b - - 0 1",
			wantDraw: true,
		},
		{
			name:     "checkmate",
			fen:      "k7/1Q6/1K6/8/8/8/8/8 b - - 0 1",
			wantDraw: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			search, err := NewSearch(tt.fen, 1, DEFAULT_MAX_NODES)
			if err != nil {
				t.Fatal(err)
			}
			search.evalNoise = 0
			score := search.Negamax(1, MIN_SCORE, MAX_SCORE)
			if tt.wantDraw {
				if score != DRAW {
					t.Errorf("expected DRAW (0), got %v", score)
				}
			} else {
				// Checkmate: score should be near MIN_SCORE
				if score > MIN_SCORE+100 {
					t.Errorf("expected near MIN_SCORE for checkmate, got %v", score)
				}
			}
		})
	}
}

func TestQSearchCheckmateInQSearch(t *testing.T) {
	// Black king on h8, White queen on g7 and king on g6 — Black is in check with no moves.
	// This position will enter QSearch where the side to move is checkmated.
	fen := "7k/6Q1/6K1/8/8/8/8/8 b - - 0 1"
	search, err := NewSearch(fen, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.evalNoise = 0
	// Call QSearch directly — Black is in check, no legal moves → mate
	score := search.QSearch(MIN_SCORE, MAX_SCORE)
	if score > MIN_SCORE+10 {
		t.Errorf("QSearch should return near MIN_SCORE for checkmate, got %v", score)
	}
}

// Priority 2: Search algorithm features

func TestSearchDepthProgression(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping depth progression test in short mode")
	}
	// Mate in 2: at depth 1, engine can't find mate. At depth 4, it should.
	fen := "3r4/2pq1p2/3k4/p5Q1/3P4/P5P1/7P/6K1 w - - 7 35"

	search, err := NewSearch(fen, 4, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.quiet = true
	search.evalNoise = 0
	search.temperature = 0
	score := search.IterativeDeepening()

	if search.completedDepth != 4 {
		t.Errorf("expected completedDepth=4, got %v", search.completedDepth)
	}
	bestUCI := search.bestMove.ToUCIString()
	if bestUCI != "g5c5" {
		t.Errorf("expected best move g5c5, got %v", bestUCI)
	}
	if score <= MATE_THRESHOLD {
		t.Errorf("expected mate score > MATE_THRESHOLD, got %v", score)
	}
}

func TestSearchTimeStop(t *testing.T) {
	search, err := NewSearch(STARTING_FEN, 10, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.quiet = true
	search.evalNoise = 0
	search.temperature = 0
	// Set stopTime to expire very soon — time check fires every 2048 nodes
	search.stopTime = time.Now().Add(1 * time.Millisecond)
	search.IterativeDeepening()

	if search.completedDepth >= 10 {
		t.Errorf("expected completedDepth < maxDepth (10) with tight timer, got %v", search.completedDepth)
	}
}

func TestSearchNullMovePruning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping null move pruning test in short mode")
	}
	// King+pawns only: hasPieces is false, NMP must not fire. Should still return valid score.
	fen := "4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1"
	search, err := NewSearch(fen, 4, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.quiet = true
	search.evalNoise = 0
	search.temperature = 0
	score := search.IterativeDeepening()

	if score < MIN_SCORE+1000 || score > MAX_SCORE-1000 {
		t.Errorf("king+pawns position should return reasonable score, got %v", score)
	}

	// Normal position with pieces: NMP should keep node count under budget at depth 5.
	// Without null-move pruning, depth 5 from the starting position would blow past 500k nodes.
	search2, err := NewSearch(STARTING_FEN, 5, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search2.quiet = true
	search2.evalNoise = 0
	search2.temperature = 0
	search2.IterativeDeepening()

	if search2.nodes >= 500_000 {
		t.Errorf("expected NMP to keep nodes < 500000, got %v", search2.nodes)
	}
}

func TestSearchCheckExtension(t *testing.T) {
	// Mate-in-2: Qc5+ gives check, so the opponent's reply gets depth-extended,
	// allowing the engine to see Qe5# at depth 2. Without check extension, depth 2 wouldn't find it.
	fen := "3r4/2pq1p2/3k4/p5Q1/3P4/P5P1/7P/6K1 w - - 7 35"
	search, err := NewSearch(fen, 2, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.quiet = true
	search.evalNoise = 0
	search.temperature = 0
	score := search.IterativeDeepening()

	bestUCI := search.bestMove.ToUCIString()
	if bestUCI != "g5c5" {
		t.Errorf("expected g5c5 (Qc5+ starting mate-in-2), got %v", bestUCI)
	}
	if score <= MATE_THRESHOLD {
		t.Errorf("expected mate score > MATE_THRESHOLD, got %v", score)
	}
}

func TestQSearchStandingPat(t *testing.T) {
	// Quiet winning position: KR vs K, no captures available. QSearch ≈ Evaluate().
	fen := "4k3/8/8/8/8/8/8/4K2R w - - 0 1"
	search, err := NewSearch(fen, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.evalNoise = 0
	staticEval := search.Evaluate()
	qScore := search.QSearch(MIN_SCORE, MAX_SCORE)

	// QSearch may find some quiet moves via check evasions, but should be close to static eval
	diff := qScore - staticEval
	if diff < -100 || diff > 100 {
		t.Errorf("QSearch score (%v) should be close to static eval (%v)", qScore, staticEval)
	}

	// Position where a capture improves score: White can capture Black's queen
	fenCapture := "4k3/8/8/3q4/4R3/8/8/4K3 w - - 0 1"
	searchCapture, err := NewSearch(fenCapture, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	searchCapture.evalNoise = 0
	staticCapture := searchCapture.Evaluate()
	qScoreCapture := searchCapture.QSearch(MIN_SCORE, MAX_SCORE)

	if qScoreCapture <= staticCapture {
		t.Errorf("QSearch with capture available (%v) should exceed standing pat (%v)", qScoreCapture, staticCapture)
	}
}

// Priority 3: Move ordering and heuristics

func TestMoveOrderingMVVLVA(t *testing.T) {
	// White rook on e4 can capture Black queen on e7 (RxQ) or Black pawn on a4 (RxP).
	// MVV-LVA should prefer capturing the queen first.
	fen := "4k3/4q3/8/8/p3R3/8/8/4K3 w - - 0 1"
	search, err := NewSearch(fen, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}

	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	search.GenerateMoves(&moves, WHITE, true) // captures only

	if len(moves) == 0 {
		t.Fatal("expected capture moves")
	}

	// Pick first move — should be the highest-value capture
	first := search.PickNextMove(0, &moves, 0)
	var mu MoveUnpacked
	first.Unpack(&mu)

	// The first capture picked should target the queen (d6), not the pawn (e5)
	if mu.dstSquare != BLACK_QUEEN {
		t.Errorf("expected first capture to target queen, got %v", mu.dstSquare)
	}
}

func TestKillerMoveUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping killer move test in short mode")
	}
	search, err := NewSearch(STARTING_FEN, 4, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.quiet = true
	search.evalNoise = 0
	search.temperature = 0
	search.IterativeDeepening()

	hasKiller := false
	for d := 0; d < KILLERS_DEPTH; d++ {
		if search.killers[0][d] != 0 {
			hasKiller = true
			break
		}
	}
	if !hasKiller {
		t.Error("expected at least one killer move to be set after depth-4 search")
	}
}

func TestHistoryHeuristicAccumulation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping history heuristic test in short mode")
	}
	search, err := NewSearch(STARTING_FEN, 4, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}
	search.quiet = true
	search.evalNoise = 0
	search.temperature = 0
	search.IterativeDeepening()

	hasHistory := false
	for from := 0; from < BOARD_SQUARES; from++ {
		for to := 0; to < BOARD_SQUARES; to++ {
			val := search.alphaHistory[from][to]
			if val > MAX_HISTORY_VALUE {
				t.Errorf("history value %v exceeds MAX_HISTORY_VALUE %v at [%v][%v]", val, MAX_HISTORY_VALUE, from, to)
			}
			if val > 0 {
				hasHistory = true
			}
		}
	}
	if !hasHistory {
		t.Error("expected some non-zero history entries after depth-4 search")
	}
}

// Priority 4: Evaluation edge cases (NoisyEvaluate and TemperatureSelect)

func TestNoisyEvaluate(t *testing.T) {
	search, err := NewSearch(STARTING_FEN, 1, DEFAULT_MAX_NODES)
	if err != nil {
		t.Fatal(err)
	}

	// With evalNoise=0, NoisyEvaluate must equal Evaluate exactly
	search.evalNoise = 0
	baseScore := search.Evaluate()
	noisy := search.NoisyEvaluate()
	if noisy != baseScore {
		t.Errorf("evalNoise=0: NoisyEvaluate (%v) != Evaluate (%v)", noisy, baseScore)
	}

	// With evalNoise=10, output should be within range and have variance
	search.evalNoise = 10
	search.rng = rand.New(rand.NewSource(42))
	min, max := baseScore, baseScore
	for i := 0; i < 200; i++ {
		s := search.NoisyEvaluate()
		if s < baseScore-10 || s > baseScore+10 {
			t.Errorf("NoisyEvaluate (%v) outside range [%v, %v]", s, baseScore-10, baseScore+10)
		}
		if s < min {
			min = s
		}
		if s > max {
			max = s
		}
	}
	if max-min < 2 {
		t.Errorf("expected variance in noisy eval, but range was only %v", max-min)
	}
}

func TestTemperatureSelect(t *testing.T) {
	t.Run("temperature=0 preserves bestMove", func(t *testing.T) {
		search, err := NewSearch(STARTING_FEN, 3, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		search.quiet = true
		search.evalNoise = 0
		search.temperature = 0
		bestScore := search.IterativeDeepening()
		originalBest := search.bestMove

		// Manually call TemperatureSelect — should be a no-op
		search.TemperatureSelect(bestScore)
		if search.bestMove != originalBest {
			t.Errorf("temperature=0 should not change bestMove")
		}
	})

	t.Run("fullMove > tempMoveLimit skips", func(t *testing.T) {
		search, err := NewSearch(STARTING_FEN, 3, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		search.quiet = true
		search.evalNoise = 0
		search.temperature = 1.0
		bestScore := search.IterativeDeepening()
		originalBest := search.bestMove
		search.fullMove = search.tempMoveLimit + 1
		search.TemperatureSelect(bestScore)
		if search.bestMove != originalBest {
			t.Error("should skip when fullMove > tempMoveLimit")
		}
	})

	t.Run("mate score skips", func(t *testing.T) {
		search, err := NewSearch(STARTING_FEN, 3, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		search.quiet = true
		search.evalNoise = 0
		search.temperature = 1.0
		search.IterativeDeepening()
		originalBest := search.bestMove

		search.TemperatureSelect(MAX_SCORE - 5) // mate score
		if search.bestMove != originalBest {
			t.Error("should skip for mate scores")
		}
	})

	t.Run("temperature>0 sometimes selects non-best", func(t *testing.T) {
		search, err := NewSearch(STARTING_FEN, 3, DEFAULT_MAX_NODES)
		if err != nil {
			t.Fatal(err)
		}
		search.quiet = true
		search.evalNoise = 0
		search.temperature = 0
		bestScore := search.IterativeDeepening()
		originalBest := search.bestMove

		if len(search.rootMoves) <= 1 {
			t.Skip("need multiple root moves for temperature test")
		}

		changed := false
		for i := 0; i < 100; i++ {
			search.bestMove = originalBest
			search.temperature = 5.0 // high temperature for more randomness
			search.rng = rand.New(rand.NewSource(int64(i)))
			search.TemperatureSelect(bestScore)
			if search.bestMove != originalBest {
				changed = true
				break
			}
		}
		if !changed {
			t.Error("with high temperature, expected at least one different move selection in 100 tries")
		}
	})
}
