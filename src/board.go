package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
   +----+----+----+----+----+----+----+----+
 8 | 56 | 57 | 58 | 59 | 60 | 61 | 62 | 63 |  8th rank
   +----+----+----+----+----+----+----+----+
 7 | 48 | 49 | 50 | 51 | 52 | 53 | 54 | 55 |  7th rank
   +----+----+----+----+----+----+----+----+
 6 | 40 | 41 | 42 | 43 | 44 | 45 | 46 | 47 |  6th rank
   +----+----+----+----+----+----+----+----+
 5 | 32 | 33 | 34 | 35 | 36 | 37 | 38 | 39 |  5th rank
   +----+----+----+----+----+----+----+----+
 4 | 24 | 25 | 26 | 27 | 28 | 29 | 30 | 31 |  4th rank
   +----+----+----+----+----+----+----+----+
 3 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 |  3rd rank
   +----+----+----+----+----+----+----+----+
 2 |  8 |  9 | 10 | 11 | 12 | 13 | 14 | 15 |  2nd rank
   +----+----+----+----+----+----+----+----+
 1 |  0 |  1 |  2 |  3 |  4 |  5 |  6 |  7 |  1st rank
   +----+----+----+----+----+----+----+----+
     A    B    C    D    E    F    G    H - file(s)
*/

type Undo struct {
	mv             Move
	capturedSquare Square
	castleRights   CastleRights
	epCoord        Coord
	halfMove       int
	hash           Hash
}

var COORD_MAP = [64]string{
	"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1",
	"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2",
	"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3",
	"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4",
	"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5",
	"A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6",
	"A7", "B7", "C7", "D7", "E7", "F7", "G7", "H7",
	"A8", "B8", "C8", "D8", "E8", "F8", "G8", "H8",
}

type Board struct {
	squares        []Square
	whiteKingCoord Coord
	blackKingCoord Coord
	sideToMove     Color
	/*
		castleRights
		 0000 	0 0 0 0
		unused	K Q k q
	*/
	castleRights CastleRights
	epCoord      Coord
	halfMove     int
	fullMove     int
	hash         Hash
	undoIndex    int
	undos        []Undo
}

func NewBoard(fen string) (*Board, error) {
	fenParts := strings.Split(fen, " ")

	if len(fenParts) != 6 {
		return nil, fmt.Errorf("FEN parts: %v != 6", len(fenParts))
	}

	b := Board{}
	b.squares = make([]Square, BOARD_SQUARES)
	b.undos = make([]Undo, MAX_GAME_MOVES)

	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("Ranks: %v != 8", len(ranks))
	}
	squareIndex := BOARD_SQUARES - 1
	for _, rank := range ranks {
		runes := []rune(rank)
		for r := len(runes) - 1; r >= 0; r-- {
			char := string(runes[r])
			switch char {
			case "1":
			case "2":
			case "3":
			case "4":
			case "5":
			case "6":
			case "7":
			case "8":
			case "P":
				b.squares[squareIndex] = WHITE_PAWN
			case "N":
				b.squares[squareIndex] = WHITE_KNIGHT
			case "B":
				b.squares[squareIndex] = WHITE_BISHOP
			case "R":
				b.squares[squareIndex] = WHITE_ROOK
			case "Q":
				b.squares[squareIndex] = WHITE_QUEEN
			case "K":
				b.squares[squareIndex] = WHITE_KING
				b.whiteKingCoord = Coord(squareIndex)
			case "p":
				b.squares[squareIndex] = BLACK_PAWN
			case "n":
				b.squares[squareIndex] = BLACK_KNIGHT
			case "b":
				b.squares[squareIndex] = BLACK_BISHOP
			case "r":
				b.squares[squareIndex] = BLACK_ROOK
			case "q":
				b.squares[squareIndex] = BLACK_QUEEN
			case "k":
				b.squares[squareIndex] = BLACK_KING
				b.blackKingCoord = Coord(squareIndex)
			default:
				return nil, fmt.Errorf("Invalid piece/digit in fen string: %v", char)
			}
			count, err := strconv.Atoi(char)
			if err != nil {
				squareIndex--
			} else {
				for i := 0; i < count; i++ {
					b.squares[squareIndex] = EMPTY
					squareIndex--
				}
			}
		}
	}

	sideToMove := fenParts[1]
	if sideToMove == "w" {
		b.sideToMove = WHITE
	} else if sideToMove == "b" {
		b.sideToMove = BLACK
	} else {
		return nil, fmt.Errorf("Invalid side to move in fen string: %v", sideToMove)
	}

	castlingRights := fenParts[2]
	for _, char := range []rune(castlingRights) {
		switch string(char) {
		case "-":
		case "K":
			b.castleRights |= 1 << 3
		case "Q":
			b.castleRights |= 1 << 2
		case "k":
			b.castleRights |= 1 << 1
		case "q":
			b.castleRights |= 1
		default:
			return nil, fmt.Errorf("Invalid castling rights in fen string: %v", string(char))
		}
	}

	var err error
	if fenParts[3] != "-" {
		b.epCoord, err = StringToCoord(fenParts[3])
		if err != nil {
			return nil, err
		}
	}

	b.halfMove, err = strconv.Atoi(fenParts[4])
	if err != nil {
		return nil, err
	}

	b.fullMove, err = strconv.Atoi(fenParts[5])
	if err != nil {
		return nil, err
	}

	b.generateBoardHash()

	return &b, nil
}

func (b Board) CopyBoard() Board {
	squares := make([]Square, len(b.squares))
	copy(squares, b.squares)

	undos := make([]Undo, len(b.undos))
	copy(undos, b.undos)

	cb := Board{
		squares:        squares,
		whiteKingCoord: b.whiteKingCoord,
		blackKingCoord: b.blackKingCoord,
		sideToMove:     b.sideToMove,
		castleRights:   b.castleRights,
		epCoord:        b.epCoord,
		halfMove:       b.halfMove,
		fullMove:       b.fullMove,
		hash:           b.hash,
		undoIndex:      b.undoIndex,
		undos:          undos,
	}
	return cb
}

var PRINT_MAP = [64]Coord{
	A8, B8, C8, D8, E8, F8, G8, H8,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A1, B1, C1, D1, E1, F1, G1, H1,
}

func (b Board) ToString() string {
	s := ""
	sep := "\n_________________________\n"
	s += sep
	rank := 8
	for i, coord := range PRINT_MAP {
		if i != 0 && i%8 == 0 {
			s += "| "
			s += strconv.Itoa(rank)
			rank--
			s += sep
		}
		s += "|"
		square := b.squares[coord]
		switch square {
		case EMPTY:
			s += " "
		case WHITE_PAWN:
			s += "♟"
		case WHITE_KNIGHT:
			s += "♞"
		case WHITE_BISHOP:
			s += "♝"
		case WHITE_ROOK:
			s += "♜"
		case WHITE_QUEEN:
			s += "♛"
		case WHITE_KING:
			s += "♚"
		case BLACK_PAWN:
			s += "♙"
		case BLACK_KNIGHT:
			s += "♘"
		case BLACK_BISHOP:
			s += "♗"
		case BLACK_ROOK:
			s += "♖"
		case BLACK_QUEEN:
			s += "♕"
		case BLACK_KING:
			s += "♔"
		default:
			s += "♠"
		}
		s += " "
	}
	s += "| 1"
	s += sep
	s += " A  B  C  D  E  F  G  H"
	return s
}
