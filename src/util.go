package main

import (
	"fmt"
	"strings"
)

func colorBySquare(s Square) Color {
	c := COLOR_NONE
	switch s {
	case WHITE_PAWN:
		fallthrough
	case WHITE_KNIGHT:
		fallthrough
	case WHITE_BISHOP:
		fallthrough
	case WHITE_ROOK:
		fallthrough
	case WHITE_QUEEN:
		fallthrough
	case WHITE_KING:
		return WHITE
	case BLACK_PAWN:
		fallthrough
	case BLACK_KNIGHT:
		fallthrough
	case BLACK_BISHOP:
		fallthrough
	case BLACK_ROOK:
		fallthrough
	case BLACK_QUEEN:
		fallthrough
	case BLACK_KING:
		return BLACK
	}
	return c
}

func squareIndexByCoord(s SquareCoord) (SquareIndex, error) {
	var f File
	var r Rank
	coordParts := []rune(s)
	if len(coordParts) != 2 {
		return 0, fmt.Errorf("Invalid chess notation: %v", s)
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
		return 0, fmt.Errorf("Invalid chess file: %v", string(coordParts[0]))
	}
	switch string(coordParts[1]) {
	case "1":
		r = RANK_ONE
	case "2":
		r = RANK_TWO
	case "3":
		r = RANK_THREE
	case "4":
		r = RANK_FOUR
	case "5":
		r = RANK_FIVE
	case "6":
		r = RANK_SIX
	case "7":
		r = RANK_SEVEN
	case "8":
		r = RANK_EIGHT
	default:
		return 0, fmt.Errorf("Invalid chess rank: %v", string(coordParts[1]))

	}
	return squareIndexByFileRank(f, r), nil
}

func squareIndexByFileRank(f File, r Rank) SquareIndex {
	return (21 + SquareIndex(f)) + SquareIndex(r)*10
}

func rankBySquareIndex(s SquareIndex) Rank {
	r := RANK_NONE
	switch (s / 10) % 10 {
	case 2:
		return RANK_EIGHT
	case 3:
		return RANK_SEVEN
	case 4:
		return RANK_SIX
	case 5:
		return RANK_FIVE
	case 6:
		return RANK_FOUR
	case 7:
		return RANK_THREE
	case 8:
		return RANK_TWO
	case 9:
		return RANK_ONE
	}
	return r
}

func fileBySquareIndex(s SquareIndex) File {
	f := FILE_NONE
	switch s % 10 {
	case 1:
		return FILE_A
	case 2:
		return FILE_B
	case 3:
		return FILE_C
	case 4:
		return FILE_D
	case 5:
		return FILE_E
	case 6:
		return FILE_F
	case 7:
		return FILE_G
	case 8:
		return FILE_H
	}
	return f
}
