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

func (b *Board) SetSquare(c Coord, sq Square) {
	b.squares[c] = sq
	b.HashSquare(sq, c)
	switch sq {
	case WHITE_PAWN:
		b.bbWP.SetBit(c)
	case WHITE_KNIGHT:
		b.bbWN.SetBit(c)
	case WHITE_BISHOP:
		b.bbWB.SetBit(c)
	case WHITE_ROOK:
		b.bbWR.SetBit(c)
	case WHITE_QUEEN:
		b.bbWQ.SetBit(c)
	case WHITE_KING:
		b.bbWK.SetBit(c)
	case BLACK_PAWN:
		b.bbBP.SetBit(c)
	case BLACK_KNIGHT:
		b.bbBN.SetBit(c)
	case BLACK_BISHOP:
		b.bbBB.SetBit(c)
	case BLACK_ROOK:
		b.bbBR.SetBit(c)
	case BLACK_QUEEN:
		b.bbBQ.SetBit(c)
	case BLACK_KING:
		b.bbBK.SetBit(c)
	default:
		panic(fmt.Errorf("set square(%v) error at coord %v", SQUARE_MAP[sq], COORD_MAP[c]))
	}
}

func (b *Board) ClearSquare(c Coord, sq Square) {
	// if sq != b.squares[c] {
	// 	panic(fmt.Errorf("clearing square mismatch %v != %v at coord %v is not allowed", SQUARE_MAP[sq], SQUARE_MAP[b.squares[c]], COORD_MAP[c]))
	// }

	b.squares[c] = EMPTY
	b.HashSquare(sq, c)
	switch sq {
	case WHITE_PAWN:
		b.bbWP.ClearBit(c)
	case WHITE_KNIGHT:
		b.bbWN.ClearBit(c)
	case WHITE_BISHOP:
		b.bbWB.ClearBit(c)
	case WHITE_ROOK:
		b.bbWR.ClearBit(c)
	case WHITE_QUEEN:
		b.bbWQ.ClearBit(c)
	case WHITE_KING:
		b.bbWK.ClearBit(c)
	case BLACK_PAWN:
		b.bbBP.ClearBit(c)
	case BLACK_KNIGHT:
		b.bbBN.ClearBit(c)
	case BLACK_BISHOP:
		b.bbBB.ClearBit(c)
	case BLACK_ROOK:
		b.bbBR.ClearBit(c)
	case BLACK_QUEEN:
		b.bbBQ.ClearBit(c)
	case BLACK_KING:
		b.bbBK.ClearBit(c)
	default:
		panic(fmt.Errorf("clear square(%v) error at coord %v", SQUARE_MAP[sq], COORD_MAP[c]))
	}
}

var CastleCheckIndexes = map[Coord]CastleRights{
	A8: 14, // KQk
	E8: 12, // KQ
	H8: 13, // KQq
	A1: 11, // Kkq
	E1: 3,  // kq
	H1: 7,  // Qkq
}

func (b *Board) UpdateCastleRights(c Coord) {
	castleBits, present := CastleCheckIndexes[c]
	if present {
		b.HashCastling()
		b.castleRights &= castleBits
		b.HashCastling()
	}
}

