package main

import (
	"testing"
)

func containsTestMove(moves []move, testMove testMove) (bool, error) {
	o, err := squareIndexByCoord(testMove.origin)
	if err != nil {
		return false, err
	}
	t, err := squareIndexByCoord(testMove.target)
	if err != nil {
		return false, err
	}
	m := move{
		origin: o,
		target: t,
		kind:   testMove.kind,
	}
	for _, move := range moves {
		if move == m {
			return true, nil
		}
	}
	return false, nil
}

type testMove struct {
	origin SquareCoord
	target SquareCoord
	kind   MoveKind
}

func TestGenerateMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, 20,
			[]testMove{
				{"a2", "a3", QUIET},
				{"b2", "b3", QUIET},
				{"c2", "c3", QUIET},
				{"d2", "d3", QUIET},
				{"e2", "e3", QUIET},
				{"f2", "f3", QUIET},
				{"g2", "g3", QUIET},
				{"h2", "h3", QUIET},
				{"a2", "a4", DOUBLE_PAWN_PUSH},
				{"b2", "b4", DOUBLE_PAWN_PUSH},
				{"c2", "c4", DOUBLE_PAWN_PUSH},
				{"d2", "d4", DOUBLE_PAWN_PUSH},
				{"e2", "e4", DOUBLE_PAWN_PUSH},
				{"f2", "f4", DOUBLE_PAWN_PUSH},
				{"g2", "g4", DOUBLE_PAWN_PUSH},
				{"h2", "h4", DOUBLE_PAWN_PUSH},
				{"g1", "f3", QUIET},
				{"g1", "h3", QUIET},
				{"b1", "a3", QUIET},
				{"b1", "c3", QUIET},
			},
		},
		{STARTING_FEN, BLACK, 20,
			[]testMove{
				{"a7", "a6", QUIET},
				{"b7", "b6", QUIET},
				{"c7", "c6", QUIET},
				{"d7", "d6", QUIET},
				{"e7", "e6", QUIET},
				{"f7", "f6", QUIET},
				{"g7", "g6", QUIET},
				{"h7", "h6", QUIET},
				{"a7", "a5", DOUBLE_PAWN_PUSH},
				{"b7", "b5", DOUBLE_PAWN_PUSH},
				{"c7", "c5", DOUBLE_PAWN_PUSH},
				{"d7", "d5", DOUBLE_PAWN_PUSH},
				{"e7", "e5", DOUBLE_PAWN_PUSH},
				{"f7", "f5", DOUBLE_PAWN_PUSH},
				{"g7", "g5", DOUBLE_PAWN_PUSH},
				{"h7", "h5", DOUBLE_PAWN_PUSH},
				{"g8", "f6", QUIET},
				{"g8", "h6", QUIET},
				{"b8", "a6", QUIET},
				{"b8", "c6", QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := b.generateMoves(tt.color)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGeneratePawnMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareCoord SquareCoord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, "c2", 2,
			[]testMove{
				{"c2", "c4", DOUBLE_PAWN_PUSH},
				{"c2", "c3", QUIET},
			},
		},
		{"rnbqkbnr/pp2pppp/8/2ppP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3", WHITE, "e5", 2,
			[]testMove{
				{"e5", "e6", QUIET},
				{"e5", "d6", EP_CAPTURE},
			},
		},
		{"rnbqkbnr/pp2pppp/3P4/2p5/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3", BLACK, "e7", 3,
			[]testMove{
				{"e7", "e5", DOUBLE_PAWN_PUSH},
				{"e7", "e6", QUIET},
				{"e7", "d6", CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P4/PPP2PPP/RNBQKBNR w KQkq - 0 3", WHITE, "e4", 3,
			[]testMove{
				{"e4", "e5", QUIET},
				{"e4", "d5", CAPTURE},
				{"e4", "f5", CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P1P2/PPP3PP/RNBQKBNR b KQkq - 0 3", BLACK, "d5", 2,
			[]testMove{
				{"d5", "d4", QUIET},
				{"d5", "e4", CAPTURE},
			},
		},
		{"rnbqkb1r/pppp2Pp/4pn2/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5", WHITE, "g7", MAX_PAWN_MOVES,
			[]testMove{
				{"g7", "f8", KNIGHT_PROMOTION_CAPTURE},
				{"g7", "f8", BISHOP_PROMOTION_CAPTURE},
				{"g7", "f8", ROOK_PROMOTION_CAPTURE},
				{"g7", "f8", QUEEN_PROMOTION_CAPTURE},
				{"g7", "g8", KNIGHT_PROMOTION},
				{"g7", "g8", BISHOP_PROMOTION},
				{"g7", "g8", ROOK_PROMOTION},
				{"g7", "g8", QUEEN_PROMOTION},
				{"g7", "h8", KNIGHT_PROMOTION_CAPTURE},
				{"g7", "h8", BISHOP_PROMOTION_CAPTURE},
				{"g7", "h8", ROOK_PROMOTION_CAPTURE},
				{"g7", "h8", QUEEN_PROMOTION_CAPTURE},
			},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/PpPP1PPP/RNBQ1RK1 b q - 1 9", BLACK, "b2", 8,
			[]testMove{
				{"b2", "a1", KNIGHT_PROMOTION_CAPTURE},
				{"b2", "a1", BISHOP_PROMOTION_CAPTURE},
				{"b2", "a1", ROOK_PROMOTION_CAPTURE},
				{"b2", "a1", QUEEN_PROMOTION_CAPTURE},
				{"b2", "c1", KNIGHT_PROMOTION_CAPTURE},
				{"b2", "c1", BISHOP_PROMOTION_CAPTURE},
				{"b2", "c1", ROOK_PROMOTION_CAPTURE},
				{"b2", "c1", QUEEN_PROMOTION_CAPTURE},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		moves := b.generatePawnMoves(tt.color, squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareCoord SquareCoord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, "g1", 2,
			[]testMove{
				{"g1", "f3", QUIET},
				{"g1", "h3", QUIET},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/3N4/8/PPPPPPPP/RNBQKB1R w KQkq - 0 3", WHITE, "d4", 6,
			[]testMove{
				{"d4", "b3", QUIET},
				{"d4", "b5", QUIET},
				{"d4", "c6", QUIET},
				{"d4", "e6", QUIET},
				{"d4", "f5", CAPTURE},
				{"d4", "f3", QUIET},
			},
		},
		{"r1bqkbnr/ppp1p1pp/8/3p4/3n1p2/2N1P1P1/PPPPQP1P/R1B1KB1R b KQkq - 1 6", BLACK, "d4", MAX_KNIGHT_MOVES,
			[]testMove{
				{"d4", "b3", QUIET},
				{"d4", "b5", QUIET},
				{"d4", "c6", QUIET},
				{"d4", "e6", QUIET},
				{"d4", "f5", QUIET},
				{"d4", "f3", QUIET},
				{"d4", "e2", CAPTURE},
				{"d4", "c2", CAPTURE},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		moves := b.generateKnightMoves(tt.color, squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareCoord SquareCoord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, "f1", 0,
			[]testMove{},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", WHITE, "f1", 5,
			[]testMove{
				{"f1", "e2", QUIET},
				{"f1", "d3", QUIET},
				{"f1", "c4", QUIET},
				{"f1", "b5", QUIET},
				{"f1", "a6", QUIET},
			},
		},
		{"rnbqkbnr/pp1pp1pp/5p2/2p5/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 3", WHITE, "c4", 10,
			[]testMove{
				{"c4", "b3", QUIET},
				{"c4", "b5", QUIET},
				{"c4", "a6", QUIET},
				{"c4", "d3", QUIET},
				{"c4", "e2", QUIET},
				{"c4", "f1", QUIET},
				{"c4", "d5", QUIET},
				{"c4", "e6", QUIET},
				{"c4", "f7", QUIET},
				{"c4", "g8", CAPTURE},
			},
		},
		{"rn2kbnr/p2qp1p1/1p1p1p1p/1bpBP3/7P/P5P1/1PPP1P1R/RNBQK1N1 w Qkq - 3 9", WHITE, "d5", MAX_BISHOP_MOVES,
			[]testMove{
				{"d5", "c6", QUIET},
				{"d5", "b7", QUIET},
				{"d5", "a8", CAPTURE},
				{"d5", "e6", QUIET},
				{"d5", "f7", QUIET},
				{"d5", "g8", CAPTURE},
				{"d5", "c4", QUIET},
				{"d5", "b3", QUIET},
				{"d5", "a2", QUIET},
				{"d5", "e4", QUIET},
				{"d5", "f3", QUIET},
				{"d5", "g2", QUIET},
				{"d5", "h1", QUIET},
			}},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		moves := b.generateBishopMoves(tt.color, squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGenerateRookMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareCoord SquareCoord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]testMove{},
		},
		{"rnbqkbnr/2pppppp/8/8/8/R7/1PPPPPPP/1NBQKBNR w Kk - 2 5", WHITE, "a3", MAX_ROOK_MOVES,
			[]testMove{
				{"a3", "a1", QUIET},
				{"a3", "a2", QUIET},
				{"a3", "a4", QUIET},
				{"a3", "a5", QUIET},
				{"a3", "a6", QUIET},
				{"a3", "a7", QUIET},
				{"a3", "a8", CAPTURE},
				{"a3", "b3", QUIET},
				{"a3", "c3", QUIET},
				{"a3", "d3", QUIET},
				{"a3", "e3", QUIET},
				{"a3", "f3", QUIET},
				{"a3", "g3", QUIET},
				{"a3", "h3", QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		moves := b.generateRookMoves(tt.color, squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareCoord SquareCoord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, "d1", 0,
			[]testMove{},
		},
		{"r2qk1r1/p1p3pp/np3p1n/3Q4/6bP/P5P1/1PP1PP1R/RNB1KBN1 w Qq - 1 10", WHITE, "d5", MAX_QUEEN_MOVES,
			[]testMove{
				{"d5", "c6", QUIET},
				{"d5", "b7", QUIET},
				{"d5", "a8", CAPTURE},
				{"d5", "e6", QUIET},
				{"d5", "f7", QUIET},
				{"d5", "g8", CAPTURE},
				{"d5", "c4", QUIET},
				{"d5", "b3", QUIET},
				{"d5", "a2", QUIET},
				{"d5", "e4", QUIET},
				{"d5", "f3", QUIET},
				{"d5", "g2", QUIET},
				{"d5", "h1", QUIET},
				{"d5", "d1", QUIET},
				{"d5", "d2", QUIET},
				{"d5", "d3", QUIET},
				{"d5", "d4", QUIET},
				{"d5", "d6", QUIET},
				{"d5", "d7", QUIET},
				{"d5", "d8", CAPTURE},
				{"d5", "a5", QUIET},
				{"d5", "b5", QUIET},
				{"d5", "c5", QUIET},
				{"d5", "e5", QUIET},
				{"d5", "f5", QUIET},
				{"d5", "g5", QUIET},
				{"d5", "h5", QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		moves := b.generateQueenMoves(tt.color, squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGenerateKingMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareCoord SquareCoord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, "e1", 0,
			[]testMove{},
		},
		{"rn1q2nr/p2kb1pp/1p6/4p1p1/4K3/8/PPPP1P1P/RNBQ1BNR w - - 0 10", WHITE, "e4", MAX_KING_MOVES,
			[]testMove{
				{"e4", "d3", QUIET},
				{"e4", "d4", QUIET},
				{"e4", "d5", QUIET},
				{"e4", "e3", QUIET},
				{"e4", "e5", CAPTURE},
				{"e4", "f3", QUIET},
				{"e4", "f4", QUIET},
				{"e4", "f5", QUIET},
			},
		},
		{"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4", WHITE, "e1", 3,
			[]testMove{
				{"e1", "e2", QUIET},
				{"e1", "f1", QUIET},
				{"e1", "g1", KING_CASTLE},
			},
		},
		{"r1bqk2r/ppppbppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 6 5", WHITE, "g1", 1,
			[]testMove{
				{"g1", "h1", QUIET},
			},
		},
		{"r1bqkb1r/ppp2ppp/2np1n2/8/4P3/2NQ4/PPP2PPP/R1B1KBNR w KQkq - 2 6", WHITE, "e1", 3,
			[]testMove{
				{"e1", "d1", QUIET},
				{"e1", "d2", QUIET},
				{"e1", "e2", QUIET},
			},
		},
		{"r1bqk2r/ppp1bppp/2np1n2/6B1/4P3/2NQ4/PPP2PPP/R3KBNR w KQkq - 4 7", WHITE, "e1", 4,
			[]testMove{
				{"e1", "d1", QUIET},
				{"e1", "d2", QUIET},
				{"e1", "e2", QUIET},
				{"e1", "c1", QUEEN_CASTLE},
			},
		},
		{"r1bq1rk1/ppp1bppp/2np1n2/6B1/4P3/2NQ4/PPP2PPP/2KR1BNR w - - 6 8", WHITE, "c1", 2,
			[]testMove{
				{"c1", "b1", QUIET},
				{"c1", "d2", QUIET},
			},
		},
		{"rnbqk2r/pppp1ppp/4pn2/8/1Q1P4/5N2/PPP1PPPP/RNB1KB1R b KQkq - 2 4", BLACK, "e8", 2,
			[]testMove{
				{"e8", "f8", QUIET},
				{"e8", "e7", QUIET},
			},
		},
		{"rnbqk2r/pppp2pp/4pp2/4N3/1Q1P4/4P2n/PPP1BPPP/RNB1K2R w KQkq - 0 8", WHITE, "e1", 3,
			[]testMove{
				{"e1", "d1", QUIET},
				{"e1", "d2", QUIET},
				{"e1", "f1", QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		moves := b.generateKingMoves(tt.color, squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			contains, err := containsTestMove(moves, testMove)
			if err != nil {
				t.Error(err)
			}
			if !contains {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}
