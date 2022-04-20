package main

type move struct {
	origin SquareIndex
	target SquareIndex
	kind   MoveKind
}

var KNIGHT_MOVE_DISTS = []int{8, 12, 19, 21}

var MOVE_DIRECTIONS = []int{POSITIVE_DIR, NEGATIVE_DIR}

var DIAGONAL_MOVE_DISTS = []int{POS_DIAG_MOVE_DIST, NEG_DIAG_MOVE_DIST}
var CARDINAL_MOVE_DISTS = []int{HORIZONTAL_MOVE_DIST, VERTICAL_MOVE_DIST}

func (b board) generateMoves(side Color) []move {
	moves := make([]move, 0, MAX_GENERATED_MOVES)
	for square, squareIndexes := range b.pieceSets {
		squareColor := colorBySquare(square)
		if squareColor != side {
			continue
		}
		for squareIndex := range squareIndexes {
			switch square {
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
				rookMoves := b.generateRookMoves(side, squareIndex)
				moves = append(moves, rookMoves...)
				break
			case WHITE_QUEEN:
				fallthrough
			case BLACK_QUEEN:
				queenMoves := b.generateQueenMoves(side, squareIndex)
				moves = append(moves, queenMoves...)
				break
			case WHITE_KING:
				fallthrough
			case BLACK_KING:
				kingMoves := b.generateKingMoves(side, squareIndex)
				moves = append(moves, kingMoves...)
				break
			}
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
func (b board) generateBishopMoves(side Color, squareIndex SquareIndex) []move {
	return b.generateSlidingMoves(side, squareIndex, DIAGONAL_MOVE_DISTS, MAX_BISHOP_MOVES, MAX_MOVE_RANGE)
}

func (b board) generateRookMoves(side Color, squareIndex SquareIndex) []move {
	return b.generateSlidingMoves(side, squareIndex, CARDINAL_MOVE_DISTS, MAX_ROOK_MOVES, MAX_MOVE_RANGE)
}

func (b board) generateQueenMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_QUEEN_MOVES)
	diagonalMoves := b.generateSlidingMoves(side, squareIndex, DIAGONAL_MOVE_DISTS, MAX_BISHOP_MOVES, MAX_MOVE_RANGE)
	moves = append(moves, diagonalMoves...)
	cardinalMoves := b.generateSlidingMoves(side, squareIndex, CARDINAL_MOVE_DISTS, MAX_ROOK_MOVES, MAX_MOVE_RANGE)
	moves = append(moves, cardinalMoves...)
	return moves
}

func (b board) generateKingMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_KING_MOVES)
	diagonalMoves := b.generateSlidingMoves(side, squareIndex, DIAGONAL_MOVE_DISTS, MAX_KING_MOVES/2, KING_MOVE_RANGE)
	moves = append(moves, diagonalMoves...)
	cardinalMoves := b.generateSlidingMoves(side, squareIndex, CARDINAL_MOVE_DISTS, MAX_KING_MOVES/2, KING_MOVE_RANGE)
	moves = append(moves, cardinalMoves...)
	kingsideCastleAvail := false
	queensideCastleAvail := false
	if side == WHITE {
		if b.castleRights&(1<<3) == 1<<3 {
			kingsideCastleAvail = true
		}
		if b.castleRights&(1<<2) == 1<<2 {
			queensideCastleAvail = true
		}
	} else {
		if b.castleRights&(1<<1) == 1<<1 {
			kingsideCastleAvail = true
		}
		if b.castleRights&1 == 1 {
			queensideCastleAvail = true
		}
	}
	if kingsideCastleAvail && b.canCastle(side, POSITIVE_DIR, squareIndex) {
		targetIndex := KING_CASTLE_MOVE_DIST + squareIndex
		move := move{
			origin: squareIndex,
			target: targetIndex,
			kind:   KING_CASTLE,
		}
		moves = append(moves, move)
	}
	if queensideCastleAvail && b.canCastle(side, NEGATIVE_DIR, squareIndex) {
		targetIndex := SquareIndex(KING_CASTLE_MOVE_DIST*NEGATIVE_DIR + int(squareIndex))
		move := move{
			origin: squareIndex,
			target: targetIndex,
			kind:   QUEEN_CASTLE,
		}
		moves = append(moves, move)
	}
	return moves
}

func (b board) canCastle(side Color, dir int, squareIndex SquareIndex) bool {
	rookDist := 3
	if dir == NEGATIVE_DIR {
		rookDist = 4
	}
	for i := 0; i <= rookDist; i++ {
		checkIndex := SquareIndex(i*dir) + squareIndex
		// Check that squares are empty between king and rook
		if i > 0 && i < rookDist && b.squares[checkIndex] != EMPTY {
			return false
		}
		// Check that king is not in check and travel squares are not attacked
		if i < 3 {
			attackers := b.squareAttackers(side, checkIndex)
			if len(attackers) > 0 {
				return false
			}
		}
	}
	return true
}
