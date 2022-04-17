package main

type move struct {
	origin SquareIndex
	target SquareIndex
	kind   MoveKind
}

var DIAGONAL_DISTS = [2]int{9, 11}
var KNIGHT_MOVES_DIST = [4]int{8, 12, 19, 21}
var VERTICAL_DIRECTIONS = [2]int{UP, DOWN}
var HORIZONTAL_DIRECTIONS = [2]int{LEFT, RIGHT}

func (b board) generateMoves(side Color) []move {
	moves := make([]move, 0, MAX_GENERATED_MOVES)
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
			bishopMoves := b.generateBishopMoves(side, squareIndex)
			moves = append(moves, bishopMoves...)
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
	moves := make([]move, 0, MAX_PAWN_MOVES)
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
		targetIndex := SquareIndex(VERTICAL_MOVE_DIST*2*dir) + squareIndex
		if b.squares[targetIndex] == EMPTY {
			move := move{
				origin: squareIndex,
				target: targetIndex,
				kind:   DOUBLE_PAWN_PUSH,
			}
			moves = append(moves, move)
		}
	}

	// Handle quiet move
	targetIndex := SquareIndex(VERTICAL_MOVE_DIST*dir) + squareIndex
	if b.squares[targetIndex] == EMPTY {
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
				target: targetIndex,
				kind:   kind,
			}
			moves = append(moves, move)
		}
	}

	// Handle captures
	for _, moveDist := range DIAGONAL_DISTS {
		targetIndex := SquareIndex(moveDist*dir) + squareIndex
		if b.colorBySquareIndex(targetIndex) == side^1 {
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
					target: targetIndex,
					kind:   kind,
				}
				moves = append(moves, move)
			}
		} else if b.epIndex == targetIndex {
			move := move{
				origin: squareIndex,
				target: targetIndex,
				kind:   EP_CAPTURE,
			}
			moves = append(moves, move)
		}
	}
	return moves
}

func (b board) generateKnightMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_KNIGHT_MOVES)
	for _, dir := range VERTICAL_DIRECTIONS {
		for _, moveDist := range KNIGHT_MOVES_DIST {
			targetIndex := SquareIndex(moveDist*dir) + squareIndex
			if b.squares[targetIndex] == EMPTY {
				move := move{
					origin: squareIndex,
					target: targetIndex,
					kind:   QUIET,
				}
				moves = append(moves, move)

			} else if b.colorBySquareIndex(targetIndex) == side^1 {
				move := move{
					origin: squareIndex,
					target: targetIndex,
					kind:   CAPTURE,
				}
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func (b board) generateBishopMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_BISHOP_MOVES)
	for _, dir := range VERTICAL_DIRECTIONS {
		for _, moveDist := range DIAGONAL_DISTS {
			for i := 1; i < BOARD_WIDTH_HEIGHT; i++ {
				targetIndex := SquareIndex(moveDist*dir*i) + squareIndex
				targetSquare := b.squares[targetIndex]
				if targetSquare == INVALID {
					break
				}
				squareColor := b.colorBySquareIndex(targetIndex)
				if squareColor == side {
					break
				} else if squareColor == side^1 {
					move := move{
						origin: squareIndex,
						target: targetIndex,
						kind:   CAPTURE,
					}
					moves = append(moves, move)
					break
				}
				move := move{
					origin: squareIndex,
					target: targetIndex,
					kind:   QUIET,
				}
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func (b board) generateRookMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_ROOK_MOVES)
	return moves
}
