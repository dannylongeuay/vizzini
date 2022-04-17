package main

type move struct {
	origin SquareIndex
	target SquareIndex
	kind   MoveKind
}

var PAWN_MOVES_CAPTURE_DIST = [2]int{9, 11}
var KNIGHT_MOVES_DIST = [4]int{8, 12, 19, 21}
var VERTICAL_DIRECTIONS = [2]int{UP, DOWN}
var CARDINAL_DIRECTIONS = [4]int{UP, DOWN, LEFT, RIGHT}
var DIAGONAL_DIRECTIONS = [4][2]int{
	{LEFT, UP},
	{RIGHT, UP},
	{LEFT, DOWN},
	{RIGHT, DOWN},
}

func (b board) generateMoves(side Color) []move {
	moves := make([]move, 0, 256)
	for _, squareIndex := range b.pieceIndexes[side] {
		switch b.squares[squareIndex] {
		case WHITE_PAWN:
			fallthrough
		case BLACK_PAWN:
			pawnMoves := b.generatePawnMoves(side, squareIndex)
			moves = append(moves, pawnMoves...)
			break
		case WHITE_KNIGHT:
			fallthrough
		case BLACK_KNIGHT:
			knightMoves := b.generateKnightMoves(side, squareIndex)
			moves = append(moves, knightMoves...)
			break
		case WHITE_BISHOP:
			fallthrough
		case BLACK_BISHOP:
			break
		case WHITE_ROOK:
			fallthrough
		case BLACK_ROOK:
			break
		case WHITE_QUEEN:
			fallthrough
		case BLACK_QUEEN:
			break
		case WHITE_KING:
			fallthrough
		case BLACK_KING:
			break
		}
	}
	return moves
}

func (b board) generatePawnMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, 16)
	dir := DOWN
	rank := rankBySquareIndex(squareIndex)
	doublePushRank := RANK_SEVEN
	promotionRank := RANK_TWO

	if side == WHITE {
		dir = UP
		doublePushRank = RANK_TWO
		promotionRank = RANK_SEVEN
	}

	// Handle double pawn push
	if rank == doublePushRank {
		target := SquareIndex(VERTICAL_MOVE_DIST*2*dir) + squareIndex
		if b.squares[target] == EMPTY {
			move := move{
				origin: squareIndex,
				target: target,
				kind:   DOUBLE_PAWN_PUSH,
			}
			moves = append(moves, move)
		}
	}

	// Handle quiet move
	target := SquareIndex(VERTICAL_MOVE_DIST*dir) + squareIndex
	if b.squares[target] == EMPTY {
		moveKinds := []MoveKind{QUIET}
		if rank == promotionRank {
			moveKinds = []MoveKind{
				KNIGHT_PROMOTION,
				BISHOP_PROMOTION,
				ROOK_PROMOTION,
				QUEEN_PROMOTION,
			}
		}
		for _, kind := range moveKinds {
			move := move{
				origin: squareIndex,
				target: target,
				kind:   kind,
			}
			moves = append(moves, move)
		}
	}

	// Handle captures
	for _, moveDist := range PAWN_MOVES_CAPTURE_DIST {
		target := SquareIndex(moveDist*dir) + squareIndex
		if b.colorBySquareIndex(target) == side^1 {
			moveKinds := []MoveKind{CAPTURE}
			if rank == promotionRank {
				moveKinds = []MoveKind{
					KNIGHT_PROMOTION_CAPTURE,
					BISHOP_PROMOTION_CAPTURE,
					ROOK_PROMOTION_CAPTURE,
					QUEEN_PROMOTION_CAPTURE,
				}
			}
			for _, kind := range moveKinds {
				move := move{
					origin: squareIndex,
					target: target,
					kind:   kind,
				}
				moves = append(moves, move)
			}
		} else if b.epIndex == target {
			move := move{
				origin: squareIndex,
				target: target,
				kind:   EP_CAPTURE,
			}
			moves = append(moves, move)
		}
	}
	return moves
}

func (b board) generateKnightMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, 8)
	for _, dir := range VERTICAL_DIRECTIONS {
		for _, moveDist := range KNIGHT_MOVES_DIST {
			target := SquareIndex(moveDist*dir) + squareIndex
			if b.squares[target] == EMPTY {
				move := move{
					origin: squareIndex,
					target: target,
					kind:   QUIET,
				}
				moves = append(moves, move)

			} else if b.colorBySquareIndex(target) == side^1 {
				move := move{
					origin: squareIndex,
					target: target,
					kind:   CAPTURE,
				}
				moves = append(moves, move)
			}
		}
	}
	return moves
}
