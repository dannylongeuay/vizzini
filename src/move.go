package main

import (
	"fmt"
)

var MOVE_KIND_MAP = [MOVE_KINDS]string{
	"QUIET", "DOUBLE_PAWN_PUSH", "KING_CASTLE",
	"QUEEN_CASTLE", "CAPTURE", "EP_CAPTURE",
	"KNIGHT_PROMOTION", "BISHOP_PROMOTION",
	"ROOK_PROMOTION", "QUEEN_PROMOTION",
	"KNIGHT_PROMOTION_CAPTURE", "BISHOP_PROMOTION_CAPTURE",
	"ROOK_PROMOTION_CAPTURE", "QUEEN_PROMOTION_CAPTURE",
}

func AppendMove(moves *[]Move, originCoord Coord, dstCoord Coord, originSquare Square, dstSquare Square, moveKind MoveKind) {
	moveOriginCoord := Move(originCoord) << MOVE_ORIGIN_COORD_SHIFT
	moveDstCoord := Move(dstCoord) << MOVE_DST_COORD_SHIFT
	moveOriginSquare := Move(originSquare) << MOVE_ORIGIN_SQUARE_SHIFT
	moveDstSquare := Move(dstSquare) << MOVE_DST_SQUARE_SHIFT
	move := moveOriginCoord | moveDstCoord | moveOriginSquare | moveDstSquare | Move(moveKind)
	*moves = append(*moves, move)
}

func NewMove(originCoord Coord, dstCoord Coord, originSquare Square, dstSquare Square, moveKind MoveKind) Move {
	moveOriginCoord := Move(originCoord) << MOVE_ORIGIN_COORD_SHIFT
	moveDstCoord := Move(dstCoord) << MOVE_DST_COORD_SHIFT
	moveOriginSquare := Move(originSquare) << MOVE_ORIGIN_SQUARE_SHIFT
	moveDstSquare := Move(dstSquare) << MOVE_DST_SQUARE_SHIFT
	return moveOriginCoord | moveDstCoord | moveOriginSquare | moveDstSquare | Move(moveKind)
}

type MoveUnpacked struct {
	originCoord  Coord
	dstCoord     Coord
	originSquare Square
	dstSquare    Square
	moveKind     MoveKind
}

func (m *Move) ToString() string {
	mu := m.Unpack()
	mk := MOVE_KIND_MAP[mu.moveKind]
	oc := COORD_MAP[mu.originCoord]
	os := SQUARE_MAP[mu.originSquare]
	dc := COORD_MAP[mu.dstCoord]
	ds := SQUARE_MAP[mu.dstSquare]
	s := fmt.Sprintf("%v: %v(%v) -> %v(%v)", mk, oc, os, dc, ds)
	return s
}

func (m *Move) Unpack() MoveUnpacked {
	mu := MoveUnpacked{
		originCoord:  Coord((*m & MOVE_ORIGIN_COORD_MASK) >> MOVE_ORIGIN_COORD_SHIFT),
		dstCoord:     Coord((*m & MOVE_DST_COORD_MASK) >> MOVE_DST_COORD_SHIFT),
		originSquare: Square((*m & MOVE_ORIGIN_SQUARE_MASK) >> MOVE_ORIGIN_SQUARE_SHIFT),
		dstSquare:    Square((*m & MOVE_DST_SQUARE_MASK) >> MOVE_DST_SQUARE_SHIFT),
		moveKind:     MoveKind(*m & MOVE_KIND_MASK),
	}
	return mu
}

func (b *Board) MakeMove(m Move) error {
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
