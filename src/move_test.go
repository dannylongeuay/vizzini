package main

import (
	"testing"
)

func CompareBoardState(bStart *Board, bEnd *Board, t *testing.T) {
	if bStart.kingCoords[WHITE] != bEnd.kingCoords[WHITE] {
		t.Errorf("board white king index: %v != %v", bStart.kingCoords[WHITE], bEnd.kingCoords[WHITE])
	}
	if bStart.kingCoords[BLACK] != bEnd.kingCoords[BLACK] {
		t.Errorf("board black king index: %v != %v", bStart.kingCoords[BLACK], bEnd.kingCoords[BLACK])
	}
	if bStart.sideToMove != bEnd.sideToMove {
		t.Errorf("board side to move: %v != %v", bStart.sideToMove, bEnd.sideToMove)
	}
	if bStart.castleRights != bEnd.castleRights {
		t.Errorf("board castle rights: %v != %v", bStart.castleRights, bEnd.castleRights)
	}
	if bStart.epCoord != bEnd.epCoord {
		t.Errorf("board en passant: %v != %v", COORD_STRINGS[bStart.epCoord], COORD_STRINGS[bEnd.epCoord])
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
	if !IsBitboardEqual(t, bStart.bbWP, bEnd.bbWP) {
		t.Error("white pawn bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbWN, bEnd.bbWN) {
		t.Error("white knight bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbWB, bEnd.bbWB) {
		t.Error("white bishop bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbWR, bEnd.bbWR) {
		t.Error("white rook bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbWQ, bEnd.bbWQ) {
		t.Error("white queen bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbWK, bEnd.bbWK) {
		t.Error("white king bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBP, bEnd.bbBP) {
		t.Error("black pawn bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBN, bEnd.bbBN) {
		t.Error("black knight bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBB, bEnd.bbBB) {
		t.Error("black bishop bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBR, bEnd.bbBR) {
		t.Error("black rook bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBQ, bEnd.bbBQ) {
		t.Error("black queen bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBK, bEnd.bbBK) {
		t.Error("black king bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbWhitePieces, bEnd.bbWhitePieces) {
		t.Error("white piece bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbBlackPieces, bEnd.bbBlackPieces) {
		t.Error("black piece bitboards are not equal")
	}
	if !IsBitboardEqual(t, bStart.bbAllPieces, bEnd.bbAllPieces) {
		t.Error("all piece bitboards are not equal")
	}
}

func TestMakeMove(t *testing.T) {
	tests := []struct {
		startFen string
		mu       MoveUnpacked
		endFen   string
	}{
		{
			STARTING_FEN,
			MoveUnpacked{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
		},
		{
			"rnbqkbnr/pp1ppppp/8/2p1P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
			MoveUnpacked{F7, F5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
			"rnbqkbnr/pp1pp1pp/8/2p1Pp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
		},
		{
			"rnbqkbnr/pppppppp/8/8/2P5/8/PP1PPPPP/RNBQKBNR b KQkq - 0 1",
			MoveUnpacked{B7, B5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
			"rnbqkbnr/p1pppppp/8/1p6/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 2",
		},
		{
			"rnbqkbnr/pppp1ppp/8/4p3/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2",
			MoveUnpacked{A1, A3, WHITE_ROOK, 0, QUIET, 0},
			"rnbqkbnr/pppp1ppp/8/4p3/P7/R7/1PPPPPPP/1NBQKBNR b Kkq - 1 2",
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			MoveUnpacked{B5, C6, WHITE_BISHOP, BLACK_KNIGHT, CAPTURE, 0},
			"r1bqk1nr/ppppbppp/2B5/4p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 0 4",
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			MoveUnpacked{E1, G1, WHITE_KING, 0, KING_CASTLE, 0},
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
		},
		{
			"rnbqkbnr/pp2pppp/8/2pP4/8/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 3",
			MoveUnpacked{D5, C6, WHITE_PAWN, 0, EP_CAPTURE, 0},
			"rnbqkbnr/pp2pppp/2P5/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, QUEEN_PROMOTION_CAPTURE, 0},
			"Qnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, ROOK_PROMOTION_CAPTURE, 0},
			"Rnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, BISHOP_PROMOTION_CAPTURE, 0},
			"Bnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, KNIGHT_PROMOTION_CAPTURE, 0},
			"Nnbqkbnr/p3pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR b KQk - 0 5",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/3P4/2N5/PPP3pP/R1BQKBNR b KQkq - 0 7",
			MoveUnpacked{G2, H1, BLACK_PAWN, WHITE_ROOK, QUEEN_PROMOTION_CAPTURE, 0},
			"rnbqkbnr/pP2pp1p/8/8/3P4/2N5/PPP4P/R1BQKBNq w Qkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, QUEEN_PROMOTION, 0},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBqR w KQkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, ROOK_PROMOTION, 0},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBrR w KQkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, BISHOP_PROMOTION, 0},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBbR w KQkq - 0 8",
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, KNIGHT_PROMOTION, 0},
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN2P/R1BQKBnR w KQkq - 0 8",
		},
	}
	for _, tt := range tests {
		bStart, err := NewBoard(tt.startFen)
		if err != nil {
			t.Error(err)
		}
		m := NewMoveFromMoveUnpacked(tt.mu)
		err = bStart.MakeMove(m)
		if err != nil {
			t.Error(err)
		}
		bEnd, err := NewBoard(tt.endFen)
		if err != nil {
			t.Error(err)
		}
		CompareBoardState(bStart, bEnd, t)
	}
}

func TestUndoMove(t *testing.T) {
	tests := []struct {
		fen string
		mu  MoveUnpacked
	}{
		{
			STARTING_FEN,
			MoveUnpacked{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
		},
		{
			"rnbqkbnr/pppp1ppp/8/4p3/P7/8/1PPPPPPP/RNBQKBNR w KQkq - 0 2",
			MoveUnpacked{A1, A3, WHITE_ROOK, 0, QUIET, 0},
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			MoveUnpacked{B5, C6, WHITE_BISHOP, BLACK_KNIGHT, CAPTURE, 0},
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			MoveUnpacked{E1, G1, WHITE_KING, 0, KING_CASTLE, 0},
		},
		{
			"rnbqkbnr/pp2pppp/8/2pP4/8/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 3",
			MoveUnpacked{D5, C6, WHITE_PAWN, 0, EP_CAPTURE, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, QUEEN_PROMOTION_CAPTURE, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, ROOK_PROMOTION_CAPTURE, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, BISHOP_PROMOTION_CAPTURE, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/6p1/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5",
			MoveUnpacked{B7, A8, WHITE_PAWN, BLACK_ROOK, KNIGHT_PROMOTION_CAPTURE, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/3P4/2N5/PPP3pP/R1BQKBNR b KQkq - 0 7",
			MoveUnpacked{G2, H1, BLACK_PAWN, WHITE_ROOK, QUEEN_PROMOTION_CAPTURE, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, QUEEN_PROMOTION, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, ROOK_PROMOTION, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, BISHOP_PROMOTION, 0},
		},
		{
			"rnbqkbnr/pP2pp1p/8/8/8/2N5/PPPPN1pP/R1BQKB1R b KQkq - 1 7",
			MoveUnpacked{G2, G1, BLACK_PAWN, 0, KNIGHT_PROMOTION, 0},
		},
	}
	for _, tt := range tests {
		bStart, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		bEnd, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		m := NewMoveFromMoveUnpacked(tt.mu)
		err = bStart.MakeMove(m)
		if err != nil {
			t.Error(err)
		}
		err = bStart.UndoMove()
		if err != nil {
			t.Error(err)
		}
		CompareBoardState(bStart, bEnd, t)
	}
}

func TestUndoMoves(t *testing.T) {
	tests := []struct {
		startFen  string
		endFen    string
		undoCount int
		mus       []MoveUnpacked
	}{
		{
			STARTING_FEN,
			STARTING_FEN,
			2,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
			},
		},
		{
			STARTING_FEN,
			"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			1,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{G1, F3, WHITE_KNIGHT, 0, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, 0, QUIET, 0},
			},
		},
		{
			STARTING_FEN,
			"r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
			2,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{G1, F3, WHITE_KNIGHT, 0, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, 0, QUIET, 0},
				{F1, B5, WHITE_BISHOP, 0, QUIET, 0},
				{G8, F6, BLACK_KNIGHT, 0, QUIET, 0},
			},
		},
		{
			STARTING_FEN,
			"r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			2,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{G1, F3, WHITE_KNIGHT, 0, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, 0, QUIET, 0},
				{F1, B5, WHITE_BISHOP, 0, QUIET, 0},
				{G8, F6, BLACK_KNIGHT, 0, QUIET, 0},
				{E1, G1, WHITE_KING, 0, KING_CASTLE, 0},
				{F8, E7, BLACK_BISHOP, 0, QUIET, 0},
			},
		},
		{
			STARTING_FEN,
			"r1bqk2r/ppppbppp/2B2n2/4p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 0 5",
			1,
			[]MoveUnpacked{
				{E2, E4, WHITE_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{E7, E5, BLACK_PAWN, 0, DOUBLE_PAWN_PUSH, 0},
				{G1, F3, WHITE_KNIGHT, 0, QUIET, 0},
				{B8, C6, BLACK_KNIGHT, 0, QUIET, 0},
				{F1, B5, WHITE_BISHOP, 0, QUIET, 0},
				{G8, F6, BLACK_KNIGHT, 0, QUIET, 0},
				{E1, G1, WHITE_KING, 0, KING_CASTLE, 0},
				{F8, E7, BLACK_BISHOP, 0, QUIET, 0},
				{B5, C6, WHITE_BISHOP, BLACK_KNIGHT, CAPTURE, 23},
				{D7, C6, BLACK_PAWN, WHITE_BISHOP, CAPTURE, 35},
			},
		},
	}
	for _, tt := range tests {
		bStart, err := NewBoard(tt.startFen)
		if err != nil {
			t.Error(err)
		}
		bEnd, err := NewBoard(tt.endFen)
		if err != nil {
			t.Error(err)
		}
		for _, tm := range tt.mus {
			moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
			bStart.GenerateMoves(&moves, bStart.sideToMove)
			move := NewMoveFromMoveUnpacked(tm)
			isLegal := false
			for _, m := range moves {
				if m == move {
					isLegal = true
				}
			}
			if !isLegal {
				t.Errorf("%v is not a valid move", tm)
			}
			err = bStart.MakeMove(move)
			if err != nil {
				t.Error(err)
			}
		}
		for i := 0; i < tt.undoCount; i++ {
			err = bStart.UndoMove()
			if err != nil {
				t.Error(err)
			}
		}
		CompareBoardState(bStart, bEnd, t)
	}
}
