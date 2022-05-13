package main

import (
	"testing"
	"time"
)

func newMoveFromTestMove(testMove testMove) (Move, error) {
	var m Move
	return m, nil
}

func containsTestMove(moves []Move, testMove testMove) (bool, error) {
	m, err := newMoveFromTestMove(testMove)

	if err != nil {
		return false, err
	}

	for _, move := range moves {
		if move == m {
			return true, nil
		}
	}

	return false, nil
}

func compareBoardState(bStart *Board, bEnd *Board, t *testing.T) {
	if bStart.whiteKingCoord != bEnd.whiteKingCoord {
		t.Errorf("board white king index: %v != %v", bStart.whiteKingCoord, bEnd.whiteKingCoord)
	}
	if bStart.blackKingCoord != bEnd.blackKingCoord {
		t.Errorf("board black king index: %v != %v", bStart.blackKingCoord, bEnd.blackKingCoord)
	}
	if bStart.sideToMove != bEnd.sideToMove {
		t.Errorf("board side to move: %v != %v", bStart.sideToMove, bEnd.sideToMove)
	}
	if bStart.castleRights != bEnd.castleRights {
		t.Errorf("board castle rights: %v != %v", bStart.castleRights, bEnd.castleRights)
	}
	if bStart.epCoord != bEnd.epCoord {
		t.Errorf("board en passant index: %v != %v", bStart.epCoord, bEnd.epCoord)
	}
	if bStart.halfMove != bEnd.halfMove {
		t.Errorf("board half move: %v != %v", bStart.halfMove, bEnd.halfMove)
	}
	if bStart.fullMove != bEnd.fullMove {
		t.Errorf("board full move: %v != %v", bStart.fullMove, bEnd.fullMove)
	}
	if bStart.hash != bEnd.hash {
		t.Errorf("board hash: %v != %v", bStart.hash, bEnd.hash)
	}

}

