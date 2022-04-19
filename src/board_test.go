package main

import (
	"testing"
)

type testSquareChecks struct {
	file   File
	rank   Rank
	square Square
}

type testPieceChecks struct {
	square      Square
	piecesCount int
	pieceCoords []SquareCoord
}

func containsPieceIndex(indexes []SquareIndex, testIndex SquareIndex) bool {
	for _, index := range indexes {
		if index == testIndex {
			return true
		}
	}
	return false
}

func TestSquareByIndexes64(t *testing.T) {
	tests := []struct {
		index    int
		expected SquareIndex
	}{
		{0, 21},
		{63, 98},
	}
	for _, tt := range tests {

		actual := SquareIndexes64[tt.index]
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
}

func TestNewBoard(t *testing.T) {
	tests := []struct {
		fen          string
		sideToMove   Color
		castleRights CastleRights
		epIndex      SquareIndex
		halfMove     int
		fullMove     int
		squareChecks []testSquareChecks
		pieceChecks  []testPieceChecks
	}{
		{STARTING_FEN, WHITE, 15, 0, 0, 1,
			[]testSquareChecks{
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
				{FILE_NONE, RANK_NONE, INVALID},
			},
			[]testPieceChecks{
				{WHITE_PAWN, 8, []SquareCoord{
					"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
				}},
				{WHITE_KNIGHT, 2, []SquareCoord{
					"b1", "g1",
				}},
				{WHITE_BISHOP, 2, []SquareCoord{
					"c1", "f1",
				}},
				{WHITE_ROOK, 2, []SquareCoord{
					"a1", "h1",
				}},
				{WHITE_QUEEN, 1, []SquareCoord{
					"d1",
				}},
				{WHITE_KING, 1, []SquareCoord{
					"e1",
				}},
				{BLACK_PAWN, 8, []SquareCoord{
					"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
				}},
				{BLACK_KNIGHT, 2, []SquareCoord{
					"b8", "g8",
				}},
				{BLACK_BISHOP, 2, []SquareCoord{
					"c8", "f8",
				}},
				{BLACK_ROOK, 2, []SquareCoord{
					"a8", "h8",
				}},
				{BLACK_QUEEN, 1, []SquareCoord{
					"d8",
				}},
				{BLACK_KING, 1, []SquareCoord{
					"e8",
				}},
			},
		},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			BLACK, 15, 75, 0, 1,
			[]testSquareChecks{
				{FILE_E, RANK_FOUR, WHITE_PAWN},
			},
			[]testPieceChecks{
				{WHITE_PAWN, 8, []SquareCoord{
					"e4",
				}},
			},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			WHITE, 15, 43, 0, 2,
			[]testSquareChecks{
				{FILE_C, RANK_FIVE, BLACK_PAWN},
			},
			[]testPieceChecks{
				{BLACK_PAWN, 8, []SquareCoord{
					"c5",
				}},
			},
		},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			BLACK, 15, 0, 1, 2,
			[]testSquareChecks{
				{FILE_F, RANK_THREE, WHITE_KNIGHT},
			},
			[]testPieceChecks{
				{WHITE_PAWN, 8, []SquareCoord{
					"e4",
				}},
				{WHITE_KNIGHT, 2, []SquareCoord{
					"f3",
				}},
			},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
			WHITE, 0, 0, 0, 1, []testSquareChecks{}, []testPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w q - 0 1",
			WHITE, 1, 0, 0, 1, []testSquareChecks{}, []testPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1",
			WHITE, 2, 0, 0, 1, []testSquareChecks{}, []testPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Q - 0 1",
			WHITE, 4, 0, 0, 1, []testSquareChecks{}, []testPieceChecks{},
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w K - 0 1",
			WHITE, 8, 0, 0, 1, []testSquareChecks{}, []testPieceChecks{},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/PpPP1PPP/RNBQ1RK1 b q - 1 9",
			BLACK, 1, 0, 1, 9,
			[]testSquareChecks{
				{FILE_B, RANK_TWO, BLACK_PAWN},
				{FILE_H, RANK_EIGHT, WHITE_QUEEN},
			},
			[]testPieceChecks{
				{WHITE_QUEEN, 2, []SquareCoord{
					"h8",
				}},
				{WHITE_BISHOP, 2, []SquareCoord{
					"d3",
				}},
				{WHITE_KNIGHT, 2, []SquareCoord{
					"f3",
				}},
				{BLACK_PAWN, 6, []SquareCoord{
					"e6", "b2",
				}},
				{BLACK_KNIGHT, 2, []SquareCoord{
					"f6",
				}},
			},
		},
		{"rnbqkb1Q/pp1p3p/4pn2/8/8/3B1N2/P1PP1PPP/qNBQ1RK1 w q - 0 10",
			WHITE, 1, 0, 0, 10,
			[]testSquareChecks{
				{FILE_A, RANK_ONE, BLACK_QUEEN},
				{FILE_H, RANK_EIGHT, WHITE_QUEEN},
			},
			[]testPieceChecks{
				{WHITE_QUEEN, 2, []SquareCoord{
					"h8",
				}},
				{WHITE_BISHOP, 2, []SquareCoord{
					"d3",
				}},
				{WHITE_KNIGHT, 2, []SquareCoord{
					"f3",
				}},
				{BLACK_PAWN, 5, []SquareCoord{
					"e6",
				}},
				{BLACK_QUEEN, 2, []SquareCoord{
					"a1",
				}},
				{BLACK_KNIGHT, 2, []SquareCoord{
					"f6",
				}},
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		for _, check := range tt.squareChecks {
			squareIndex := squareIndexByFileRank(check.file, check.rank)

			if b.squares[squareIndex] != check.square {
				t.Errorf("square: %v != %v", b.squares[squareIndex], check.square)
			}
		}
		for _, check := range tt.pieceChecks {
			if len(b.pieceIndexes[check.square]) != check.piecesCount {
				t.Errorf("piece indexes length: %v != %v", len(b.pieceIndexes[check.square]), check.piecesCount)
			}
			for _, pieceCoord := range check.pieceCoords {
				pieceIndex, err := squareIndexByCoord(pieceCoord)
				if err != nil {
					t.Error(err)
				}
				if !containsPieceIndex(b.pieceIndexes[check.square], pieceIndex) {
					t.Errorf("unable to find piece index %v in %v", pieceIndex, b.pieceIndexes[check.square])

				}
			}
		}

		if b.sideToMove != tt.sideToMove {
			t.Errorf("side: %v != %v", b.sideToMove, tt.sideToMove)
		}

		if b.castleRights != tt.castleRights {
			t.Errorf("castling rights: %v != %v", b.castleRights, tt.castleRights)
		}

		if b.epIndex != tt.epIndex {
			t.Errorf("en passant index: %v != %v", b.epIndex, tt.epIndex)
		}

		if b.halfMove != tt.halfMove {
			t.Errorf("half move clock: %v != %v", b.halfMove, tt.halfMove)
		}

		if b.fullMove != tt.fullMove {
			t.Errorf("full move clock: %v != %v", b.fullMove, tt.fullMove)
		}
	}
}

func TestColorBySquareIndex(t *testing.T) {
	tests := []struct {
		fen         string
		squareIndex SquareIndex
		expected    Color
	}{
		{STARTING_FEN, 21, BLACK},
		{STARTING_FEN, 91, WHITE},
		{STARTING_FEN, 55, COLOR_NONE},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actual := b.colorBySquareIndex(tt.squareIndex)
		if actual != tt.expected {
			t.Errorf("color: %v != %v", actual, tt.expected)
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

	board, err := newBoard(STARTING_FEN)
	if err != nil {
		t.Error(err)
	}
	if board.toString() != expected {
		t.Errorf("board to string: %v != %v", board.toString(), expected)
	}
}