func (b *Board) MakeMove(m Move) error {
	mu := m.Unpack()

	u := Undo{
		move:         m,
		castleRights: b.castleRights,
		epCoord:      b.epCoord,
		halfMove:     b.halfMove,
		hash:         b.hash,
	}
	b.undos[b.undoIndex] = u
	b.undoIndex++

	b.ClearSquare(mu.originCoord, mu.originSquare)

	if b.castleRights != 0 {
		b.UpdateCastleRights(mu.originCoord)
		b.UpdateCastleRights(mu.dstCoord)
	}

	b.halfMove++

	switch mu.originSquare {
	case WHITE_PAWN:
		fallthrough
	case BLACK_PAWN:
		b.halfMove = 0
	case WHITE_KING:
		b.whiteKingCoord = mu.dstCoord
	case BLACK_KING:
		b.blackKingCoord = mu.dstCoord
	}

	if mu.dstSquare != EMPTY {
		b.ClearSquare(mu.dstCoord, mu.dstSquare)
		b.halfMove = 0
	}

	if b.epCoord != A1 {
		b.HashEnPassant()
	}
	b.epCoord = 0

	if mu.moveKind == DOUBLE_PAWN_PUSH {
		if b.sideToMove == WHITE &&
			(b.squares[int(mu.dstCoord)-SHIFT_HORIZONTAL] == BLACK_PAWN ||
				b.squares[int(mu.dstCoord)+SHIFT_HORIZONTAL] == BLACK_PAWN) {
			b.epCoord = mu.dstCoord - Coord(SHIFT_VERTICAL)
			b.HashEnPassant()
		} else if b.sideToMove == BLACK &&
			(b.squares[int(mu.dstCoord)-SHIFT_HORIZONTAL] == WHITE_PAWN ||
				b.squares[int(mu.dstCoord)+SHIFT_HORIZONTAL] == WHITE_PAWN) {
			b.epCoord = mu.dstCoord + Coord(SHIFT_VERTICAL)
			b.HashEnPassant()
		}
	}

	switch mu.moveKind {
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch mu.dstCoord {
		case G1:
			b.ClearSquare(H1, WHITE_ROOK)
			b.SetSquare(F1, WHITE_ROOK)
		case C1:
			b.ClearSquare(A1, WHITE_ROOK)
			b.SetSquare(D1, WHITE_ROOK)
		case G8:
			b.ClearSquare(H8, BLACK_ROOK)
			b.SetSquare(F8, BLACK_ROOK)
		case C8:
			b.ClearSquare(A8, BLACK_ROOK)
			b.SetSquare(D8, BLACK_ROOK)
		}
	case EP_CAPTURE:
		if b.sideToMove == WHITE {
			b.ClearSquare(mu.dstCoord-Coord(SHIFT_VERTICAL), BLACK_PAWN)
		} else {
			b.ClearSquare(mu.dstCoord+Coord(SHIFT_VERTICAL), WHITE_PAWN)
		}
	case KNIGHT_PROMOTION_CAPTURE:
		fallthrough
	case KNIGHT_PROMOTION:
		if b.sideToMove == WHITE {
			mu.originSquare = WHITE_KNIGHT
		} else {
			mu.originSquare = BLACK_KNIGHT
		}
	case BISHOP_PROMOTION_CAPTURE:
		fallthrough
	case BISHOP_PROMOTION:
		if b.sideToMove == WHITE {
			mu.originSquare = WHITE_BISHOP
		} else {
			mu.originSquare = BLACK_BISHOP
		}
	case ROOK_PROMOTION_CAPTURE:
		fallthrough
	case ROOK_PROMOTION:
		if b.sideToMove == WHITE {
			mu.originSquare = WHITE_ROOK
		} else {
			mu.originSquare = BLACK_ROOK
		}
	case QUEEN_PROMOTION_CAPTURE:
		fallthrough
	case QUEEN_PROMOTION:
		if b.sideToMove == WHITE {
			mu.originSquare = WHITE_QUEEN
		} else {
			mu.originSquare = BLACK_QUEEN
		}
	}

	b.SetSquare(mu.dstCoord, mu.originSquare)

	kingCoord := b.whiteKingCoord
	if b.sideToMove == BLACK {
		b.fullMove++
		kingCoord = b.blackKingCoord
	}

	kingAttacked := b.CoordAttacked(kingCoord, b.sideToMove)

	b.HashSide()
	b.sideToMove ^= 1

	if kingAttacked {
		return fmt.Errorf("king is attacked")
	}

	return nil
}

func (b *Board) UndoMove() error {
	if b.undoIndex <= 0 {
		return fmt.Errorf("invalid undo index: %v", b.undoIndex)
	}
	b.undoIndex--
	u := b.undos[b.undoIndex]
	mu := u.move.Unpack()

	b.sideToMove ^= 1

	if b.sideToMove == BLACK {
		b.fullMove--
	}

	b.ClearSquare(mu.dstCoord, mu.originSquare)

	if mu.dstSquare != EMPTY {
		b.SetSquare(mu.dstCoord, mu.dstSquare)
	}

	switch mu.moveKind {
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch mu.dstCoord {
		case G1:
			b.ClearSquare(F1, WHITE_ROOK)
			b.SetSquare(H1, WHITE_ROOK)
		case C1:
			b.ClearSquare(D1, WHITE_ROOK)
			b.SetSquare(A1, WHITE_ROOK)
		case G8:
			b.ClearSquare(F8, BLACK_ROOK)
			b.SetSquare(H8, BLACK_ROOK)
		case C8:
			b.ClearSquare(D8, BLACK_ROOK)
			b.SetSquare(A8, BLACK_ROOK)
		}
	case EP_CAPTURE:
		epCaptureCoord := mu.dstCoord - Coord(SHIFT_VERTICAL)
		epSquare := WHITE_PAWN
		if b.sideToMove == WHITE {
			epCaptureCoord = mu.dstCoord + Coord(SHIFT_VERTICAL)
			epSquare = BLACK_PAWN
		}
		b.SetSquare(epCaptureCoord, epSquare)
	}
	b.SetSquare(mu.originCoord, mu.originSquare)

	switch mu.dstSquare {
	case WHITE_KING:
		b.whiteKingCoord = mu.originCoord
	case BLACK_KING:
		b.blackKingCoord = mu.originCoord
	}

	b.castleRights = u.castleRights
	b.epCoord = u.epCoord
	b.halfMove = u.halfMove
	b.hash = u.hash
	return nil
}
