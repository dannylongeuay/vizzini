package main

import (
	"fmt"
)

func (b *Board) generateMoves(side Color) []Move {
	moves := make([]Move, 0, MAX_GENERATED_MOVES)
	return moves
}

func (b *Board) generatePawnMoves(moves *[]Move, side Color) {
	// Handle quiet move
	// Handle double pawn push
	// Handle captures
	// Handle enPassant
}

func (b *Board) generateKnightMoves(moves *[]Move, side Color) {
}

func (b *Board) generateBishopMoves(moves *[]Move, side Color) {
}

func (b *Board) generateRookMoves(moves *[]Move, side Color) {
}

func (b *Board) generateQueenMoves(moves *[]Move, side Color) {
}

func (b *Board) generateKingMoves(moves *[]Move, side Color) {
}

func (b *Board) makeMove(m Move) error {
	originSquare := b.squares[A1]
	targetSquare := b.squares[A2]

	u := Undo{
		mv:             m,
		capturedSquare: targetSquare,
		castleRights:   b.castleRights,
		epCoord:        b.epCoord,
		halfMove:       b.halfMove,
		hash:           b.hash,
	}
	b.undos[b.undoIndex] = u
	b.undoIndex++

	// b.clearSquare(originSquare, m.origin)

	if b.castleRights != 0 {
		// b.updateCastleRights(m.origin)
		// b.updateCastleRights(m.target)
	}

	b.halfMove++

	switch originSquare {
	case WHITE_PAWN:
		fallthrough
	case BLACK_PAWN:
		b.halfMove = 0
	case WHITE_KING:
		// b.whiteKingIndex = m.target
	case BLACK_KING:
		// b.blackKingIndex = m.target
	}

	if targetSquare != EMPTY {
		// b.clearSquare(targetSquare, m.target)
		b.halfMove = 0
	}

	// if b.epIndex != 0 {
	// 	b.hashEnPassant()
	// }
	// b.epIndex = 0

	// if m.kind == DOUBLE_PAWN_PUSH {
	// 	if b.sideToMove == WHITE &&
	// 		(b.squares[m.target-SquareIndex(HORIZONTAL_MOVE_DIST)] == BLACK_PAWN ||
	// 			b.squares[m.target+SquareIndex(HORIZONTAL_MOVE_DIST)] == BLACK_PAWN) {
	// 		b.epIndex = m.target - SquareIndex(VERTICAL_MOVE_DIST)
	// 		b.hashEnPassant()
	// 	} else if b.sideToMove == BLACK &&
	// 		(b.squares[m.target-SquareIndex(HORIZONTAL_MOVE_DIST)] == WHITE_PAWN ||
	// 			b.squares[m.target+SquareIndex(HORIZONTAL_MOVE_DIST)] == WHITE_PAWN) {
	// 		b.epIndex = m.target - SquareIndex(VERTICAL_MOVE_DIST)
	// 		b.hashEnPassant()
	// 	}
	// }

	// switch m.kind {
	// case KING_CASTLE:
	// 	fallthrough
	// case QUEEN_CASTLE:
	// 	switch m.target {
	// 	case 97:
	// 		b.clearSquare(WHITE_ROOK, 98)
	// 		b.setSquare(WHITE_ROOK, 96)
	// 	case 93:
	// 		b.clearSquare(WHITE_ROOK, 91)
	// 		b.setSquare(WHITE_ROOK, 94)
	// 	case 27:
	// 		b.clearSquare(BLACK_ROOK, 28)
	// 		b.setSquare(BLACK_ROOK, 26)
	// 	case 23:
	// 		b.clearSquare(BLACK_ROOK, 21)
	// 		b.setSquare(BLACK_ROOK, 24)
	// 	}
	// case EP_CAPTURE:
	// 	if b.sideToMove == WHITE {
	// 		b.clearSquare(BLACK_PAWN, m.target+SquareIndex(VERTICAL_MOVE_DIST))
	// 	} else {
	// 		b.clearSquare(WHITE_PAWN, m.target-SquareIndex(VERTICAL_MOVE_DIST))
	// 	}
	// case KNIGHT_PROMOTION_CAPTURE:
	// 	fallthrough
	// case KNIGHT_PROMOTION:
	// 	if b.sideToMove == WHITE {
	// 		originSquare = WHITE_KNIGHT
	// 	} else {
	// 		originSquare = BLACK_KNIGHT
	// 	}
	// case BISHOP_PROMOTION_CAPTURE:
	// 	fallthrough
	// case BISHOP_PROMOTION:
	// 	if b.sideToMove == WHITE {
	// 		originSquare = WHITE_BISHOP
	// 	} else {
	// 		originSquare = BLACK_BISHOP
	// 	}
	// case ROOK_PROMOTION_CAPTURE:
	// 	fallthrough
	// case ROOK_PROMOTION:
	// 	if b.sideToMove == WHITE {
	// 		originSquare = WHITE_ROOK
	// 	} else {
	// 		originSquare = BLACK_ROOK
	// 	}
	// case QUEEN_PROMOTION_CAPTURE:
	// 	fallthrough
	// case QUEEN_PROMOTION:
	// 	if b.sideToMove == WHITE {
	// 		originSquare = WHITE_QUEEN
	// 	} else {
	// 		originSquare = BLACK_QUEEN
	// 	}
	// }

	// b.setSquare(originSquare, m.target)

	// kingIndex := b.whiteKingIndex
	// if b.sideToMove == BLACK {
	// 	b.fullMove++
	// 	kingIndex = b.blackKingIndex
	// }

	// kingAttackers := make(map[SquareIndex]bool, MAX_SQUARE_ATTACKERS)
	// b.squareAttackers(&kingAttackers, b.sideToMove, kingIndex)

	// b.hashSide()
	// b.sideToMove ^= 1

	// if len(kingAttackers) > 0 {
	// 	return fmt.Errorf("king has %v attacker(s)", len(kingAttackers))
	// }

	return nil
}

