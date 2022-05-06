package main

import (
	"fmt"
	"sort"
)

type move struct {
	origin SquareIndex
	target SquareIndex
	kind   MoveKind
}

var KNIGHT_MOVE_DISTS = []int{8, 12, 19, 21}

var MOVE_DIRECTIONS = []int{POSITIVE_DIR, NEGATIVE_DIR}

var DIAGONAL_MOVE_DISTS = []int{POS_DIAG_MOVE_DIST, NEG_DIAG_MOVE_DIST}
var CARDINAL_MOVE_DISTS = []int{HORIZONTAL_MOVE_DIST, VERTICAL_MOVE_DIST}

func (b *board) generateMoves(side Color) []move {
	moves := make([]move, 0, MAX_GENERATED_MOVES)
	for square, squareIndexes := range b.pieceSets {
		squareColor := colorBySquare(square)
		if squareColor != side {
			continue
		}
		var sortedKeys []int
		for key := range squareIndexes {
			sortedKeys = append(sortedKeys, int(key))
		}

		sort.Ints(sortedKeys)

		for _, key := range sortedKeys {
			squareIndex := SquareIndex(key)
			switch square {
			case WHITE_PAWN:
				fallthrough
			case BLACK_PAWN:
				pawnMoves := b.generatePawnMoves(side, squareIndex)
				moves = append(moves, pawnMoves...)
			case WHITE_KNIGHT:
				fallthrough
			case BLACK_KNIGHT:
				knightMoves := b.generateKnightMoves(side, squareIndex)
				moves = append(moves, knightMoves...)
			case WHITE_BISHOP:
				fallthrough
			case BLACK_BISHOP:
				bishopMoves := b.generateBishopMoves(side, squareIndex)
				moves = append(moves, bishopMoves...)
			case WHITE_ROOK:
				fallthrough
			case BLACK_ROOK:
				rookMoves := b.generateRookMoves(side, squareIndex)
				moves = append(moves, rookMoves...)
			case WHITE_QUEEN:
				fallthrough
			case BLACK_QUEEN:
				queenMoves := b.generateQueenMoves(side, squareIndex)
				moves = append(moves, queenMoves...)
			case WHITE_KING:
				fallthrough
			case BLACK_KING:
				kingMoves := b.generateKingMoves(side, squareIndex)
				moves = append(moves, kingMoves...)
			}
		}
	}
	return moves
}

func (b *board) generateLegalMoves(moves []move, maxSize int) []move {
	legalMoves := make([]move, 0, maxSize)
	for _, m := range moves {
		err := b.makeMove(m)
		b.undoMove()
		if err != nil {
			continue
		}
		legalMoves = append(legalMoves, m)
	}
	return legalMoves
}

