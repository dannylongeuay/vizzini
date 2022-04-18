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

func TestSquareKnightAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2", BLACK, "e5", 1,
			[]SquareCoord{"f3"},
		},
		{"k7/8/2n1n3/1n3n2/3Q4/1n3n2/2n1n3/7K w - - 0 1", WHITE, "d4", MAX_SQUARE_KNIGHT_ATTACKERS,
			[]SquareCoord{"b3", "b5", "c2", "c6", "e2", "e6", "f3", "f5"},
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
		attackers := b.squareKnightAttackers(tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			if !containsPieceIndex(attackers, attackerIndex) {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}

func TestSquareDiagonalAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3", BLACK, "c6", 1,
			[]SquareCoord{"b5"},
		},
		{"k7/8/5Q2/2K5/3n4/8/5P2/B7 b - - 0 1", BLACK, "d4", 3,
			[]SquareCoord{"a1", "c5", "f6"},
		},
		{"k6B/8/8/2P5/3q4/4B3/1Q6/7K b - - 0 1", BLACK, "d4", 3,
			[]SquareCoord{"b2", "e3", "h8"},
		},
		{"k6B/8/8/2B5/3q4/4P3/1Q6/7K b - - 0 1", BLACK, "d4", MAX_SQUARE_DIAGONAL_ATTACKERS,
			[]SquareCoord{"b2", "c5", "e3", "h8"},
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
		attackers := b.squareDiagonalAttackers(tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			if !containsPieceIndex(attackers, attackerIndex) {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}

func TestSquareCardinalAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"k7/3Q4/8/8/1R1p1K2/8/3R4/8 b - - 0 1", BLACK, "d4", 3,
			[]SquareCoord{"b4", "d2", "d7"},
		},
		{"k7/3Q4/8/8/1R1pK3/8/3R4/8 b - - 0 1", BLACK, "d4", MAX_SQUARE_CARDINAL_ATTACKERS,
			[]SquareCoord{"b4", "d2", "d7", "e4"},
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
		attackers := b.squareCardinalAttackers(tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			if !containsPieceIndex(attackers, attackerIndex) {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}

func TestSquareAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"k7/8/5Q2/2K5/3n2R1/4P3/2N2P2/B7 b - - 0 1", BLACK, "d4", 6,
			[]SquareCoord{"a1", "c2", "c5", "e3", "f6", "g4"},
		},
		{"k7/8/2N1NQ2/1NBK1N2/1R1n2R1/1N1QPN2/2N1NP2/B7 b - - 0 1", BLACK, "d4", MAX_SQUARE_ATTACKERS,
			[]SquareCoord{
				"a1", "b3", "b4", "b5", "c2", "c5", "c6", "d3",
				"d5", "e2", "e3", "e6", "f3", "f5", "f6", "g4",
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
		attackers := b.squareAttackers(tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			if !containsPieceIndex(attackers, attackerIndex) {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
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

func TestSquareIndexByFileRank(t *testing.T) {
	tests := []struct {
		file     File
		rank     Rank
		expected SquareIndex
	}{
		{FILE_A, RANK_EIGHT, 21},
		{FILE_H, RANK_ONE, 98},
	}
	for _, tt := range tests {
		actual := squareIndexByFileRank(tt.file, tt.rank)
		if actual != tt.expected {
			t.Errorf("square: %v != %v", actual, tt.expected)
		}
	}
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

func TestRankBySquareIndex(t *testing.T) {
	tests := []struct {
		squareIndex SquareIndex
		expected    Rank
	}{
		{21, RANK_EIGHT},
		{32, RANK_SEVEN},
		{43, RANK_SIX},
		{54, RANK_FIVE},
		{65, RANK_FOUR},
		{76, RANK_THREE},
		{87, RANK_TWO},
		{98, RANK_ONE},
		{0, RANK_NONE},
	}
	for _, tt := range tests {
		actual := rankBySquareIndex(tt.squareIndex)
		if actual != tt.expected {
			t.Errorf("file: %v != %v", actual, tt.expected)
		}
	}
}

func TestFileBySquareIndex(t *testing.T) {
	tests := []struct {
		squareIndex SquareIndex
		expected    File
	}{
		{21, FILE_A},
		{32, FILE_B},
		{43, FILE_C},
		{54, FILE_D},
		{65, FILE_E},
		{76, FILE_F},
		{87, FILE_G},
		{98, FILE_H},
		{0, FILE_NONE},
	}
	for _, tt := range tests {
		actual := fileBySquareIndex(tt.squareIndex)
		if actual != tt.expected {
			t.Errorf("file: %v != %v", actual, tt.expected)
		}
	}
}
