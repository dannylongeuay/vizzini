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

var COORD_STRINGS = [BOARD_SQUARES]string{
	"A1", "B1", "C1", "D1", "E1", "F1", "G1", "H1",
	"A2", "B2", "C2", "D2", "E2", "F2", "G2", "H2",
	"A3", "B3", "C3", "D3", "E3", "F3", "G3", "H3",
	"A4", "B4", "C4", "D4", "E4", "F4", "G4", "H4",
	"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5",
	"A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6",
	"A7", "B7", "C7", "D7", "E7", "F7", "G7", "H7",
	"A8", "B8", "C8", "D8", "E8", "F8", "G8", "H8",
}

var MIRROR_COORDS = [64]Coord{
	A8, B8, C8, D8, E8, F8, G8, H8,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A1, B1, C1, D1, E1, F1, G1, H1,
}

var SQUARES = [SQUARE_TYPES]string{
	"EMPTY", "WHITE_PAWN", "WHITE_KNIGHT", "WHITE_BISHOP",
	"WHITE_ROOK", "WHITE_QUEEN", "WHITE_KING", "BLACK_PAWN",
	"BLACK_KNIGHT", "BLACK_BISHOP", "BLACK_ROOK", "BLACK_QUEEN",
	"BLACK_KING",
}

type Board struct {
	squares       []Square
	bbWP          Bitboard
	bbWN          Bitboard
	bbWB          Bitboard
	bbWR          Bitboard
	bbWQ          Bitboard
	bbWK          Bitboard
	bbBP          Bitboard
	bbBN          Bitboard
	bbBB          Bitboard
	bbBR          Bitboard
	bbBQ          Bitboard
	bbBK          Bitboard
	bbWhitePieces Bitboard
	bbBlackPieces Bitboard
	bbAllPieces   Bitboard
	kingCoords    []Coord
	sideToMove    Color
	castleRights  CastleRights
	epCoord       Coord
	halfMove      HalfMove
	fullMove      int
	hash          Hash
	undoIndex     int
	undos         []Undo
}

func NewBoard(fen string) (*Board, error) {
	InitBitboards()
	fenParts := strings.Split(fen, " ")

	if len(fenParts) != 6 {
		return nil, fmt.Errorf("FEN parts: %v != 6", len(fenParts))
	}

	b := Board{}
	b.squares = make([]Square, BOARD_SQUARES)
	b.undos = make([]Undo, MAX_GAME_MOVES)
	b.kingCoords = make([]Coord, PLAYERS)

	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("Ranks: %v != 8", len(ranks))
	}
	coord := Coord(BOARD_SQUARES - 1)
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
				b.squares[coord] = WHITE_PAWN
				b.bbWP.SetBit(coord)
			case "N":
				b.squares[coord] = WHITE_KNIGHT
				b.bbWN.SetBit(coord)
			case "B":
				b.squares[coord] = WHITE_BISHOP
				b.bbWB.SetBit(coord)
			case "R":
				b.squares[coord] = WHITE_ROOK
				b.bbWR.SetBit(coord)
			case "Q":
				b.squares[coord] = WHITE_QUEEN
				b.bbWQ.SetBit(coord)
			case "K":
				b.squares[coord] = WHITE_KING
				b.bbWK.SetBit(coord)
				b.kingCoords[WHITE] = coord
			case "p":
				b.squares[coord] = BLACK_PAWN
				b.bbBP.SetBit(coord)
			case "n":
				b.squares[coord] = BLACK_KNIGHT
				b.bbBN.SetBit(coord)
			case "b":
				b.squares[coord] = BLACK_BISHOP
				b.bbBB.SetBit(coord)
			case "r":
				b.squares[coord] = BLACK_ROOK
				b.bbBR.SetBit(coord)
			case "q":
				b.squares[coord] = BLACK_QUEEN
				b.bbBQ.SetBit(coord)
			case "k":
				b.squares[coord] = BLACK_KING
				b.bbBK.SetBit(coord)
				b.kingCoords[BLACK] = coord
			default:
				return nil, fmt.Errorf("Invalid piece/digit in fen string: %v", char)
			}
			count, err := strconv.Atoi(char)
			if err != nil {
				coord--
			} else {
				for i := 0; i < count; i++ {
					b.squares[coord] = EMPTY
					coord--
				}
			}
		}
	}
	b.UpdateUnionBitboards()

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
			b.castleRights |= CASTLING_RIGHTS_WHITE_KING_MASK
		case "Q":
			b.castleRights |= CASTLING_RIGHTS_WHITE_QUEEN_MASK
		case "k":
			b.castleRights |= CASTLING_RIGHTS_BLACK_KING_MASK
		case "q":
			b.castleRights |= CASTLING_RIGHTS_BLACK_QUEEN_MASK
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

	halfMove, err := strconv.Atoi(fenParts[4])
	if err != nil {
		return nil, err
	}
	b.halfMove = HalfMove(halfMove)

	b.fullMove, err = strconv.Atoi(fenParts[5])
	if err != nil {
		return nil, err
	}

	b.GenerateBoardHash()

	return &b, nil
}

