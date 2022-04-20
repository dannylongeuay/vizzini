package main

import (
	"fmt"
	"strconv"
	"strings"
)

var SquareIndexes64 = [64]SquareIndex{
	21, 22, 23, 24, 25, 26, 27, 28,
	31, 32, 33, 34, 35, 36, 37, 38,
	41, 42, 43, 44, 45, 46, 47, 48,
	51, 52, 53, 54, 55, 56, 57, 58,
	61, 62, 63, 64, 65, 66, 67, 68,
	71, 72, 73, 74, 75, 76, 77, 78,
	81, 82, 83, 84, 85, 86, 87, 88,
	91, 92, 93, 94, 95, 96, 97, 98,
}

type board struct {
	squares        []Square
	pieceSets      map[Square]map[SquareIndex]bool
	whiteKingIndex SquareIndex
	blackKindIndex SquareIndex
	sideToMove     Color
	fiftyMove      int
	/*
		castleRights
		0000 0001 = Black king can castle queenside
		0000 0010 = Black king can castle kingside
		0000 0100 = White king can castle queenside
		0000 1000 = White king can castle kingside
	*/
	castleRights CastleRights
	epIndex      SquareIndex
	halfMove     int
	fullMove     int
	hash         uint64
}

func newBoard(fen string) (*board, error) {
	fenParts := strings.Split(fen, " ")

	if len(fenParts) != 6 {
		return nil, fmt.Errorf("FEN parts: %v != 6", len(fenParts))
	}

	b := board{}
	b.squares = make([]Square, BOARD_SQUARES)
	b.pieceSets = makePieceSets()

	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("Ranks: %v != 8", len(ranks))
	}
	squareIndex64 := 0
	for _, rank := range ranks {
		for _, char := range []rune(rank) {
			squareIndex := SquareIndexes64[squareIndex64]
			switch string(char) {
			case "1":
			case "2":
			case "3":
			case "4":
			case "5":
			case "6":
			case "7":
			case "8":
				break
			case "P":
				b.squares[squareIndex] = WHITE_PAWN
				break
			case "N":
				b.squares[squareIndex] = WHITE_KNIGHT
				break
			case "B":
				b.squares[squareIndex] = WHITE_BISHOP
				break
			case "R":
				b.squares[squareIndex] = WHITE_ROOK
				break
			case "Q":
				b.squares[squareIndex] = WHITE_QUEEN
				break
			case "K":
				b.squares[squareIndex] = WHITE_KING
				b.whiteKingIndex = squareIndex
				break
			case "p":
				b.squares[squareIndex] = BLACK_PAWN
				break
			case "n":
				b.squares[squareIndex] = BLACK_KNIGHT
				break
			case "b":
				b.squares[squareIndex] = BLACK_BISHOP
				break
			case "r":
				b.squares[squareIndex] = BLACK_ROOK
				break
			case "q":
				b.squares[squareIndex] = BLACK_QUEEN
				break
			case "k":
				b.squares[squareIndex] = BLACK_KING
				b.blackKindIndex = squareIndex
				break
			default:
				return nil, fmt.Errorf("Invalid piece/digit in fen string: %v", string(char))
			}
			count, err := strconv.Atoi(string(char))
			if err != nil {
				square := b.squares[squareIndex]
				b.pieceSets[square][squareIndex] = true
				squareIndex64++
			} else {
				for i := 0; i < count; i++ {
					emptySquareIndex := SquareIndexes64[squareIndex64]
					b.squares[emptySquareIndex] = EMPTY
					squareIndex64++
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
			break
		case "K":
			b.castleRights |= 1 << 3
			break
		case "Q":
			b.castleRights |= 1 << 2
			break
		case "k":
			b.castleRights |= 1 << 1
			break
		case "q":
			b.castleRights |= 1
			break
		default:
			return nil, fmt.Errorf("Invalid castling rights in fen string: %v", string(char))
		}
	}

	if fenParts[3] != "-" {
		var err error
		b.epIndex, err = squareIndexByCoord(SquareCoord(fenParts[3]))
		if err != nil {
			return nil, err
		}
	}

	var err error
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

func (b board) colorBySquareIndex(s SquareIndex) Color {
	return colorBySquare(b.squares[s])
}

func (b board) toString() string {
	s := ""
	sep := "\n_________________________\n"
	s += sep
	rank := 8
	for i, index := range SquareIndexes64 {
		square := b.squares[index]
		if i != 0 && i%8 == 0 {
			s += "| "
			s += strconv.Itoa(rank)
			rank--
			s += sep
		}
		s += "|"
		switch square {
		case EMPTY:
			s += " "
			break
		case WHITE_PAWN:
			s += "♟"
			break
		case WHITE_KNIGHT:
			s += "♞"
			break
		case WHITE_BISHOP:
			s += "♝"
			break
		case WHITE_ROOK:
			s += "♜"
			break
		case WHITE_QUEEN:
			s += "♛"
			break
		case WHITE_KING:
			s += "♚"
			break
		case BLACK_PAWN:
			s += "♙"
			break
		case BLACK_KNIGHT:
			s += "♘"
			break
		case BLACK_BISHOP:
			s += "♗"
			break
		case BLACK_ROOK:
			s += "♖"
			break
		case BLACK_QUEEN:
			s += "♕"
			break
		case BLACK_KING:
			s += "♔"
			break
		default:
			s += "♠"
			break
		}
		s += " "
	}
	s += "| 1"
	s += sep
	s += " A  B  C  D  E  F  G  H"
	return s
}

func makePieceSets() map[Square]map[SquareIndex]bool {
	pieceSets := make(map[Square]map[SquareIndex]bool)
	pieceSets[WHITE_PAWN] = make(map[SquareIndex]bool)
	pieceSets[WHITE_KNIGHT] = make(map[SquareIndex]bool)
	pieceSets[WHITE_BISHOP] = make(map[SquareIndex]bool)
	pieceSets[WHITE_ROOK] = make(map[SquareIndex]bool)
	pieceSets[WHITE_QUEEN] = make(map[SquareIndex]bool)
	pieceSets[WHITE_KING] = make(map[SquareIndex]bool)
	pieceSets[BLACK_PAWN] = make(map[SquareIndex]bool)
	pieceSets[BLACK_KNIGHT] = make(map[SquareIndex]bool)
	pieceSets[BLACK_BISHOP] = make(map[SquareIndex]bool)
	pieceSets[BLACK_ROOK] = make(map[SquareIndex]bool)
	pieceSets[BLACK_QUEEN] = make(map[SquareIndex]bool)
	pieceSets[BLACK_KING] = make(map[SquareIndex]bool)
	return pieceSets
}
