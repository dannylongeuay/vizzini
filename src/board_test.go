package main

import (
	"testing"
)

type TestSquareChecks struct {
	file   File
	rank   Rank
	square Square
}

type TestPieceChecks struct {
	square Square
	count  int
	coords []Coord
}

func TestNewBoard(t *testing.T) {
	tests := []struct {
		fen            string
		whiteKingCoord Coord
		blackKingCoord Coord
		sideToMove     Color
		castleRights   CastleRights
		epCoord        Coord
		halfMove       HalfMove
		fullMove       int
		squareChecks   []TestSquareChecks
		pieceChecks    []TestPieceChecks
	}{
		{STARTING_FEN, E1, E8, WHITE, 15, 0, 0, 1,
			[]TestSquareChecks{
				{FILE_A, RANK_EIGHT, BLACK_ROOK},
				{FILE_B, RANK_EIGHT, BLACK_KNIGHT},
				{FILE_C, RANK_EIGHT, BLACK_BISHOP},
				{FILE_D, RANK_EIGHT, BLACK_QUEEN},
				{FILE_E, RANK_EIGHT, BLACK_KING},
				{FILE_F, RANK_EIGHT, BLACK_BISHOP},
				{FILE_G, RANK_EIGHT, BLACK_KNIGHT},
				{FILE_H, RANK_EIGHT, BLACK_ROOK},
				{FILE_A, RANK_SEVEN, BLACK_PAWN},
				{FILE_B, RANK_SEVEN, BLACK_PAWN},
				{FILE_C, RANK_SEVEN, BLACK_PAWN},
				{FILE_D, RANK_SEVEN, BLACK_PAWN},
				{FILE_E, RANK_SEVEN, BLACK_PAWN},
				{FILE_F, RANK_SEVEN, BLACK_PAWN},
				{FILE_G, RANK_SEVEN, BLACK_PAWN},
				{FILE_H, RANK_SEVEN, BLACK_PAWN},
				{FILE_A, RANK_SIX, EMPTY},
				{FILE_H, RANK_THREE, EMPTY},
				{FILE_A, RANK_TWO, WHITE_PAWN},
				{FILE_B, RANK_TWO, WHITE_PAWN},
				{FILE_C, RANK_TWO, WHITE_PAWN},
				{FILE_D, RANK_TWO, WHITE_PAWN},
				{FILE_E, RANK_TWO, WHITE_PAWN},
				{FILE_F, RANK_TWO, WHITE_PAWN},
				{FILE_G, RANK_TWO, WHITE_PAWN},
				{FILE_H, RANK_TWO, WHITE_PAWN},
				{FILE_A, RANK_ONE, WHITE_ROOK},
				{FILE_B, RANK_ONE, WHITE_KNIGHT},
				{FILE_C, RANK_ONE, WHITE_BISHOP},
				{FILE_D, RANK_ONE, WHITE_QUEEN},
				{FILE_E, RANK_ONE, WHITE_KING},
				{FILE_F, RANK_ONE, WHITE_BISHOP},
				{FILE_G, RANK_ONE, WHITE_KNIGHT},
				{FILE_H, RANK_ONE, WHITE_ROOK},
			},
			[]TestPieceChecks{
				{WHITE_PAWN, 8, []Coord{
					A2, B2, C2, D2, E2, F2, G2, H2,
				}},
				{WHITE_KNIGHT, 2, []Coord{
					B1, G1,
				}},
				{WHITE_BISHOP, 2, []Coord{
					C1, F1,
				}},
				{WHITE_ROOK, 2, []Coord{
					A1, H1,
				}},
				{WHITE_QUEEN, 1, []Coord{
					D1,
				}},
				{WHITE_KING, 1, []Coord{
					E1,
				}},
				{BLACK_PAWN, 8, []Coord{
					A7, B7, C7, D7, E7, F7, G7, H7,
				}},
				{BLACK_KNIGHT, 2, []Coord{
					B8, G8,
				}},
				{BLACK_BISHOP, 2, []Coord{
					C8, F8,
				}},
				{BLACK_ROOK, 2, []Coord{
					A8, H8,
				}},
				{BLACK_QUEEN, 1, []Coord{
					D8,
				}},
				{BLACK_KING, 1, []Coord{
					E8,
				}},
			},
		},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			E1, E8, BLACK, 15, E3, 0, 1,
			[]TestSquareChecks{
				{FILE_E, RANK_FOUR, WHITE_PAWN},
			},
			[]TestPieceChecks{
				{WHITE_PAWN, 8, []Coord{
					E4,
				}},
			},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			E1, E8, WHITE, 15, C6, 0, 2,
			[]TestSquareChecks{
				{FILE_C, RANK_FIVE, BLACK_PAWN},
			},
			[]TestPieceChecks{
				{BLACK_PAWN, 8, []Coord{
					C5,
				}},
			},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			E1, E8, BLACK, 15, 0, 1, 2,
			[]TestSquareChecks{
				{FILE_F, RANK_THREE, WHITE_KNIGHT},
			},
			[]TestPieceChecks{
				{WHITE_PAWN, 8, []Coord{
					E4,
				}},
				{WHITE_KNIGHT, 2, []Coord{
					F3,
				}},
			},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
			E1, E8, WHITE, 0, 0, 0, 1, []TestSquareChecks{}, []TestPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w q - 0 1",
			E1, E8, WHITE, 1, 0, 0, 1, []TestSquareChecks{}, []TestPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1",
			E1, E8, WHITE, 2, 0, 0, 1, []TestSquareChecks{}, []TestPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Q - 0 1",
			E1, E8, WHITE, 4, 0, 0, 1, []TestSquareChecks{}, []TestPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w K - 0 1",
			E1, E8, WHITE, 8, 0, 0, 1, []TestSquareChecks{}, []TestPieceChecks{},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/PpPP1PPP/RNBQ1RK1 b q - 1 9",
			G1, E8, BLACK, 1, 0, 1, 9,
			[]TestSquareChecks{
				{FILE_B, RANK_TWO, BLACK_PAWN},
				{FILE_H, RANK_EIGHT, WHITE_QUEEN},
			},
			[]TestPieceChecks{
				{WHITE_QUEEN, 2, []Coord{
					H8,
				}},
				{WHITE_BISHOP, 2, []Coord{
					D3,
				}},
				{WHITE_KNIGHT, 2, []Coord{
					F3,
				}},
				{BLACK_PAWN, 6, []Coord{
					E6, B2,
				}},
				{BLACK_KNIGHT, 2, []Coord{
					F6,
				}},
			},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/P1PP1PPP/qNBQ1RK1 w q - 0 10",
			G1, E8, WHITE, 1, 0, 0, 10,
			[]TestSquareChecks{
				{FILE_A, RANK_ONE, BLACK_QUEEN},
				{FILE_H, RANK_EIGHT, WHITE_QUEEN},
			},
			[]TestPieceChecks{
				{WHITE_QUEEN, 2, []Coord{
					H8,
				}},
				{WHITE_BISHOP, 2, []Coord{
					D3,
				}},
				{WHITE_KNIGHT, 2, []Coord{
					F3,
				}},
				{BLACK_PAWN, 5, []Coord{
					E6,
				}},
				{BLACK_QUEEN, 2, []Coord{
					A1,
				}},
				{BLACK_KNIGHT, 2, []Coord{
					F6,
				}},
			},
		},
	}
	for _, tt := range tests {
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.kingCoords[WHITE] != tt.whiteKingCoord {
			t.Errorf("white king index: %v != %v", b.kingCoords[WHITE], tt.whiteKingCoord)
		}
		if err != nil {
			t.Error(err)
		}
		if b.squares[tt.whiteKingCoord] != WHITE_KING {
			t.Errorf("white king square: %v != %v", b.squares[tt.whiteKingCoord], WHITE_KING)
		}
		if b.kingCoords[BLACK] != tt.blackKingCoord {
			t.Errorf("black king index: %v != %v", b.kingCoords[BLACK], tt.blackKingCoord)
		}
		if err != nil {
			t.Error(err)
		}
		if b.squares[tt.blackKingCoord] != BLACK_KING {
			t.Errorf("black king square: %v != %v", b.squares[tt.blackKingCoord], BLACK_KING)
		}
		if b.sideToMove != tt.sideToMove {
			t.Errorf("side: %v != %v", b.sideToMove, tt.sideToMove)
		}

		if b.castleRights != tt.castleRights {
			t.Errorf("castling rights: %v != %v", b.castleRights, tt.castleRights)
		}

		if b.epCoord != tt.epCoord {
			t.Errorf("en passant index: %v != %v", b.epCoord, tt.epCoord)
		}

		if b.halfMove != tt.halfMove {
			t.Errorf("half move clock: %v != %v", b.halfMove, tt.halfMove)
		}

		if b.fullMove != tt.fullMove {
			t.Errorf("full move clock: %v != %v", b.fullMove, tt.fullMove)
		}
	}
}

func TestBoardToString(t *testing.T) {
	expected := `
_________________________
|♖ |♘ |♗ |♕ |♔ |♗ |♘ |♖ | 8
_________________________
|♙ |♙ |♙ |♙ |♙ |♙ |♙ |♙ | 7
_________________________
|  |  |  |  |  |  |  |  | 6
_________________________
|  |  |  |  |  |  |  |  | 5
_________________________
|  |  |  |  |  |  |  |  | 4
_________________________
|  |  |  |  |  |  |  |  | 3
_________________________
|♟ |♟ |♟ |♟ |♟ |♟ |♟ |♟ | 2
_________________________
|♜ |♞ |♝ |♛ |♚ |♝ |♞ |♜ | 1
_________________________
 A  B  C  D  E  F  G  H`

	board, err := NewBoard(STARTING_FEN)
	if err != nil {
		t.Error(err)
	}
	if board.ToString() != expected {
		t.Errorf("board to string: %v \n\n!=\n %v", board.ToString(), expected)
	}
}