func (b *board) generatePawnMoves(side Color, squareIndex SquareIndex) []move {
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
		jumpedIndex := SquareIndex(VERTICAL_MOVE_DIST*dir) + squareIndex
		if b.squares[targetIndex] == EMPTY && b.squares[jumpedIndex] == EMPTY {
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
	return b.generateLegalMoves(moves, MAX_PAWN_MOVES)
}

func (b *board) generateKnightMoves(side Color, squareIndex SquareIndex) []move {
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
	return b.generateLegalMoves(moves, MAX_KNIGHT_MOVES)
}

func (b *board) generateSlidingMoves(side Color, squareIndex SquareIndex, moveDists []int, maxMoves int, moveRange int) []move {
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
	return b.generateLegalMoves(moves, maxMoves)
}
func (b *board) generateBishopMoves(side Color, squareIndex SquareIndex) []move {
	return b.generateSlidingMoves(side, squareIndex, DIAGONAL_MOVE_DISTS, MAX_BISHOP_MOVES, MAX_MOVE_RANGE)
}

func (b *board) generateRookMoves(side Color, squareIndex SquareIndex) []move {
	return b.generateSlidingMoves(side, squareIndex, CARDINAL_MOVE_DISTS, MAX_ROOK_MOVES, MAX_MOVE_RANGE)
}

func (b *board) generateQueenMoves(side Color, squareIndex SquareIndex) []move {
	moves := make([]move, 0, MAX_QUEEN_MOVES)
	diagonalMoves := b.generateSlidingMoves(side, squareIndex, DIAGONAL_MOVE_DISTS, MAX_BISHOP_MOVES, MAX_MOVE_RANGE)
	moves = append(moves, diagonalMoves...)
	cardinalMoves := b.generateSlidingMoves(side, squareIndex, CARDINAL_MOVE_DISTS, MAX_ROOK_MOVES, MAX_MOVE_RANGE)
	moves = append(moves, cardinalMoves...)
	return moves
}

func (b *board) generateKingMoves(side Color, squareIndex SquareIndex) []move {
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

func (b *board) canCastle(side Color, dir int, squareIndex SquareIndex) bool {
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

func (b *board) clearSquare(square Square, squareIndex SquareIndex) error {
	if square == EMPTY || square == INVALID || square != b.squares[squareIndex] {
		return fmt.Errorf("clearing square %v at index %v is not allowed", square, squareIndex)
	}
	b.squares[squareIndex] = EMPTY
	delete(b.pieceSets[square], squareIndex)
	b.hashSquare(square, squareIndex)
	return nil
}

func (b *board) setSquare(square Square, squareIndex SquareIndex) error {
	if square == EMPTY || square == INVALID {
		return fmt.Errorf("setting square %v at index %v is not allowed", square, squareIndex)
	}
	b.squares[squareIndex] = square
	b.pieceSets[square][squareIndex] = true
	b.hashSquare(square, squareIndex)
	return nil
}

var CastleCheckIndexes = map[SquareIndex]CastleRights{
	21: 14, // a8 : KQk
	25: 12, // e8 : KQ
	28: 13, // h8 : KQq
	91: 11, // a1 : Kkq
	95: 3,  // e1 : kq
	98: 7,  // h1 : Qkq
}

func (b *board) updateCastleRights(squareIndex SquareIndex) {
	castleBits, present := CastleCheckIndexes[squareIndex]
	if present {
		b.hashCastling()
		b.castleRights &= castleBits
		b.hashCastling()
	}
}

func (b *board) makeMove(m move) error {
	originSquare := b.squares[m.origin]
	targetSquare := b.squares[m.target]

	u := undo{
		mv:             m,
		capturedSquare: targetSquare,
		castleRights:   b.castleRights,
		epIndex:        b.epIndex,
		halfMove:       b.halfMove,
		hash:           b.hash,
	}
	b.undos[b.undoIndex] = u
	b.undoIndex++

	b.clearSquare(originSquare, m.origin)

	if b.castleRights != 0 {
		b.updateCastleRights(m.origin)
		b.updateCastleRights(m.target)
	}

	b.halfMove++

	switch originSquare {
	case WHITE_PAWN:
		fallthrough
	case BLACK_PAWN:
		b.halfMove = 0
	case WHITE_KING:
		b.whiteKingIndex = m.target
	case BLACK_KING:
		b.blackKingIndex = m.target
	}

	if targetSquare != EMPTY {
		b.clearSquare(targetSquare, m.target)
		b.halfMove = 0
	}

	switch m.kind {
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch m.target {
		case 97:
			b.clearSquare(WHITE_ROOK, 98)
			b.setSquare(WHITE_ROOK, 96)
		case 93:
			b.clearSquare(WHITE_ROOK, 91)
			b.setSquare(WHITE_ROOK, 94)
		case 27:
			b.clearSquare(BLACK_ROOK, 28)
			b.setSquare(BLACK_ROOK, 26)
		case 23:
			b.clearSquare(BLACK_ROOK, 21)
			b.setSquare(BLACK_ROOK, 24)
		}
	case EP_CAPTURE:
		if b.sideToMove == WHITE {
			b.clearSquare(BLACK_PAWN, m.target+SquareIndex(VERTICAL_MOVE_DIST))
		} else {
			b.clearSquare(WHITE_PAWN, m.target-SquareIndex(VERTICAL_MOVE_DIST))
		}
		b.hashEnPassant()
		b.epIndex = 0
	case KNIGHT_PROMOTION_CAPTURE:
		fallthrough
	case KNIGHT_PROMOTION:
		if b.sideToMove == WHITE {
			originSquare = WHITE_KNIGHT
		} else {
			originSquare = BLACK_KNIGHT
		}
	case BISHOP_PROMOTION_CAPTURE:
		fallthrough
	case BISHOP_PROMOTION:
		if b.sideToMove == WHITE {
			originSquare = WHITE_BISHOP
		} else {
			originSquare = BLACK_BISHOP
		}
	case ROOK_PROMOTION_CAPTURE:
		fallthrough
	case ROOK_PROMOTION:
		if b.sideToMove == WHITE {
			originSquare = WHITE_ROOK
		} else {
			originSquare = BLACK_ROOK
		}
	case QUEEN_PROMOTION_CAPTURE:
		fallthrough
	case QUEEN_PROMOTION:
		if b.sideToMove == WHITE {
			originSquare = WHITE_QUEEN
		} else {
			originSquare = BLACK_QUEEN
		}
	}

	b.setSquare(originSquare, m.target)

	kingIndex := b.whiteKingIndex
	if b.sideToMove == BLACK {
		b.fullMove++
		kingIndex = b.blackKingIndex
	}

	kingAttackers := b.squareAttackers(b.sideToMove, kingIndex)

	b.hashSide()
	b.sideToMove ^= 1

	if len(kingAttackers) > 0 {
		return fmt.Errorf("king has %v attacker(s)", len(kingAttackers))
	}

	return nil
}

func (b *board) undoMove() error {
	if b.undoIndex <= 0 {
		return fmt.Errorf("invalid undo index: %v", b.undoIndex)
	}
	b.undoIndex--
	u := b.undos[b.undoIndex]

	b.sideToMove ^= 1
	targetSquare := b.squares[u.mv.target]

	if b.sideToMove == BLACK {
		b.fullMove--
	}

	b.clearSquare(targetSquare, u.mv.target)

	if u.capturedSquare != EMPTY {
		b.setSquare(u.capturedSquare, u.mv.target)
	}

	switch u.mv.kind {
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch u.mv.target {
		case 97:
			b.clearSquare(WHITE_ROOK, 96)
			b.setSquare(WHITE_ROOK, 98)
		case 93:
			b.clearSquare(WHITE_ROOK, 94)
			b.setSquare(WHITE_ROOK, 91)
		case 27:
			b.clearSquare(BLACK_ROOK, 26)
			b.setSquare(BLACK_ROOK, 28)
		case 23:
			b.clearSquare(BLACK_ROOK, 24)
			b.setSquare(BLACK_ROOK, 21)
		}
	case EP_CAPTURE:
		epCaptureIndex := u.mv.target - SquareIndex(VERTICAL_MOVE_DIST)
		epSquare := WHITE_PAWN
		if b.sideToMove == WHITE {
			epCaptureIndex = u.mv.target + SquareIndex(VERTICAL_MOVE_DIST)
			epSquare = BLACK_PAWN
		}
		b.setSquare(epSquare, epCaptureIndex)
	case KNIGHT_PROMOTION_CAPTURE:
		fallthrough
	case KNIGHT_PROMOTION:
		fallthrough
	case BISHOP_PROMOTION_CAPTURE:
		fallthrough
	case BISHOP_PROMOTION:
		fallthrough
	case ROOK_PROMOTION_CAPTURE:
		fallthrough
	case ROOK_PROMOTION:
		fallthrough
	case QUEEN_PROMOTION_CAPTURE:
		fallthrough
	case QUEEN_PROMOTION:
		if b.sideToMove == WHITE {
			targetSquare = WHITE_PAWN
		} else {
			targetSquare = BLACK_PAWN
		}
	}
	b.setSquare(targetSquare, u.mv.origin)

	switch targetSquare {
	case WHITE_KING:
		b.whiteKingIndex = u.mv.origin
	case BLACK_KING:
		b.blackKingIndex = u.mv.origin
	}

	b.castleRights = u.castleRights
	b.epIndex = u.epIndex
	b.halfMove = u.halfMove
	b.hash = u.hash
	return nil
}
