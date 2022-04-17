package main

type move struct {
	origin SquareIndex
	target SquareIndex
	kind   MoveKind
}

var DIAGONAL_MOVE_DISTS = []int{POS_DIAG_MOVE_DIST, NEG_DIAG_MOVE_DIST}
var KNIGHT_MOVE_DISTS = []int{8, 12, 19, 21}
var CARDINAL_MOVE_DISTS = []int{HORIZONTAL_MOVE_DIST, VERTICAL_MOVE_DIST}
var CARD_DIAG_MOVE_DISTS = []int{
	HORIZONTAL_MOVE_DIST,
	VERTICAL_MOVE_DIST,
	POS_DIAG_MOVE_DIST,
	NEG_DIAG_MOVE_DIST,
}
var MOVE_DIRECTIONS = []int{POSITIVE_DIR, NEGATIVE_DIR}

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
			bishopMoves := b.generateSlidingMoves(side, squareIndex, DIAGONAL_MOVE_DISTS, MAX_BISHOP_MOVES, MAX_MOVE_RANGE)
			moves = append(moves, bishopMoves...)
			break
		case WHITE_ROOK:
			fallthrough
		case BLACK_ROOK:
			rookMoves := b.generateSlidingMoves(side, squareIndex, CARDINAL_MOVE_DISTS, MAX_ROOK_MOVES, MAX_MOVE_RANGE)
			moves = append(moves, rookMoves...)
			break
		case WHITE_QUEEN:
			fallthrough
		case BLACK_QUEEN:
			queenMoves := b.generateSlidingMoves(side, squareIndex, CARD_DIAG_MOVE_DISTS, MAX_QUEEN_MOVES, MAX_MOVE_RANGE)
			moves = append(moves, queenMoves...)
			break
		case WHITE_KING:
			fallthrough
		case BLACK_KING:
			// TODO: add castling moves
			kingMoves := b.generateSlidingMoves(side, squareIndex, CARD_DIAG_MOVE_DISTS, MAX_KING_MOVES, KING_MOVE_RANGE)
			moves = append(moves, kingMoves...)
			break
		}
	}
	return moves
}

func (b board) generatePawnMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_PAWN_MOVES)
	dir := POSITIVE_DIR
	rank := rankBySquareIndex(squareIndex)
	doublePushRank := RANK_SEVEN
	promotionRank := RANK_TWO

	if side == WHITE {
		dir = NEGATIVE_DIR
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
	for _, moveDist := range DIAGONAL_MOVE_DISTS {
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
	for _, dir := range MOVE_DIRECTIONS {
		for _, moveDist := range KNIGHT_MOVE_DISTS {
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

func (b board) generateSlidingMoves(side Color, squareIndex SquareIndex, moveDists []int, maxMoves int, moveRange int) []move {
	moves := make([]move, 0, maxMoves)
	for _, dir := range MOVE_DIRECTIONS {
		for _, moveDist := range moveDists {
			for i := 1; i < moveRange; i++ {
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
