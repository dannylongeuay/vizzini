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

func (b *Board) UCIParseMove(s string) (Move, error) {
	sParts := []rune(s)

	if len(sParts) < 4 || len(sParts) > 5 {
		return 0, fmt.Errorf("Invalid UCI move notation: %v", s)
	}

	originCoord, err := StringToCoord(string(sParts[:2]))
	if err != nil {
		return 0, err
	}

	dstCoord, err := StringToCoord(string(sParts[2:4]))
	if err != nil {
		return 0, err
	}

	originSquare := b.squares[originCoord]
	dstSquare := b.squares[dstCoord]

	moveKind := QUIET
	if dstSquare != EMPTY {
		moveKind = CAPTURE
	}

	if originSquare == WHITE_PAWN && dstCoord-originCoord == 16 {
		moveKind = DOUBLE_PAWN_PUSH
	}

	if originSquare == BLACK_PAWN && originCoord-dstCoord == 16 {
		moveKind = DOUBLE_PAWN_PUSH
	}

	if dstCoord == b.epCoord && (originSquare == WHITE_PAWN || originSquare == BLACK_PAWN) {
		moveKind = EP_CAPTURE
	}

	if originSquare == WHITE_KING {
		switch dstCoord {
		case G1:
			if b.castleRights&CASTLING_RIGHTS_WHITE_KING_MASK > 0 {
				moveKind = KING_CASTLE
			}
		case C1:
			if b.castleRights&CASTLING_RIGHTS_WHITE_QUEEN_MASK > 0 {
				moveKind = QUEEN_CASTLE
			}
		}
	}

	if originSquare == BLACK_KING {
		switch dstCoord {
		case G8:
			if b.castleRights&CASTLING_RIGHTS_BLACK_KING_MASK > 0 {
				moveKind = KING_CASTLE
			}
		case C8:
			if b.castleRights&CASTLING_RIGHTS_BLACK_QUEEN_MASK > 0 {
				moveKind = QUEEN_CASTLE
			}
		}
	}

	if len(sParts) == 5 {
		switch sParts[4] {
		case 'n':
			if moveKind == CAPTURE {
				moveKind = KNIGHT_PROMOTION_CAPTURE
			} else {
				moveKind = KNIGHT_PROMOTION
			}
		case 'b':
			if moveKind == CAPTURE {
				moveKind = BISHOP_PROMOTION_CAPTURE
			} else {
				moveKind = BISHOP_PROMOTION
			}
		case 'r':
			if moveKind == CAPTURE {
				moveKind = ROOK_PROMOTION_CAPTURE
			} else {
				moveKind = ROOK_PROMOTION
			}
		case 'q':
			if moveKind == CAPTURE {
				moveKind = QUEEN_PROMOTION_CAPTURE
			} else {
				moveKind = QUEEN_PROMOTION
			}
		default:
			return 0, fmt.Errorf("Invalid promotion: %v", s)
		}
	}
	return NewMove(originCoord, dstCoord, originSquare, dstSquare, moveKind), nil
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