type testMove struct {
	origin Coord
	target Coord
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
				{A2, A3, QUIET},
				{B2, B3, QUIET},
				{C2, C3, QUIET},
				{D2, D3, QUIET},
				{E2, E3, QUIET},
				{F2, F3, QUIET},
				{G2, G3, QUIET},
				{H2, H3, QUIET},
				{A2, A4, DOUBLE_PAWN_PUSH},
				{B2, B4, DOUBLE_PAWN_PUSH},
				{C2, C4, DOUBLE_PAWN_PUSH},
				{D2, D4, DOUBLE_PAWN_PUSH},
				{E2, E4, DOUBLE_PAWN_PUSH},
				{F2, F4, DOUBLE_PAWN_PUSH},
				{G2, G4, DOUBLE_PAWN_PUSH},
				{H2, H4, DOUBLE_PAWN_PUSH},
				{G1, F3, QUIET},
				{G1, H3, QUIET},
				{B1, A3, QUIET},
				{B1, C3, QUIET},
			},
		},
		{STARTING_FEN, BLACK, 20,
			[]testMove{
				{A7, A6, QUIET},
				{B7, B6, QUIET},
				{C7, C6, QUIET},
				{D7, D6, QUIET},
				{E7, E6, QUIET},
				{F7, F6, QUIET},
				{G7, G6, QUIET},
				{H7, H6, QUIET},
				{A7, A5, DOUBLE_PAWN_PUSH},
				{B7, B5, DOUBLE_PAWN_PUSH},
				{C7, C5, DOUBLE_PAWN_PUSH},
				{D7, D5, DOUBLE_PAWN_PUSH},
				{E7, E5, DOUBLE_PAWN_PUSH},
				{F7, F5, DOUBLE_PAWN_PUSH},
				{G7, G5, DOUBLE_PAWN_PUSH},
				{H7, H5, DOUBLE_PAWN_PUSH},
				{G8, F6, QUIET},
				{G8, H6, QUIET},
				{B8, A6, QUIET},
				{B8, C6, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
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
		squareCoord Coord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, C2, 2,
			[]testMove{
				{C2, C4, DOUBLE_PAWN_PUSH},
				{C2, C3, QUIET},
			},
		},
		{"rnbqkbnr/pp2pppp/8/2ppP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3", WHITE, E5, 2,
			[]testMove{
				{E5, E6, QUIET},
				{E5, D6, EP_CAPTURE},
			},
		},
		{"rnbqkbnr/pp2pppp/3P4/2p5/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3", BLACK, E7, 3,
			[]testMove{
				{E7, E5, DOUBLE_PAWN_PUSH},
				{E7, E6, QUIET},
				{E7, D6, CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P4/PPP2PPP/RNBQKBNR w KQkq - 0 3", WHITE, E4, 3,
			[]testMove{
				{E4, E5, QUIET},
				{E4, D5, CAPTURE},
				{E4, F5, CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P1P2/PPP3PP/RNBQKBNR b KQkq - 0 3", BLACK, D5, 2,
			[]testMove{
				{D5, D4, QUIET},
				{D5, E4, CAPTURE},
			},
		},
		{"rnbqkb1r/pppp2Pp/4pn2/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5", WHITE, G7, MAX_PAWN_MOVES,
			[]testMove{
				{G7, F8, KNIGHT_PROMOTION_CAPTURE},
				{G7, F8, BISHOP_PROMOTION_CAPTURE},
				{G7, F8, ROOK_PROMOTION_CAPTURE},
				{G7, F8, QUEEN_PROMOTION_CAPTURE},
				{G7, G8, KNIGHT_PROMOTION},
				{G7, G8, BISHOP_PROMOTION},
				{G7, G8, ROOK_PROMOTION},
				{G7, G8, QUEEN_PROMOTION},
				{G7, H8, KNIGHT_PROMOTION_CAPTURE},
				{G7, H8, BISHOP_PROMOTION_CAPTURE},
				{G7, H8, ROOK_PROMOTION_CAPTURE},
				{G7, H8, QUEEN_PROMOTION_CAPTURE},
			},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/PpPP1PPP/RNBQ1RK1 b q - 1 9", BLACK, B2, 8,
			[]testMove{
				{B2, A1, KNIGHT_PROMOTION_CAPTURE},
				{B2, A1, BISHOP_PROMOTION_CAPTURE},
				{B2, A1, ROOK_PROMOTION_CAPTURE},
				{B2, A1, QUEEN_PROMOTION_CAPTURE},
				{B2, C1, KNIGHT_PROMOTION_CAPTURE},
				{B2, C1, BISHOP_PROMOTION_CAPTURE},
				{B2, C1, ROOK_PROMOTION_CAPTURE},
				{B2, C1, QUEEN_PROMOTION_CAPTURE},
			},
		},
		{"rnbqkbnr/1ppppppp/p7/8/8/2N5/PPPPPPPP/R1BQKBNR w KQkq - 0 2", WHITE, C2, 0,
			[]testMove{},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_PAWN_MOVES)
		b.generatePawnMoves(&moves, tt.color)
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
		squareCoord Coord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, G1, 2,
			[]testMove{
				{G1, F3, QUIET},
				{G1, H3, QUIET},
			}},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/3N4/8/PPPPPPPP/RNBQKB1R w KQkq - 0 3", WHITE, D4, 6,
			[]testMove{
				{D4, B3, QUIET},
				{D4, B5, QUIET},
				{D4, C6, QUIET},
				{D4, E6, QUIET},
				{D4, F5, CAPTURE},
				{D4, F3, QUIET},
			},
		},
		{"r1bqkbnr/ppp1p1pp/8/3p4/3n1p2/2N1P1P1/PPPPQP1P/R1B1KB1R b KQkq - 1 6", BLACK, D4, MAX_KNIGHT_MOVES,
			[]testMove{
				{D4, B3, QUIET},
				{D4, B5, QUIET},
				{D4, C6, QUIET},
				{D4, E6, QUIET},
				{D4, F5, QUIET},
				{D4, F3, QUIET},
				{D4, E2, CAPTURE},
				{D4, C2, CAPTURE},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_KNIGHT_MOVES)
		b.generateKnightMoves(&moves, tt.color)
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
		squareCoord Coord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, F1, 0,
			[]testMove{},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", WHITE, F1, 5,
			[]testMove{
				{F1, E2, QUIET},
				{F1, D3, QUIET},
				{F1, C4, QUIET},
				{F1, B5, QUIET},
				{F1, A6, QUIET},
			},
		},
		{"rnbqkbnr/pp1pp1pp/5p2/2p5/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 3", WHITE, C4, 10,
			[]testMove{
				{C4, B3, QUIET},
				{C4, B5, QUIET},
				{C4, A6, QUIET},
				{C4, D3, QUIET},
				{C4, E2, QUIET},
				{C4, F1, QUIET},
				{C4, D5, QUIET},
				{C4, E6, QUIET},
				{C4, F7, QUIET},
				{C4, G8, CAPTURE},
			},
		},
		{"rn2kbnr/p2qp1p1/1p1p1p1p/1bpBP3/7P/P5P1/1PPP1P1R/RNBQK1N1 w Qkq - 3 9", WHITE, D5, MAX_BISHOP_MOVES,
			[]testMove{
				{D5, C6, QUIET},
				{D5, B7, QUIET},
				{D5, A8, CAPTURE},
				{D5, E6, QUIET},
				{D5, F7, QUIET},
				{D5, G8, CAPTURE},
				{D5, C4, QUIET},
				{D5, B3, QUIET},
				{D5, A2, QUIET},
				{D5, E4, QUIET},
				{D5, F3, QUIET},
				{D5, G2, QUIET},
				{D5, H1, QUIET},
			}},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_BISHOP_MOVES)
		b.generateBishopMoves(&moves, tt.color)
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
		squareCoord Coord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, A1, 0,
			[]testMove{},
		},
		{"rnbqkbnr/2pppppp/8/8/8/R7/1PPPPPPP/1NBQKBNR w Kk - 2 5", WHITE, A3, MAX_ROOK_MOVES,
			[]testMove{
				{A3, A1, QUIET},
				{A3, A2, QUIET},
				{A3, A4, QUIET},
				{A3, A5, QUIET},
				{A3, A6, QUIET},
				{A3, A7, QUIET},
				{A3, A8, CAPTURE},
				{A3, B3, QUIET},
				{A3, C3, QUIET},
				{A3, D3, QUIET},
				{A3, E3, QUIET},
				{A3, F3, QUIET},
				{A3, G3, QUIET},
				{A3, H3, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_ROOK_MOVES)
		b.generateRookMoves(&moves, tt.color)
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
		squareCoord Coord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, D1, 0,
			[]testMove{},
		},
		{"r2qk1r1/p1p3pp/np3p1n/3Q4/6bP/P5P1/1PP1PP1R/RNB1KBN1 w Qq - 1 10", WHITE, D5, MAX_QUEEN_MOVES,
			[]testMove{
				{D5, C6, QUIET},
				{D5, B7, QUIET},
				{D5, A8, CAPTURE},
				{D5, E6, QUIET},
				{D5, F7, QUIET},
				{D5, G8, CAPTURE},
				{D5, C4, QUIET},
				{D5, B3, QUIET},
				{D5, A2, QUIET},
				{D5, E4, QUIET},
				{D5, F3, QUIET},
				{D5, G2, QUIET},
				{D5, H1, QUIET},
				{D5, D1, QUIET},
				{D5, D2, QUIET},
				{D5, D3, QUIET},
				{D5, D4, QUIET},
				{D5, D6, QUIET},
				{D5, D7, QUIET},
				{D5, D8, CAPTURE},
				{D5, A5, QUIET},
				{D5, B5, QUIET},
				{D5, C5, QUIET},
				{D5, E5, QUIET},
				{D5, F5, QUIET},
				{D5, G5, QUIET},
				{D5, H5, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_QUEEN_MOVES)
		b.generateQueenMoves(&moves, tt.color)
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
		squareCoord Coord
		movesLength int
		moves       []testMove
	}{
		{STARTING_FEN, WHITE, E1, 0,
			[]testMove{},
		},
		{"rn1q2nr/p2kb1pp/1p6/8/4K3/2N2p1N/PPPP4/R1BQ1B1R w - - 0 14", WHITE, E4, MAX_KING_MOVES,
			[]testMove{
				{E4, D3, QUIET},
				{E4, D4, QUIET},
				{E4, D5, QUIET},
				{E4, E3, QUIET},
				{E4, E5, QUIET},
				{E4, F3, CAPTURE},
				{E4, F4, QUIET},
				{E4, F5, QUIET},
			},
		},
		// TODO: add illegal moves
		{"rn1q2nr/p2kb1pp/1p6/4p1p1/4K3/8/PPPP1P1P/RNBQ1BNR w - - 0 10", WHITE, E4, 8,
			[]testMove{
				{E4, D3, QUIET},
				{E4, D5, QUIET},
				{E4, E3, QUIET},
				{E4, E5, CAPTURE},
				{E4, F3, QUIET},
				{E4, F5, QUIET},
			},
		},
		{"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4", WHITE, E1, 3,
			[]testMove{
				{E1, E2, QUIET},
				{E1, F1, QUIET},
				{E1, G1, KING_CASTLE},
			},
		},
		{"r1bqk2r/ppppbppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 6 5", WHITE, G1, 1,
			[]testMove{
				{G1, H1, QUIET},
			},
		},
		{"r1bqkb1r/ppp2ppp/2np1n2/8/4P3/2NQ4/PPP2PPP/R1B1KBNR w KQkq - 2 6", WHITE, E1, 3,
			[]testMove{
				{E1, D1, QUIET},
				{E1, D2, QUIET},
				{E1, E2, QUIET},
			},
		},
		{"r1bqk2r/ppp1bppp/2np1n2/6B1/4P3/2NQ4/PPP2PPP/R3KBNR w KQkq - 4 7", WHITE, E1, 4,
			[]testMove{
				{E1, D1, QUIET},
				{E1, D2, QUIET},
				{E1, E2, QUIET},
				{E1, C1, QUEEN_CASTLE},
			},
		},
		{"r1bq1rk1/ppp1bppp/2np1n2/6B1/4P3/2NQ4/PPP2PPP/2KR1BNR w - - 6 8", WHITE, C1, 2,
			[]testMove{
				{C1, B1, QUIET},
				{C1, D2, QUIET},
			},
		},
		// TODO: add illegal moves
		{"rnbqk2r/pppp1ppp/4pn2/8/1Q1P4/5N2/PPP1PPPP/RNB1KB1R b KQkq - 2 4", BLACK, E8, 2,
			[]testMove{},
		},
		{"rnbqk2r/pppp2pp/4pp2/4N3/1Q1P4/4P2n/PPP1BPPP/RNB1K2R w KQkq - 0 8", WHITE, E1, 3,
			[]testMove{
				{E1, D1, QUIET},
				{E1, D2, QUIET},
				{E1, F1, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_KING_MOVES)
		b.generateKingMoves(&moves, tt.color)
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

func TestMakeMove(t *testing.T) {
	tests := []struct {
		startFen string
		move     testMove
		endFen   string
	}{
		{
			STARTING_FEN,
			testMove{E2, E4, DOUBLE_PAWN_PUSH},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			"rnbqkbnr/pp1ppppp/8/2p1P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
			testMove{F7, F5, DOUBLE_PAWN_PUSH},
			"rnbqkbnr/pp1pp1pp/8/2p1Pp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
		},
		{
			"rnbqkbnr/pppp1ppp/8/4p3/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2",
			testMove{A1, A3, QUIET},
			"rnbqkbnr/pppp1ppp/8/4p3/P7/R7/1PPPPPPP/1NBQKBNR b Kkq - 1 2",
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			testMove{B5, C6, CAPTURE},
			"r1bqk1nr/ppppbppp/2B5/4p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 0 4",
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			testMove{E1, G1, KING_CASTLE},
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
		},
		{
			"rnbqkbnr/pp2pppp/8/2pP4/8/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 3",
			testMove{D5, C6, EP_CAPTURE},
			"rnbqkbnr/pp2pppp/2P5/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, QUEEN_PROMOTION_CAPTURE},
			"Qnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, ROOK_PROMOTION_CAPTURE},
			"Rnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, BISHOP_PROMOTION_CAPTURE},
			"Bnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, KNIGHT_PROMOTION_CAPTURE},
			"Nnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/3P4/2N5/PPP3pP/R1BQKBNR b KQkq - 0 7",
			testMove{G2, H1, QUEEN_PROMOTION_CAPTURE},
			"rnbqkbnr/pP2pp1p/8/8/3P4/2N5/PPP4P/R1BQKBNq w Qkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, QUEEN_PROMOTION},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBqR w KQkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, ROOK_PROMOTION},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBrR w KQkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, BISHOP_PROMOTION},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBbR w KQkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, KNIGHT_PROMOTION},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBnR w KQkq - 0 8",
		},
	}
	seedKeys(time.Now().UTC().UnixNano())
	for _, tt := range tests {
		bStart, err := NewBoard(tt.startFen)
		if err != nil {
			t.Error(err)
		}
		m, err := newMoveFromTestMove(tt.move)
		if err != nil {
			t.Error(err)
		}
		err = bStart.makeMove(m)
		if err != nil {
			t.Error(err)
		}
		bEnd, err := NewBoard(tt.endFen)
		if err != nil {
			t.Error(err)
		}
		compareBoardState(bStart, bEnd, t)
	}
}

func TestUndoMove(t *testing.T) {
	tests := []struct {
		fen  string
		move testMove
	}{
		{
			STARTING_FEN,
			testMove{E2, E4, DOUBLE_PAWN_PUSH},
		},
		{
			"rnbqkbnr/pppp1ppp/8/4p3/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2",
			testMove{A1, A3, QUIET},
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			testMove{B5, C6, CAPTURE},
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			testMove{E1, G1, KING_CASTLE},
		},
		{
			"rnbqkbnr/pp2pppp/8/2pP4/8/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 3",
			testMove{D5, C6, EP_CAPTURE},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, QUEEN_PROMOTION_CAPTURE},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, ROOK_PROMOTION_CAPTURE},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, BISHOP_PROMOTION_CAPTURE},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			testMove{B7, A8, KNIGHT_PROMOTION_CAPTURE},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/3P4/2N5/PPP3pP/R1BQKBNR b KQkq - 0 7",
			testMove{G2, H1, QUEEN_PROMOTION_CAPTURE},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, QUEEN_PROMOTION},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, ROOK_PROMOTION},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, BISHOP_PROMOTION},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			testMove{G2, G1, KNIGHT_PROMOTION},
		},
	}
	seedKeys(time.Now().UTC().UnixNano())
	for _, tt := range tests {
		bStart, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		bEnd, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		m, err := newMoveFromTestMove(tt.move)
		if err != nil {
			t.Error(err)
		}
		err = bStart.makeMove(m)
		if err != nil {
			t.Error(err)
		}
		err = bStart.undoMove()
		if err != nil {
			t.Error(err)
		}
		compareBoardState(bStart, bEnd, t)
	}
}

func TestUndoMoves(t *testing.T) {
	tests := []struct {
		startFen  string
		endFen    string
		undoCount int
		moves     []testMove
	}{
		{
			STARTING_FEN,
			STARTING_FEN,
			2,
			[]testMove{
				{E2, E4, DOUBLE_PAWN_PUSH},
				{E7, E5, DOUBLE_PAWN_PUSH},
			},
		},
		{
			STARTING_FEN,
			"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			1,
			[]testMove{
				{E2, E4, DOUBLE_PAWN_PUSH},
				{E7, E5, DOUBLE_PAWN_PUSH},
				{G1, F3, QUIET},
				{B8, C6, QUIET},
			},
		},
		{
			STARTING_FEN,
			"r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
			2,
			[]testMove{
				{E2, E4, DOUBLE_PAWN_PUSH},
				{E7, E5, DOUBLE_PAWN_PUSH},
				{G1, F3, QUIET},
				{B8, C6, QUIET},
				{F1, B5, QUIET},
				{G8, F6, QUIET},
			},
		},
		{
			STARTING_FEN,
			"r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			2,
			[]testMove{
				{E2, E4, DOUBLE_PAWN_PUSH},
				{E7, E5, DOUBLE_PAWN_PUSH},
				{G1, F3, QUIET},
				{B8, C6, QUIET},
				{F1, B5, QUIET},
				{G8, F6, QUIET},
				{E1, G1, KING_CASTLE},
				{F8, E7, QUIET},
			},
		},
		{
			STARTING_FEN,
			"r1bqk2r/ppppbppp/2B2n2/4p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 0 5",
			1,
			[]testMove{
				{E2, E4, DOUBLE_PAWN_PUSH},
				{E7, E5, DOUBLE_PAWN_PUSH},
				{G1, F3, QUIET},
				{B8, C6, QUIET},
				{F1, B5, QUIET},
				{G8, F6, QUIET},
				{E1, G1, KING_CASTLE},
				{F8, E7, QUIET},
				{B5, C6, CAPTURE},
				{D7, C6, CAPTURE},
			},
		},
	}
	seedKeys(time.Now().UTC().UnixNano())
	for _, tt := range tests {
		bStart, err := NewBoard(tt.startFen)
		if err != nil {
			t.Error(err)
		}
		bEnd, err := NewBoard(tt.endFen)
		if err != nil {
			t.Error(err)
		}
		for _, tm := range tt.moves {
			legalMoves := bStart.generateMoves(bStart.sideToMove)
			m, err := newMoveFromTestMove(tm)
			if err != nil {
				t.Error(err)
			}
			isLegal := false
			for _, lm := range legalMoves {
				if lm == m {
					isLegal = true
				}
			}
			if !isLegal {
				t.Errorf("%v is not a valid move", tm)
			}
			err = bStart.makeMove(m)
			if err != nil {
				t.Error(err)
			}
		}
		for i := 0; i < tt.undoCount; i++ {
			err = bStart.undoMove()
			if err != nil {
				t.Error(err)
			}
		}
		compareBoardState(bStart, bEnd, t)
	}
}