func (b *Board) undoMove() error {
	if b.undoIndex <= 0 {
		return fmt.Errorf("invalid undo index: %v", b.undoIndex)
	}
	b.undoIndex--
	u := b.undos[b.undoIndex]

	b.sideToMove ^= 1
	// targetSquare := b.squares[u.mv.target]

	if b.sideToMove == BLACK {
		b.fullMove--
	}

	// b.clearSquare(targetSquare, u.mv.target)

	if u.capturedSquare != EMPTY {
		// b.setSquare(u.capturedSquare, u.mv.target)
	}

	// switch u.mv.kind {
	// case KING_CASTLE:
	// 	fallthrough
	// case QUEEN_CASTLE:
	// 	switch u.mv.target {
	// 	case 97:
	// 		b.clearSquare(WHITE_ROOK, 96)
	// 		b.setSquare(WHITE_ROOK, 98)
	// 	case 93:
	// 		b.clearSquare(WHITE_ROOK, 94)
	// 		b.setSquare(WHITE_ROOK, 91)
	// 	case 27:
	// 		b.clearSquare(BLACK_ROOK, 26)
	// 		b.setSquare(BLACK_ROOK, 28)
	// 	case 23:
	// 		b.clearSquare(BLACK_ROOK, 24)
	// 		b.setSquare(BLACK_ROOK, 21)
	// 	}
	// case EP_CAPTURE:
	// 	epCaptureIndex := u.mv.target - SquareIndex(VERTICAL_MOVE_DIST)
	// 	epSquare := WHITE_PAWN
	// 	if b.sideToMove == WHITE {
	// 		epCaptureIndex = u.mv.target + SquareIndex(VERTICAL_MOVE_DIST)
	// 		epSquare = BLACK_PAWN
	// 	}
	// 	b.setSquare(epSquare, epCaptureIndex)
	// case KNIGHT_PROMOTION_CAPTURE:
	// 	fallthrough
	// case KNIGHT_PROMOTION:
	// 	fallthrough
	// case BISHOP_PROMOTION_CAPTURE:
	// 	fallthrough
	// case BISHOP_PROMOTION:
	// 	fallthrough
	// case ROOK_PROMOTION_CAPTURE:
	// 	fallthrough
	// case ROOK_PROMOTION:
	// 	fallthrough
	// case QUEEN_PROMOTION_CAPTURE:
	// 	fallthrough
	// case QUEEN_PROMOTION:
	// 	if b.sideToMove == WHITE {
	// 		targetSquare = WHITE_PAWN
	// 	} else {
	// 		targetSquare = BLACK_PAWN
	// 	}
	// }
	// b.setSquare(targetSquare, u.mv.origin)

	// switch targetSquare {
	// case WHITE_KING:
	// 	b.whiteKingIndex = u.mv.origin
	// case BLACK_KING:
	// 	b.blackKingIndex = u.mv.origin
	// }

	// b.castleRights = u.castleRights
	b.epCoord = u.epCoord
	b.halfMove = u.halfMove
	b.hash = u.hash
	return nil
}
