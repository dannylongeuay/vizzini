package main

var SQUARE_SCORES = [SQUARE_TYPES]int{
	0, 100, 300, 300, 500, 900, 10000,
	-100, -300, -300, -500, -900, -10000,
}

var PAWN_COORD_SCORES = [BOARD_SQUARES]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 5, 5, -10, -10, 5, 5, 5,
	5, 0, 0, 5, 5, 0, 0, 5,
	5, 5, 5, 20, 20, 5, 5, 5,
	10, 10, 15, 25, 25, 15, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	30, 30, 30, 40, 40, 30, 30, 30,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KNIGHT_COORD_SCORES = [BOARD_SQUARES]int{
	-5, -10, -5, -5, -5, -5, -10, -5,
	-5, 0, 0, 5, 5, 0, 0, -5,
	-5, 5, 15, 15, 15, 15, 5, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 5, 15, 20, 20, 15, 5, -5,
	-5, 0, 0, 10, 10, 0, 0, -5,
	-5, -5, -5, -5, -5, -5, -5, -5,
}

var BISHOP_COORD_SCORES = [BOARD_SQUARES]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 0, 0, 10, 0, 0,
	0, 10, 0, 0, 0, 0, 10, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var ROOK_COORD_SCORES = [BOARD_SQUARES]int{
	0, 0, 5, 10, 10, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	20, 20, 20, 20, 20, 20, 20, 20,
	10, 10, 10, 10, 10, 10, 10, 10,
}

var QUEEN_COORD_SCORES = [BOARD_SQUARES]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-20, -10, -10, 0, 0, -10, -10, -20,
}

var KING_COORD_SCORES = [BOARD_SQUARES]int{
	0, 0, 10, -5, -5, -5, 10, 0,
	0, 0, 0, -5, -5, -5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func (b *Board) Evaluate() int {

	var score int

	for bbSquare := range b.BitboardSquares() {
		for bbSquare.bb > 0 {
			score += SQUARE_SCORES[bbSquare.sq]

			coord := bbSquare.bb.PopLSB()
			mirror := MIRROR_COORDS[coord]

			switch bbSquare.sq {
			case WHITE_PAWN:
				score += PAWN_COORD_SCORES[coord]
			case WHITE_KNIGHT:
				score += KNIGHT_COORD_SCORES[coord]
			case WHITE_BISHOP:
				score += BISHOP_COORD_SCORES[coord]
			case WHITE_ROOK:
				score += ROOK_COORD_SCORES[coord]
			case WHITE_QUEEN:
				score += QUEEN_COORD_SCORES[coord]
			case WHITE_KING:
				score += KING_COORD_SCORES[coord]
			case BLACK_PAWN:
				score -= PAWN_COORD_SCORES[mirror]
			case BLACK_KNIGHT:
				score -= KNIGHT_COORD_SCORES[mirror]
			case BLACK_BISHOP:
				score -= BISHOP_COORD_SCORES[mirror]
			case BLACK_ROOK:
				score -= ROOK_COORD_SCORES[mirror]
			case BLACK_QUEEN:
				score -= QUEEN_COORD_SCORES[mirror]
			case BLACK_KING:
				score -= KING_COORD_SCORES[mirror]
			}
		}
	}

	if b.sideToMove == WHITE {
		return score
	}
	return -score
}
