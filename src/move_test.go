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
		{"rnbqkb1r/pppp2Pp/4pn2/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5", WHITE, "g7", 12,
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
		{"r1bqkbnr/ppp1p1pp/8/3p4/3n1p2/2N1P1P1/PPPPQP1P/R1B1KB1R b KQkq - 1 6", BLACK, "d4", 8,
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
