package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Piece int
type Color int
type File int
type Rank int

const (
	Empty Piece = iota
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

var PiecesIndex64 = [64]int{
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
	pieces       []Piece
	side         Color
	fiftyMove    int
	castleRights int
}

func newBoard(fen string) (*board, error) {
	fenParts := strings.Split(fen, " ")
	if len(fenParts) != 6 {
		return nil, fmt.Errorf("FEN parts: %v != 6", len(fenParts))
	}
	b := board{}
	b.pieces = make([]Piece, BoardSquares)
	ranks := strings.Split(fenParts[0], "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("Ranks: %v != 8", len(ranks))
	}
	pieceIndex64 := 0
	for _, rank := range ranks {
		for _, char := range []rune(rank) {
			pieceIndex := PiecesIndex64[pieceIndex64]
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
				b.pieces[pieceIndex] = WhitePawn
				break
			case "N":
				b.pieces[pieceIndex] = WhiteKnight
				break
			case "B":
				b.pieces[pieceIndex] = WhiteBishop
				break
			case "R":
				b.pieces[pieceIndex] = WhiteRook
				break
			case "Q":
				b.pieces[pieceIndex] = WhiteQueen
				break
			case "K":
				b.pieces[pieceIndex] = WhiteKing
				break
			case "p":
				b.pieces[pieceIndex] = BlackPawn
				break
			case "n":
				b.pieces[pieceIndex] = BlackKnight
				break
			case "b":
				b.pieces[pieceIndex] = BlackBishop
				break
			case "r":
				b.pieces[pieceIndex] = BlackRook
				break
			case "q":
				b.pieces[pieceIndex] = BlackQueen
				break
			case "k":
				b.pieces[pieceIndex] = BlackKing
				break
			default:
				return nil, fmt.Errorf("Invalid char in fen string: %v", string(char))
			}
			i, err := strconv.Atoi(string(char))
			if err != nil {
				pieceIndex64++
			} else {
				pieceIndex64 += i
			}
		}
	}
	return &b, nil
}

func squareByFileRank(f File, r Rank) int {
	return (21 + int(f)) + int(r)*10
}