func (b *Board) UpdateUnionBitboards() {
	b.bbWhitePieces = b.bbWP | b.bbWN | b.bbWB | b.bbWR | b.bbWQ | b.bbWK
	b.bbBlackPieces = b.bbBP | b.bbBN | b.bbBB | b.bbBR | b.bbBQ | b.bbBK
	b.bbAllPieces = b.bbWhitePieces | b.bbBlackPieces
}

func (b Board) CopyBoard() Board {
	squares := make([]Square, len(b.squares))
	copy(squares, b.squares)

	undos := make([]Undo, len(b.undos))
	copy(undos, b.undos)

	kingCoords := make([]Coord, len(b.kingCoords))
	copy(kingCoords, b.kingCoords)

	cb := Board{
		squares:       squares,
		bbWP:          b.bbWP,
		bbWN:          b.bbWN,
		bbWB:          b.bbWB,
		bbWR:          b.bbWR,
		bbWQ:          b.bbWQ,
		bbWK:          b.bbWK,
		bbBP:          b.bbBP,
		bbBN:          b.bbBN,
		bbBB:          b.bbBB,
		bbBR:          b.bbBR,
		bbBQ:          b.bbBQ,
		bbBK:          b.bbBK,
		bbWhitePieces: b.bbWhitePieces,
		bbBlackPieces: b.bbBlackPieces,
		bbAllPieces:   b.bbAllPieces,
		kingCoords:    kingCoords,
		sideToMove:    b.sideToMove,
		castleRights:  b.castleRights,
		epCoord:       b.epCoord,
		halfMove:      b.halfMove,
		fullMove:      b.fullMove,
		hash:          b.hash,
		undoIndex:     b.undoIndex,
		undos:         undos,
	}
	return cb
}

type BitboardSquareResult struct {
	bb Bitboard
	sq Square
}

func (b *Board) BitboardSquares() chan BitboardSquareResult {
	c := make(chan BitboardSquareResult, 12)
	c <- BitboardSquareResult{b.bbWP, WHITE_PAWN}
	c <- BitboardSquareResult{b.bbWN, WHITE_KNIGHT}
	c <- BitboardSquareResult{b.bbWB, WHITE_BISHOP}
	c <- BitboardSquareResult{b.bbWR, WHITE_ROOK}
	c <- BitboardSquareResult{b.bbWQ, WHITE_QUEEN}
	c <- BitboardSquareResult{b.bbWK, WHITE_KING}
	c <- BitboardSquareResult{b.bbBP, BLACK_PAWN}
	c <- BitboardSquareResult{b.bbBN, BLACK_KNIGHT}
	c <- BitboardSquareResult{b.bbBB, BLACK_BISHOP}
	c <- BitboardSquareResult{b.bbBR, BLACK_ROOK}
	c <- BitboardSquareResult{b.bbBQ, BLACK_QUEEN}
	c <- BitboardSquareResult{b.bbBK, BLACK_KING}
	close(c)
	return c
}

func (b Board) ToString() string {
	s := ""
	sep := "\n_________________________\n"
	s += sep
	rank := 8
	for i, coord := range MIRROR_COORDS {
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
