package main

func (b board) squareKnightAttackers(attackersPtr *map[SquareIndex]bool, side Color, squareIndex SquareIndex) {
	attackers := *attackersPtr

	enemyKnight := BLACK_KNIGHT

	if side == BLACK {
		enemyKnight = WHITE_KNIGHT
	}

	for _, dir := range MOVE_DIRECTIONS {
		for _, moveDist := range KNIGHT_MOVE_DISTS {
			originIndex := SquareIndex(moveDist*dir) + squareIndex
			if b.squares[originIndex] == enemyKnight {
				attackers[originIndex] = true
			}
		}
	}
}

func (b board) squareDiagonalAttackers(attackersPtr *map[SquareIndex]bool, side Color, squareIndex SquareIndex) {
	attackers := *attackersPtr

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
					attackers[originIndex] = true
				case enemyPawn:
					if i == 1 && dir*-1 == pawnAttackDir {
						attackers[originIndex] = true
					}
				case enemyKing:
					if i == 1 {
						attackers[originIndex] = true
					}
				}
				// Reached invalid, friendly, or non-diagonal moving piece
				break
			}
		}
	}
}

func (b board) squareCardinalAttackers(attackersPtr *map[SquareIndex]bool, side Color, squareIndex SquareIndex) {
	attackers := *attackersPtr

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
					attackers[originIndex] = true
				case enemyKing:
					if i == 1 {
						attackers[originIndex] = true
					}
				}
				// Reached invalid, friendly, or non-cardinal moving piece
				break
			}
		}
	}
}

func (b board) squareAttackers(attackersPtr *map[SquareIndex]bool, side Color, squareIndex SquareIndex) {
	b.squareKnightAttackers(attackersPtr, side, squareIndex)
	b.squareDiagonalAttackers(attackersPtr, side, squareIndex)
	b.squareCardinalAttackers(attackersPtr, side, squareIndex)
}
