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

func StringToCoord(s string) (Coord, error) {
	var coord Coord

	coordParts := []rune(s)
	if len(coordParts) != 2 {
		return 0, fmt.Errorf("Invalid chess notation: %v", s)
	}

	switch strings.ToLower(string(coordParts[0])) {
	case "a":
	case "b":
		coord += 1
	case "c":
		coord += 2
	case "d":
		coord += 3
	case "e":
		coord += 4
	case "f":
		coord += 5
	case "g":
		coord += 6
	case "h":
		coord += 7
	default:
		return 0, fmt.Errorf("Invalid chess file: %v", string(coordParts[0]))
	}

	switch string(coordParts[1]) {
	case "1":
	case "2":
		coord += 8
	case "3":
		coord += 16
	case "4":
		coord += 24
	case "5":
		coord += 32
	case "6":
		coord += 40
	case "7":
		coord += 48
	case "8":
		coord += 56
	default:
		return 0, fmt.Errorf("Invalid chess rank: %v", string(coordParts[1]))
	}

	return coord, nil
}

func IntPow(a int, b int) int {
	result := 1

	for b > 0 {
		if (b & 1) > 0 {
			result *= a
		}
		b = b >> 1
		a *= a
	}

	return result

}

func FilterMovesByKind(kind MoveKind, moves *[]Move) {
	n := 0
	for _, m := range *moves {
		var mu MoveUnpacked
		m.Unpack(&mu)
		if mu.moveKind == kind {
			(*moves)[n] = m
			n++
		}
	}
	*moves = (*moves)[:n]
}
