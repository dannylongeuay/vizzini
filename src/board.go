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
	squares      []Square
	pieceIndexes map[Square][]SquareIndex // piecesIndexes[WHITE_PAWN] = [81, 82, ...]
	sideToMove   Color
	fiftyMove    int
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
}

func newBoard(fen string) (*board, error) {
	fenParts := strings.Split(fen, " ")

	if len(fenParts) != 6 {
		return nil, fmt.Errorf("FEN parts: %v != 6", len(fenParts))
	}

	b := board{}
	b.squares = make([]Square, BOARD_SQUARES)
	b.pieceIndexes = make(map[Square][]SquareIndex)

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
				break
			default:
				return nil, fmt.Errorf("Invalid piece/digit in fen string: %v", string(char))
			}
			count, err := strconv.Atoi(string(char))
			if err != nil {
				square := b.squares[squareIndex]
				b.pieceIndexes[square] = append(b.pieceIndexes[square], squareIndex)
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

	return &b, nil
}

func (b board) colorBySquareIndex(s SquareIndex) Color {
	return colorBySquare(b.squares[s])
}

func (b board) squareKnightAttackers(side Color, squareIndex SquareIndex) []SquareIndex {
	attackers := make([]SquareIndex, 0, MAX_SQUARE_KNIGHT_ATTACKERS)

	enemyKnight := BLACK_KNIGHT

	if side == BLACK {
		enemyKnight = WHITE_KNIGHT
	}

	for _, dir := range MOVE_DIRECTIONS {
		for _, moveDist := range KNIGHT_MOVE_DISTS {
			originIndex := SquareIndex(moveDist*dir) + squareIndex
			if b.squares[originIndex] == enemyKnight {
				attackers = append(attackers, originIndex)
			}
		}
	}
	return attackers
}

func (b board) squareDiagonalAttackers(side Color, squareIndex SquareIndex) []SquareIndex {
	attackers := make([]SquareIndex, 0, MAX_SQUARE_DIAGONAL_ATTACKERS)

	pawnAttackDir := POSITIVE_DIR

	enemyPawn := BLACK_PAWN
	enemyBishop := BLACK_BISHOP
	enemyQueen := BLACK_QUEEN
	enemyKing := BLACK_KING

	if side == BLACK {
		pawnAttackDir = NEGATIVE_DIR

		enemyPawn = WHITE_PAWN
		enemyBishop = WHITE_BISHOP
		enemyQueen = WHITE_QUEEN
		enemyKing = WHITE_KING
	}

	for _, dir := range MOVE_DIRECTIONS {
		for _, moveDist := range DIAGONAL_MOVE_DISTS {
			for i := 1; i < MAX_MOVE_RANGE; i++ {
				originIndex := SquareIndex(moveDist*dir*i) + squareIndex
				originSquare := b.squares[originIndex]
				if originSquare == EMPTY {
					continue
				}
				switch originSquare {
				case enemyBishop:
					fallthrough
				case enemyQueen:
					attackers = append(attackers, originIndex)
					break
				case enemyPawn:
					if i == 1 && dir*-1 == pawnAttackDir {
						attackers = append(attackers, originIndex)
					}
					break
				case enemyKing:
					if i == 1 {
						attackers = append(attackers, originIndex)
					}
					break
				}
				// Reached invalid, friendly, or non-diagonal moving piece
				break
			}
		}
	}
	return attackers
}

func (b board) squareCardinalAttackers(side Color, squareIndex SquareIndex) []SquareIndex {
	attackers := make([]SquareIndex, 0, MAX_SQUARE_CARDINAL_ATTACKERS)

	enemyRook := BLACK_ROOK
	enemyQueen := BLACK_QUEEN
	enemyKing := BLACK_KING

	if side == BLACK {
		enemyRook = WHITE_ROOK
		enemyQueen = WHITE_QUEEN
		enemyKing = WHITE_KING
	}

	for _, dir := range MOVE_DIRECTIONS {
		for _, moveDist := range CARDINAL_MOVE_DISTS {
			for i := 1; i < MAX_MOVE_RANGE; i++ {
				originIndex := SquareIndex(moveDist*dir*i) + squareIndex
				originSquare := b.squares[originIndex]
				if originSquare == EMPTY {
					continue
				}
				switch originSquare {
				case enemyRook:
					fallthrough
				case enemyQueen:
					attackers = append(attackers, originIndex)
					break
				case enemyKing:
					if i == 1 {
						attackers = append(attackers, originIndex)
					}
					break
				}
				// Reached invalid, friendly, or non-cardinal moving piece
				break
			}
		}
	}
	return attackers
}

func (b board) squareAttackers(side Color, squareIndex SquareIndex) []SquareIndex {
	attackers := make([]SquareIndex, 0, MAX_SQUARE_ATTACKERS)
	knightAttackers := b.squareKnightAttackers(side, squareIndex)
	attackers = append(attackers, knightAttackers...)
	diagonalAttackers := b.squareDiagonalAttackers(side, squareIndex)
	attackers = append(attackers, diagonalAttackers...)
	cardinalAttackers := b.squareCardinalAttackers(side, squareIndex)
	attackers = append(attackers, cardinalAttackers...)
	return attackers
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
