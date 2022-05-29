package main

import (
	"strings"
	"testing"
	"time"
)

func TestUCIParseMove(t *testing.T) {
	tests := []struct {
		fen          string
		moveNotation string
		mu           MoveUnpacked
	}{
		{
			STARTING_FEN,
			"e2e4",
			MoveUnpacked{
				E2,
				E4,
				WHITE_PAWN,
				EMPTY,
				DOUBLE_PAWN_PUSH,
			},
		},
		{
			"rnbqkbnr/pp1pp1pp/8/2p1Pp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3",
			"e5f6",
			MoveUnpacked{
				E5,
				F6,
				WHITE_PAWN,
				EMPTY,
				EP_CAPTURE,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7g8n",
			MoveUnpacked{
				G7,
				G8,
				WHITE_PAWN,
				EMPTY,
				KNIGHT_PROMOTION,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7g8b",
			MoveUnpacked{
				G7,
				G8,
				WHITE_PAWN,
				EMPTY,
				BISHOP_PROMOTION,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7g8r",
			MoveUnpacked{
				G7,
				G8,
				WHITE_PAWN,
				EMPTY,
				ROOK_PROMOTION,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7g8q",
			MoveUnpacked{
				G7,
				G8,
				WHITE_PAWN,
				EMPTY,
				QUEEN_PROMOTION,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7h8n",
			MoveUnpacked{
				G7,
				H8,
				WHITE_PAWN,
				BLACK_ROOK,
				KNIGHT_PROMOTION_CAPTURE,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7h8b",
			MoveUnpacked{
				G7,
				H8,
				WHITE_PAWN,
				BLACK_ROOK,
				BISHOP_PROMOTION_CAPTURE,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7h8r",
			MoveUnpacked{
				G7,
				H8,
				WHITE_PAWN,
				BLACK_ROOK,
				ROOK_PROMOTION_CAPTURE,
			},
		},
		{
			"r1bqkb1r/pp1pp1Pp/2n2n2/2p5/8/8/PPPP1PPP/RNBQKBNR w KQkq - 1 5",
			"g7h8q",
			MoveUnpacked{
				G7,
				H8,
				WHITE_PAWN,
				BLACK_ROOK,
				QUEEN_PROMOTION_CAPTURE,
			},
		},
		{
			"r3kb1Q/pp1bp2p/1qnp1n2/2p5/2B5/5N2/PPPP1PPP/RNBQK2R w KQq - 2 8",
			"e1g1",
			MoveUnpacked{
				E1,
				G1,
				WHITE_KING,
				EMPTY,
				KING_CASTLE,
			},
		},
		{
			"r3kb1Q/pp1bp2p/1qnp1n2/2p5/2B5/5N2/PPPP1PPP/RNBQ1RK1 b q - 3 8",
			"e8c8",
			MoveUnpacked{
				E8,
				C8,
				BLACK_KING,
				EMPTY,
				QUEEN_CASTLE,
			},
		},
		{
			"2kr1b1Q/pp1bp2p/1qnp1n2/2pB4/8/5N2/PPPP1PPP/RNBQ1RK1 b - - 5 9",
			"f6d5",
			MoveUnpacked{
				F6,
				D5,
				BLACK_KNIGHT,
				WHITE_BISHOP,
				CAPTURE,
			},
		},
		{
			"2kr1b1Q/pp1bp2p/1qnp4/2pn4/8/5N2/PPPP1PPP/RNBQ1RK1 w - - 0 10",
			"f3g5",
			MoveUnpacked{
				F3,
				G5,
				WHITE_KNIGHT,
				EMPTY,
				QUIET,
			},
		},
	}
	for _, tt := range tests {
		board, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actual, err := board.ParseUCIMove(tt.moveNotation)
		if err != nil {
			t.Error(err)
		}
		expected := NewMove(tt.mu.originCoord, tt.mu.dstCoord, tt.mu.originSquare, tt.mu.dstSquare, tt.mu.moveKind)

		if actual != expected {
			t.Errorf("move notation: %v != %v", actual.ToString(), expected.ToString())
		}

	}

}

func TestUCISetGoParams(t *testing.T) {
	tests := []struct {
		s        string
		duration time.Duration
	}{
		{
			"infinite",
			time.Duration(time.Minute - time.Second),
		},
		{
			"wtime 10000 btime 8000 winc 1000 binc 1000 movestogo 40 depth 4 nodes 100000 movetime 5000",
			time.Duration(time.Millisecond * time.Duration(4900)),
		},
		{
			"wtime 10000 btime 8000 winc 1000 binc 1000 movestogo 40 depth 4 nodes 100000",
			time.Duration(time.Millisecond * time.Duration(1100)),
		},
	}
	for _, tt := range tests {
		args := strings.Split(tt.s, " ")
		uci, err := NewUCI()
		if err != nil {
			t.Error(err)
		}
		uci.SetGoParams(args)
		approxStopTime := time.Now().Add(tt.duration)
		if !uci.stopTime.After(approxStopTime) {
			t.Errorf("stopTime %v is not after %v", uci.stopTime, approxStopTime)
		}
	}
}
