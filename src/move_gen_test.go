package main

import "testing"

func ContainsTestMove(t *testing.T, moves []Move, mu MoveUnpacked) {
	m := NewMoveFromMoveUnpacked(mu)
	var s string
	for _, move := range moves {
		if move == m {
			return
		}
		s += "\n" + move.ToString()
	}
	t.Errorf("\n\n%v\n\nnot found in\n%v\n\n", m.ToString(), s)
}

func FilterMovesByOrigin(c Coord, moves *[]Move) {
	n := 0
	for _, m := range *moves {
		var mu MoveUnpacked
		m.Unpack(&mu)
		if mu.originCoord == c {
			(*moves)[n] = m
			n++
		}
	}
	*moves = (*moves)[:n]
}

func TestGenerateMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, 20,
			[]MoveUnpacked{
				{A2, A3, WHITE_PAWN, 0, QUIET},
				{B2, B3, WHITE_PAWN, 0, QUIET},
				{C2, C3, WHITE_PAWN, 0, QUIET},
				{D2, D3, WHITE_PAWN, 0, QUIET},
				{E2, E3, WHITE_PAWN, 0, QUIET},
				{F2, F3, WHITE_PAWN, 0, QUIET},
				{G2, G3, WHITE_PAWN, 0, QUIET},
				{H2, H3, WHITE_PAWN, 0, QUIET},
				{A2, A4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{B2, B4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{C2, C4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{D2, D4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{F2, F4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{G2, G4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{H2, H4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{G1, F3, WHITE_KNIGHT, 0, QUIET},
				{G1, H3, WHITE_KNIGHT, 0, QUIET},
				{B1, A3, WHITE_KNIGHT, 0, QUIET},
				{B1, C3, WHITE_KNIGHT, 0, QUIET},
			},
		},
		{STARTING_FEN, BLACK, 20,
			[]MoveUnpacked{
				{A7, A6, BLACK_PAWN, 0, QUIET},
				{B7, B6, BLACK_PAWN, 0, QUIET},
				{C7, C6, BLACK_PAWN, 0, QUIET},
				{D7, D6, BLACK_PAWN, 0, QUIET},
				{E7, E6, BLACK_PAWN, 0, QUIET},
				{F7, F6, BLACK_PAWN, 0, QUIET},
				{G7, G6, BLACK_PAWN, 0, QUIET},
				{H7, H6, BLACK_PAWN, 0, QUIET},
				{A7, A5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{B7, B5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{C7, C5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{D7, D5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{F7, F5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{G7, G5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{H7, H5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{G8, F6, BLACK_KNIGHT, 0, QUIET},
				{G8, H6, BLACK_KNIGHT, 0, QUIET},
				{B8, A6, BLACK_KNIGHT, 0, QUIET},
				{B8, C6, BLACK_KNIGHT, 0, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		b.GenerateMoves(&moves, tt.color, false)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGenerateCaptureMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, 0,
			[]MoveUnpacked{},
		},
		{"r1bqkbnr/pp1ppppp/2n5/2p5/3PP3/5N2/PPP2PPP/RNBQKB1R b KQkq - 0 3", BLACK, 2,
			[]MoveUnpacked{
				{C5, D4, BLACK_PAWN, WHITE_PAWN, CAPTURE},
				{C6, D4, BLACK_KNIGHT, WHITE_PAWN, CAPTURE},
			},
		},
		{"r1bqkb1r/pp2pppp/2n2n2/2Pp4/4P3/5N2/PPP2PPP/RNBQKB1R w KQkq d6 0 5", WHITE, 3,
			[]MoveUnpacked{
				{E4, D5, WHITE_PAWN, BLACK_PAWN, CAPTURE},
				{C5, D6, WHITE_PAWN, EMPTY, EP_CAPTURE},
				{D1, D5, WHITE_QUEEN, BLACK_PAWN, CAPTURE},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		b.GenerateMoves(&moves, tt.color, true)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGeneratePawnMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		coord       Coord
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, C2, 2,
			[]MoveUnpacked{
				{C2, C4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH},
				{C2, C3, WHITE_PAWN, 0, QUIET},
			},
		},
		{"rnbqkbnr/pp2pppp/8/2ppP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3", WHITE, E5, 2,
			[]MoveUnpacked{
				{E5, E6, WHITE_PAWN, 0, QUIET},
				{E5, D6, WHITE_PAWN, 0, EP_CAPTURE},
			},
		},
		{"rnbqkbnr/pp2pppp/3P4/2p5/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3", BLACK, E7, 3,
			[]MoveUnpacked{
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH},
				{E7, E6, BLACK_PAWN, 0, QUIET},
				{E7, D6, BLACK_PAWN, WHITE_PAWN, CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P4/PPP2PPP/RNBQKBNR w KQkq - 0 3", WHITE, E4, 3,
			[]MoveUnpacked{
				{E4, E5, WHITE_PAWN, 0, QUIET},
				{E4, D5, WHITE_PAWN, BLACK_PAWN, CAPTURE},
				{E4, F5, WHITE_PAWN, BLACK_PAWN, CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P1P2/PPP3PP/RNBQKBNR b KQkq - 0 3", BLACK, D5, 2,
			[]MoveUnpacked{
				{D5, D4, BLACK_PAWN, 0, QUIET},
				{D5, E4, BLACK_PAWN, WHITE_PAWN, CAPTURE},
			},
		},
		{"rnbqkb1r/pppp2Pp/4pn2/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5", WHITE, G7, MAX_PAWN_MOVES,
			[]MoveUnpacked{
				{G7, F8, WHITE_PAWN, BLACK_BISHOP, KNIGHT_PROMOTION_CAPTURE},
				{G7, F8, WHITE_PAWN, BLACK_BISHOP, BISHOP_PROMOTION_CAPTURE},
				{G7, F8, WHITE_PAWN, BLACK_BISHOP, ROOK_PROMOTION_CAPTURE},
				{G7, F8, WHITE_PAWN, BLACK_BISHOP, QUEEN_PROMOTION_CAPTURE},
				{G7, G8, WHITE_PAWN, 0, KNIGHT_PROMOTION},
				{G7, G8, WHITE_PAWN, 0, BISHOP_PROMOTION},
				{G7, G8, WHITE_PAWN, 0, ROOK_PROMOTION},
				{G7, G8, WHITE_PAWN, 0, QUEEN_PROMOTION},
				{G7, H8, WHITE_PAWN, BLACK_ROOK, KNIGHT_PROMOTION_CAPTURE},
				{G7, H8, WHITE_PAWN, BLACK_ROOK, BISHOP_PROMOTION_CAPTURE},
				{G7, H8, WHITE_PAWN, BLACK_ROOK, ROOK_PROMOTION_CAPTURE},
				{G7, H8, WHITE_PAWN, BLACK_ROOK, QUEEN_PROMOTION_CAPTURE},
			},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/PpPP1PPP/RNBQ1RK1 b q - 1 9", BLACK, B2, 8,
			[]MoveUnpacked{
				{B2, A1, BLACK_PAWN, WHITE_ROOK, KNIGHT_PROMOTION_CAPTURE},
				{B2, A1, BLACK_PAWN, WHITE_ROOK, BISHOP_PROMOTION_CAPTURE},
				{B2, A1, BLACK_PAWN, WHITE_ROOK, ROOK_PROMOTION_CAPTURE},
				{B2, A1, BLACK_PAWN, WHITE_ROOK, QUEEN_PROMOTION_CAPTURE},
				{B2, C1, BLACK_PAWN, WHITE_BISHOP, KNIGHT_PROMOTION_CAPTURE},
				{B2, C1, BLACK_PAWN, WHITE_BISHOP, BISHOP_PROMOTION_CAPTURE},
				{B2, C1, BLACK_PAWN, WHITE_BISHOP, ROOK_PROMOTION_CAPTURE},
				{B2, C1, BLACK_PAWN, WHITE_BISHOP, QUEEN_PROMOTION_CAPTURE},
			},
		},
		{"rnbqkbnr/1ppppppp/p7/8/8/2N5/PPPPPPPP/R1BQKBNR w KQkq - 0 2", WHITE, C2, 0,
			[]MoveUnpacked{},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_PAWN_MOVES)
		b.GeneratePawnMoves(&moves, tt.color, false)
		FilterMovesByOrigin(tt.coord, &moves)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGenerateKnightMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		coord       Coord
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, G1, 2,
			[]MoveUnpacked{
				{G1, F3, WHITE_KNIGHT, 0, QUIET},
				{G1, H3, WHITE_KNIGHT, 0, QUIET},
			}},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/3N4/8/PPPPPPPP/RNBQKB1R w KQkq - 0 3", WHITE, D4, 6,
			[]MoveUnpacked{
				{D4, B3, WHITE_KNIGHT, 0, QUIET},
				{D4, B5, WHITE_KNIGHT, 0, QUIET},
				{D4, C6, WHITE_KNIGHT, 0, QUIET},
				{D4, E6, WHITE_KNIGHT, 0, QUIET},
				{D4, F5, WHITE_KNIGHT, BLACK_PAWN, CAPTURE},
				{D4, F3, WHITE_KNIGHT, 0, QUIET},
			},
		},
		{"r1bqkbnr/ppp1p1pp/8/3p4/3n1p2/2N1P1P1/PPPPQP1P/R1B1KB1R b KQkq - 1 6", BLACK, D4, MAX_KNIGHT_MOVES,
			[]MoveUnpacked{
				{D4, B3, BLACK_KNIGHT, 0, QUIET},
				{D4, B5, BLACK_KNIGHT, 0, QUIET},
				{D4, C6, BLACK_KNIGHT, 0, QUIET},
				{D4, E6, BLACK_KNIGHT, 0, QUIET},
				{D4, F5, BLACK_KNIGHT, 0, QUIET},
				{D4, F3, BLACK_KNIGHT, 0, QUIET},
				{D4, E2, BLACK_KNIGHT, WHITE_QUEEN, CAPTURE},
				{D4, C2, BLACK_KNIGHT, WHITE_PAWN, CAPTURE},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_KNIGHT_MOVES)
		b.GenerateKnightMoves(&moves, tt.color, false)
		FilterMovesByOrigin(tt.coord, &moves)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGenerateBishopMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		coord       Coord
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, F1, 0,
			[]MoveUnpacked{},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", WHITE, F1, 5,
			[]MoveUnpacked{
				{F1, E2, WHITE_BISHOP, 0, QUIET},
				{F1, D3, WHITE_BISHOP, 0, QUIET},
				{F1, C4, WHITE_BISHOP, 0, QUIET},
				{F1, B5, WHITE_BISHOP, 0, QUIET},
				{F1, A6, WHITE_BISHOP, 0, QUIET},
			},
		},
		{"rnbqkbnr/pp1pp1pp/5p2/2p5/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 3", WHITE, C4, 10,
			[]MoveUnpacked{
				{C4, B3, WHITE_BISHOP, 0, QUIET},
				{C4, B5, WHITE_BISHOP, 0, QUIET},
				{C4, A6, WHITE_BISHOP, 0, QUIET},
				{C4, D3, WHITE_BISHOP, 0, QUIET},
				{C4, E2, WHITE_BISHOP, 0, QUIET},
				{C4, F1, WHITE_BISHOP, 0, QUIET},
				{C4, D5, WHITE_BISHOP, 0, QUIET},
				{C4, E6, WHITE_BISHOP, 0, QUIET},
				{C4, F7, WHITE_BISHOP, 0, QUIET},
				{C4, G8, WHITE_BISHOP, BLACK_KNIGHT, CAPTURE},
			},
		},
		{"rn2kbnr/p2qp1p1/1p1p1p1p/1bpBP3/7P/P5P1/1PPP1P1R/RNBQK1N1 w Qkq - 3 9", WHITE, D5, MAX_BISHOP_MOVES,
			[]MoveUnpacked{
				{D5, C6, WHITE_BISHOP, 0, QUIET},
				{D5, B7, WHITE_BISHOP, 0, QUIET},
				{D5, A8, WHITE_BISHOP, BLACK_ROOK, CAPTURE},
				{D5, E6, WHITE_BISHOP, 0, QUIET},
				{D5, F7, WHITE_BISHOP, 0, QUIET},
				{D5, G8, WHITE_BISHOP, BLACK_KNIGHT, CAPTURE},
				{D5, C4, WHITE_BISHOP, 0, QUIET},
				{D5, B3, WHITE_BISHOP, 0, QUIET},
				{D5, A2, WHITE_BISHOP, 0, QUIET},
				{D5, E4, WHITE_BISHOP, 0, QUIET},
				{D5, F3, WHITE_BISHOP, 0, QUIET},
				{D5, G2, WHITE_BISHOP, 0, QUIET},
				{D5, H1, WHITE_BISHOP, 0, QUIET},
			}},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_BISHOP_MOVES)
		b.GenerateBishopMoves(&moves, tt.color, false)
		FilterMovesByOrigin(tt.coord, &moves)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGenerateRookMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		coord       Coord
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, A1, 0,
			[]MoveUnpacked{},
		},
		{"rnbqkbnr/2pppppp/8/8/8/R7/1PPPPPPP/1NBQKBNR w Kk - 2 5", WHITE, A3, MAX_ROOK_MOVES,
			[]MoveUnpacked{
				{A3, A1, WHITE_ROOK, 0, QUIET},
				{A3, A2, WHITE_ROOK, 0, QUIET},
				{A3, A4, WHITE_ROOK, 0, QUIET},
				{A3, A5, WHITE_ROOK, 0, QUIET},
				{A3, A6, WHITE_ROOK, 0, QUIET},
				{A3, A7, WHITE_ROOK, 0, QUIET},
				{A3, A8, WHITE_ROOK, BLACK_ROOK, CAPTURE},
				{A3, B3, WHITE_ROOK, 0, QUIET},
				{A3, C3, WHITE_ROOK, 0, QUIET},
				{A3, D3, WHITE_ROOK, 0, QUIET},
				{A3, E3, WHITE_ROOK, 0, QUIET},
				{A3, F3, WHITE_ROOK, 0, QUIET},
				{A3, G3, WHITE_ROOK, 0, QUIET},
				{A3, H3, WHITE_ROOK, 0, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_ROOK_MOVES)
		b.GenerateRookMoves(&moves, tt.color, false)
		FilterMovesByOrigin(tt.coord, &moves)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGenerateQueenMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		coord       Coord
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, D1, 0,
			[]MoveUnpacked{},
		},
		{"r2qk1r1/p1p3pp/np3p1n/3Q4/6bP/P5P1/1PP1PP1R/RNB1KBN1 w Qq - 1 10", WHITE, D5, MAX_QUEEN_MOVES,
			[]MoveUnpacked{
				{D5, C6, WHITE_QUEEN, 0, QUIET},
				{D5, B7, WHITE_QUEEN, 0, QUIET},
				{D5, A8, WHITE_QUEEN, BLACK_ROOK, CAPTURE},
				{D5, E6, WHITE_QUEEN, 0, QUIET},
				{D5, F7, WHITE_QUEEN, 0, QUIET},
				{D5, G8, WHITE_QUEEN, BLACK_ROOK, CAPTURE},
				{D5, C4, WHITE_QUEEN, 0, QUIET},
				{D5, B3, WHITE_QUEEN, 0, QUIET},
				{D5, A2, WHITE_QUEEN, 0, QUIET},
				{D5, E4, WHITE_QUEEN, 0, QUIET},
				{D5, F3, WHITE_QUEEN, 0, QUIET},
				{D5, G2, WHITE_QUEEN, 0, QUIET},
				{D5, H1, WHITE_QUEEN, 0, QUIET},
				{D5, D1, WHITE_QUEEN, 0, QUIET},
				{D5, D2, WHITE_QUEEN, 0, QUIET},
				{D5, D3, WHITE_QUEEN, 0, QUIET},
				{D5, D4, WHITE_QUEEN, 0, QUIET},
				{D5, D6, WHITE_QUEEN, 0, QUIET},
				{D5, D7, WHITE_QUEEN, 0, QUIET},
				{D5, D8, WHITE_QUEEN, BLACK_QUEEN, CAPTURE},
				{D5, A5, WHITE_QUEEN, 0, QUIET},
				{D5, B5, WHITE_QUEEN, 0, QUIET},
				{D5, C5, WHITE_QUEEN, 0, QUIET},
				{D5, E5, WHITE_QUEEN, 0, QUIET},
				{D5, F5, WHITE_QUEEN, 0, QUIET},
				{D5, G5, WHITE_QUEEN, 0, QUIET},
				{D5, H5, WHITE_QUEEN, 0, QUIET},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_QUEEN_MOVES)
		b.GenerateQueenMoves(&moves, tt.color, false)
		FilterMovesByOrigin(tt.coord, &moves)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}

func TestGenerateKingMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		coord       Coord
		movesLength int
		moves       []MoveUnpacked
	}{
		{STARTING_FEN, WHITE, E1, 0,
			[]MoveUnpacked{},
		},
		{"rn1q2nr/p2kb1pp/1p6/8/4K3/2N2p1N/PPPP4/R1BQ1B1R w - - 0 14", WHITE, E4, MAX_KING_MOVES,
			[]MoveUnpacked{
				{E4, D3, WHITE_KING, 0, QUIET},
				{E4, D4, WHITE_KING, 0, QUIET},
				{E4, D5, WHITE_KING, 0, QUIET},
				{E4, E3, WHITE_KING, 0, QUIET},
				{E4, E5, WHITE_KING, 0, QUIET},
				{E4, F3, WHITE_KING, BLACK_PAWN, CAPTURE},
				{E4, F4, WHITE_KING, 0, QUIET},
				{E4, F5, WHITE_KING, 0, QUIET},
			},
		},
		{"rn1q2nr/p2kb1pp/1p6/4p1p1/4K3/8/PPPP1P1P/RNBQ1BNR w - - 0 10", WHITE, E4, 8,
			[]MoveUnpacked{
				{E4, D3, WHITE_KING, 0, QUIET},
				{E4, D4, WHITE_KING, 0, QUIET},
				{E4, D5, WHITE_KING, 0, QUIET},
				{E4, E3, WHITE_KING, 0, QUIET},
				{E4, E5, WHITE_KING, BLACK_PAWN, CAPTURE},
				{E4, F3, WHITE_KING, 0, QUIET},
				{E4, F4, WHITE_KING, 0, QUIET},
				{E4, F5, WHITE_KING, 0, QUIET},
			},
		},
		{"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4", WHITE, E1, 3,
			[]MoveUnpacked{
				{E1, E2, WHITE_KING, 0, QUIET},
				{E1, F1, WHITE_KING, 0, QUIET},
				{E1, G1, WHITE_KING, 0, KING_CASTLE},
			},
		},
		{"r1bqk2r/ppppbppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 w kq - 6 5", WHITE, G1, 1,
			[]MoveUnpacked{
				{G1, H1, WHITE_KING, 0, QUIET},
			},
		},
		{"r1bqkb1r/ppp2ppp/2np1n2/8/4P3/2NQ4/PPP2PPP/R1B1KBNR w KQkq - 2 6", WHITE, E1, 3,
			[]MoveUnpacked{
				{E1, D1, WHITE_KING, 0, QUIET},
				{E1, D2, WHITE_KING, 0, QUIET},
				{E1, E2, WHITE_KING, 0, QUIET},
			},
		},
		{"r1bqk2r/ppp1bppp/2np1n2/6B1/4P3/2NQ4/PPP2PPP/R3KBNR w KQkq - 4 7", WHITE, E1, 4,
			[]MoveUnpacked{
				{E1, D1, WHITE_KING, 0, QUIET},
				{E1, D2, WHITE_KING, 0, QUIET},
				{E1, E2, WHITE_KING, 0, QUIET},
				{E1, C1, WHITE_KING, 0, QUEEN_CASTLE},
			},
		},
		{"r1bq1rk1/ppp1bppp/2np1n2/6B1/4P3/2NQ4/PPP2PPP/2KR1BNR w - - 6 8", WHITE, C1, 2,
			[]MoveUnpacked{
				{C1, B1, WHITE_KING, 0, QUIET},
				{C1, D2, WHITE_KING, 0, QUIET},
			},
		},
		{"rnbqk2r/pppp1ppp/4pn2/8/1Q1P4/5N2/PPP1PPPP/RNB1KB1R b KQkq - 2 4", BLACK, E8, 2,
			[]MoveUnpacked{
				{E8, F8, BLACK_KING, 0, QUIET},
				{E8, E7, BLACK_KING, 0, QUIET},
			},
		},
		{"rnbqk2r/pppp2pp/4pp2/4N3/1Q1P4/4P2n/PPP1BPPP/RNB1K2R w KQkq - 0 8", WHITE, E1, 3,
			[]MoveUnpacked{
				{E1, D1, WHITE_KING, 0, QUIET},
				{E1, D2, WHITE_KING, 0, QUIET},
				{E1, F1, WHITE_KING, 0, QUIET},
			},
		},
		{"r1bqk2r/2ppbppp/p1n2n2/1p2p3/4P3/1B3N2/PPPP1PPP/RNBQR1K1 b kq - 5 7", BLACK, E8, 2,
			[]MoveUnpacked{
				{E8, F8, BLACK_KING, 0, QUIET},
				{E8, G8, BLACK_KING, 0, KING_CASTLE},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := make([]Move, 0, MAX_KING_MOVES)
		b.GenerateKingMoves(&moves, tt.color, false)
		FilterMovesByOrigin(tt.coord, &moves)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, MoveUnpacked := range tt.moves {
			ContainsTestMove(t, moves, MoveUnpacked)
		}
	}
}
