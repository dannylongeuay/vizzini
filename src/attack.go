package main

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
