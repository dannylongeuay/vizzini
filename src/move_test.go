package main

import (
	"testing"
)

func containsMove(moves []move, testMove move) bool {
	for _, move := range moves {
		if move == testMove {
			return true
		}
	}
	return false
}

func TestGenerateMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		movesLength int
		moves       []move
	}{
		{STARTING_FEN, WHITE, 16,
			[]move{
				{81, 71, QUIET},
				{82, 72, QUIET},
				{83, 73, QUIET},
				{84, 74, QUIET},
				{85, 75, QUIET},
				{86, 76, QUIET},
				{87, 77, QUIET},
				{88, 78, QUIET},
				{81, 61, DOUBLE_PAWN_PUSH},
				{82, 62, DOUBLE_PAWN_PUSH},
				{83, 63, DOUBLE_PAWN_PUSH},
				{84, 64, DOUBLE_PAWN_PUSH},
				{85, 65, DOUBLE_PAWN_PUSH},
				{86, 66, DOUBLE_PAWN_PUSH},
				{87, 67, DOUBLE_PAWN_PUSH},
				{88, 68, DOUBLE_PAWN_PUSH},
			},
		},
		{STARTING_FEN, BLACK, 16,
			[]move{
				{31, 41, QUIET},
				{32, 42, QUIET},
				{33, 43, QUIET},
				{34, 44, QUIET},
				{35, 45, QUIET},
				{36, 46, QUIET},
				{37, 47, QUIET},
				{38, 48, QUIET},
				{31, 51, DOUBLE_PAWN_PUSH},
				{32, 52, DOUBLE_PAWN_PUSH},
				{33, 53, DOUBLE_PAWN_PUSH},
				{34, 54, DOUBLE_PAWN_PUSH},
				{35, 55, DOUBLE_PAWN_PUSH},
				{36, 56, DOUBLE_PAWN_PUSH},
				{37, 57, DOUBLE_PAWN_PUSH},
				{38, 58, DOUBLE_PAWN_PUSH},
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
			if !containsMove(moves, testMove) {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}

func TestGeneratePawnMoves(t *testing.T) {
	tests := []struct {
		fen         string
		color       Color
		squareIndex SquareIndex
		movesLength int
		moves       []move
	}{
		{STARTING_FEN, WHITE, 83, 2,
			[]move{
				{83, 63, DOUBLE_PAWN_PUSH},
				{83, 73, QUIET},
			},
		},
		{"rnbqkbnr/pp2pppp/8/2ppP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3", WHITE, 55, 2,
			[]move{
				{55, 45, QUIET},
				{55, 44, EP_CAPTURE},
			},
		},
		{"rnbqkbnr/pp2pppp/3P4/2p5/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3", BLACK, 35, 3,
			[]move{
				{35, 55, DOUBLE_PAWN_PUSH},
				{35, 45, QUIET},
				{35, 44, CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P4/PPP2PPP/RNBQKBNR w KQkq - 0 3", WHITE, 65, 3,
			[]move{
				{65, 55, QUIET},
				{65, 54, CAPTURE},
				{65, 56, CAPTURE},
			},
		},
		{"rnbqkbnr/ppp1p1pp/8/3p1p2/4P3/3P1P2/PPP3PP/RNBQKBNR b KQkq - 0 3", BLACK, 54, 2,
			[]move{
				{54, 64, QUIET},
				{54, 65, CAPTURE},
			},
		},
		{"rnbqkb1r/pppp2Pp/4pn2/8/8/8/PPPP1PPP/RNBQKBNR w KQkq - 0 5", WHITE, 37, 12,
			[]move{
				{37, 26, KNIGHT_PROMOTION_CAPTURE},
				{37, 26, BISHOP_PROMOTION_CAPTURE},
				{37, 26, ROOK_PROMOTION_CAPTURE},
				{37, 26, QUEEN_PROMOTION_CAPTURE},
				{37, 27, KNIGHT_PROMOTION},
				{37, 27, BISHOP_PROMOTION},
				{37, 27, ROOK_PROMOTION},
				{37, 27, QUEEN_PROMOTION},
				{37, 28, KNIGHT_PROMOTION_CAPTURE},
				{37, 28, BISHOP_PROMOTION_CAPTURE},
				{37, 28, ROOK_PROMOTION_CAPTURE},
				{37, 28, QUEEN_PROMOTION_CAPTURE},
			},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/PpPP1PPP/RNBQ1RK1 b q - 1 9", BLACK, 82, 8,
			[]move{
				{82, 91, KNIGHT_PROMOTION_CAPTURE},
				{82, 91, BISHOP_PROMOTION_CAPTURE},
				{82, 91, ROOK_PROMOTION_CAPTURE},
				{82, 91, QUEEN_PROMOTION_CAPTURE},
				{82, 93, KNIGHT_PROMOTION_CAPTURE},
				{82, 93, BISHOP_PROMOTION_CAPTURE},
				{82, 93, ROOK_PROMOTION_CAPTURE},
				{82, 93, QUEEN_PROMOTION_CAPTURE},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		moves := b.generatePawnMoves(tt.color, tt.squareIndex)
		if len(moves) != tt.movesLength {
			t.Errorf("moves length: %v != %v", len(moves), tt.movesLength)
		}
		for _, testMove := range tt.moves {
			if !containsMove(moves, testMove) {
				t.Errorf("unable to find move %v in %v", testMove, moves)
			}
		}
	}
}
