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
			i, err := strconv.Atoi(string(char))
			if err != nil {
				squareIndex64++
			} else {
				squareIndex64 += i
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
	return &b, nil
}

func squareByFileRank(f File, r Rank) int {
	return (21 + int(f)) + int(r)*10
}
