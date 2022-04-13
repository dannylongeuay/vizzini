package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Square int
type Color int
type File int
type Rank int
type Bits uint8

const (
	Invalid Square = iota
	WhitePawn
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
	Empty
)

const (
	White Color = iota
	Black
	Both
)

const (
	FILE_A File = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
	FILE_None
)

const (
	RANK_8 Rank = iota
	RANK_7
	RANK_6
	RANK_5
	RANK_4
	RANK_3
	RANK_2
	RANK_1
	RANK_None
)

const BoardSquares int = 120

const StartingFEN string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var SquareIndexes64 = [64]int{
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
	squares      []Square
	sideToMove   Color
	fiftyMove    int
	castleRights Bits
	epIndex      int
	halfMove     int
	fullMove     int
}

func newBoard(fen string) (*board, error) {
	fenParts := strings.Split(fen, " ")

	if len(fenParts) != 6 {
		return nil, fmt.Errorf("FEN parts: %v != 6", len(fenParts))
	}

	b := board{}
	b.squares = make([]Square, BoardSquares)

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
				b.squares[squareIndex] = WhitePawn
				break
			case "N":
				b.squares[squareIndex] = WhiteKnight
				break
			case "B":
				b.squares[squareIndex] = WhiteBishop
				break
			case "R":
				b.squares[squareIndex] = WhiteRook
				break
			case "Q":
				b.squares[squareIndex] = WhiteQueen
				break
			case "K":
				b.squares[squareIndex] = WhiteKing
				break
			case "p":
				b.squares[squareIndex] = BlackPawn
				break
			case "n":
				b.squares[squareIndex] = BlackKnight
				break
			case "b":
				b.squares[squareIndex] = BlackBishop
				break
			case "r":
				b.squares[squareIndex] = BlackRook
				break
			case "q":
				b.squares[squareIndex] = BlackQueen
				break
			case "k":
				b.squares[squareIndex] = BlackKing
				break
			default:
				return nil, fmt.Errorf("Invalid piece/digit in fen string: %v", string(char))
			}
			count, err := strconv.Atoi(string(char))
			if err != nil {
				squareIndex64++
			} else {
				for i := 0; i < count; i++ {
					b.squares[squareIndex] = Empty
					squareIndex64++
					squareIndex = SquareIndexes64[squareIndex64]
				}
			}
		}
	}

	sideToMove := fenParts[1]
	if sideToMove == "w" {
		b.sideToMove = White
	} else if sideToMove == "b" {
		b.sideToMove = Black
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

	b.epIndex = -1
	if fenParts[3] != "-" {
		var err error
		b.epIndex, err = squareIndexByCoord(fenParts[3])
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

	return &b, nil
}

func (b board) toString() string {
	s := ""
	s += "\n_________________________\n"
	for i, index := range SquareIndexes64 {
		square := b.squares[index]
		if i != 0 && i%8 == 0 {
			s += "|\n_________________________\n"
		}
		s += "|"
		switch square {
		case Empty:
			s += " "
			break
		case WhitePawn:
			s += "♟"
			break
		case WhiteKnight:
			s += "♞"
			break
		case WhiteBishop:
			s += "♝"
			break
		case WhiteRook:
			s += "♜"
			break
		case WhiteQueen:
			s += "♛"
			break
		case WhiteKing:
			s += "♚"
			break
		case BlackPawn:
			s += "♙"
			break
		case BlackKnight:
			s += "♘"
			break
		case BlackBishop:
			s += "♗"
			break
		case BlackRook:
			s += "♖"
			break
		case BlackQueen:
			s += "♕"
			break
		case BlackKing:
			s += "♔"
			break
		}
		s += " "
	}
	s += "|\n_________________________\n"
	return s
}

func squareIndexByCoord(s string) (int, error) {
	var f File
	var r Rank
	coordParts := []rune(s)
	if len(coordParts) != 2 {
		return -1, fmt.Errorf("Invalid chess notation: %v", s)
	}
	switch strings.ToLower(string(coordParts[0])) {
	case "a":
		f = FILE_A
	case "b":
		f = FILE_B
	case "c":
		f = FILE_C
	case "d":
		f = FILE_D
	case "e":
		f = FILE_E
	case "f":
		f = FILE_F
	case "g":
		f = FILE_G
	case "h":
		f = FILE_H
	default:
		return -1, fmt.Errorf("Invalid chess file: %v", string(coordParts[0]))
	}
	switch string(coordParts[1]) {
	case "1":
		r = RANK_1
	case "2":
		r = RANK_2
	case "3":
		r = RANK_3
	case "4":
		r = RANK_4
	case "5":
		r = RANK_5
	case "6":
		r = RANK_6
	case "7":
		r = RANK_7
	case "8":
		r = RANK_8
	default:
		return -1, fmt.Errorf("Invalid chess rank: %v", string(coordParts[1]))

	}
	return squareIndexByFileRank(f, r), nil
}

func squareIndexByFileRank(f File, r Rank) int {
	return (21 + int(f)) + int(r)*10
}
