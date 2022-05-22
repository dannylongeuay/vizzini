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

func NewUndo(move Move, clearSquare Square, halfMove HalfMove, castleRights CastleRights, epCoord Coord) Undo {
	undoClearSquare := Undo(clearSquare) << UNDO_CLEAR_SQUARE_SHIFT
	undoHalfMove := Undo(halfMove) << UNDO_HALF_MOVE_SHIFT
	// hm := (undoHalfMove & UNDO_HALF_MOVE_MASK) >> UNDO_HALF_MOVE_SHIFT
	// fmt.Println(hm)
	undoCastleRights := Undo(castleRights) << UNDO_CASTLE_RIGHTS_SHIFT
	undoEpCoord := Undo(epCoord) << UNDO_EP_COORD_SHIFT
	return Undo(move) | undoClearSquare | undoHalfMove | undoCastleRights | undoEpCoord
}

type MoveUnpacked struct {
	originCoord  Coord
	dstCoord     Coord
	originSquare Square
	dstSquare    Square
	moveKind     MoveKind
}

type UndoUnpacked struct {
	MoveUnpacked
	clearSquare  Square
	halfMove     HalfMove
	castleRights CastleRights
	epCoord      Coord
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

func (u *Undo) ToString() string {
	move := Move(*u & UNDO_MOVE_MASK)
	uu := u.Unpack()
	cs := SQUARE_MAP[uu.clearSquare]
	hm := uu.halfMove
	cr := uu.castleRights
	ec := COORD_MAP[uu.epCoord]
	s := fmt.Sprintf("ClearSquare: %v | HalfMove: %v | CastleRights: %v | EP Coord: %v | Move: %v", cs, hm, cr, ec, move.ToString())
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

func (u *Undo) Unpack() UndoUnpacked {
	move := Move(*u & UNDO_MOVE_MASK)
	uu := UndoUnpacked{
		MoveUnpacked: move.Unpack(),
		clearSquare:  Square((*u & UNDO_CLEAR_SQUARE_MASK) >> UNDO_CLEAR_SQUARE_SHIFT),
		halfMove:     HalfMove((*u & UNDO_HALF_MOVE_MASK) >> UNDO_HALF_MOVE_SHIFT),
		castleRights: CastleRights((*u & UNDO_CASTLE_RIGHTS_MASK) >> UNDO_CASTLE_RIGHTS_SHIFT),
		epCoord:      Coord((*u & UNDO_EP_COORD_MASK) >> UNDO_EP_COORD_SHIFT),
	}
	return uu
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
		b.kingCoords[b.sideToMove] = c
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
		b.kingCoords[b.sideToMove] = c
	default:
		panic(fmt.Errorf("set square(%v) error at coord %v", SQUARE_MAP[sq], COORD_MAP[c]))
	}
}

func (b *Board) ClearSquare(c Coord, sq Square) {
	if sq != b.squares[c] {
		panic(fmt.Errorf("clearing square mismatch %v != %v at coord %v is not allowed", SQUARE_MAP[sq], SQUARE_MAP[b.squares[c]], COORD_MAP[c]))
	}

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

	moveDstSetSquare := mu.originSquare

	epCoord := A1
	epCaptureSquare := EMPTY
	epCaptureCoord := A1

	castleRookSquare := EMPTY
	castleRookClearCoord := A1
	castleRookSetCoord := A1

	switch mu.moveKind {
	case DOUBLE_PAWN_PUSH:
		if b.sideToMove == WHITE &&
			(b.squares[int(mu.dstCoord)-SHIFT_HORIZONTAL] == BLACK_PAWN ||
				b.squares[int(mu.dstCoord)+SHIFT_HORIZONTAL] == BLACK_PAWN) {
			epCoord = mu.dstCoord - Coord(SHIFT_VERTICAL)
		} else if b.sideToMove == BLACK &&
			(b.squares[int(mu.dstCoord)-SHIFT_HORIZONTAL] == WHITE_PAWN ||
				b.squares[int(mu.dstCoord)+SHIFT_HORIZONTAL] == WHITE_PAWN) {
			epCoord = mu.dstCoord + Coord(SHIFT_VERTICAL)
		}
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch mu.dstCoord {
		case G1:
			castleRookSquare = WHITE_ROOK
			castleRookClearCoord = H1
			castleRookSetCoord = F1
		case C1:
			castleRookSquare = WHITE_ROOK
			castleRookClearCoord = A1
			castleRookSetCoord = D1
		case G8:
			castleRookSquare = BLACK_ROOK
			castleRookClearCoord = H8
			castleRookSetCoord = F8
		case C8:
			castleRookSquare = BLACK_ROOK
			castleRookClearCoord = A8
			castleRookSetCoord = D8
		}
	case EP_CAPTURE:
		if b.sideToMove == WHITE {
			epCaptureSquare = BLACK_PAWN
			epCaptureCoord = mu.dstCoord - Coord(SHIFT_VERTICAL)
		} else {
			epCaptureSquare = WHITE_PAWN
			epCaptureCoord = mu.dstCoord + Coord(SHIFT_VERTICAL)
		}
	case KNIGHT_PROMOTION_CAPTURE:
		fallthrough
	case KNIGHT_PROMOTION:
		if b.sideToMove == WHITE {
			moveDstSetSquare = WHITE_KNIGHT
		} else {
			moveDstSetSquare = BLACK_KNIGHT
		}
	case BISHOP_PROMOTION_CAPTURE:
		fallthrough
	case BISHOP_PROMOTION:
		if b.sideToMove == WHITE {
			moveDstSetSquare = WHITE_BISHOP
		} else {
			moveDstSetSquare = BLACK_BISHOP
		}
	case ROOK_PROMOTION_CAPTURE:
		fallthrough
	case ROOK_PROMOTION:
		if b.sideToMove == WHITE {
			moveDstSetSquare = WHITE_ROOK
		} else {
			moveDstSetSquare = BLACK_ROOK
		}
	case QUEEN_PROMOTION_CAPTURE:
		fallthrough
	case QUEEN_PROMOTION:
		if b.sideToMove == WHITE {
			moveDstSetSquare = WHITE_QUEEN
		} else {
			moveDstSetSquare = BLACK_QUEEN
		}
	}

	// Update Undo (should happen before any state is modified)
	u := NewUndo(m, moveDstSetSquare, b.halfMove, b.castleRights, b.epCoord)
	b.undos[b.undoIndex] = u
	b.undoIndex++

	// Update board
	b.ClearSquare(mu.originCoord, mu.originSquare)
	if mu.dstSquare != EMPTY {
		b.ClearSquare(mu.dstCoord, mu.dstSquare)
	}
	b.SetSquare(mu.dstCoord, moveDstSetSquare)
	if castleRookSquare != EMPTY {
		b.ClearSquare(castleRookClearCoord, castleRookSquare)
		b.SetSquare(castleRookSetCoord, castleRookSquare)
	}
	if epCaptureSquare != EMPTY {
		b.ClearSquare(epCaptureCoord, epCaptureSquare)
	}
	b.UpdateUnionBitboards()

	// Update CastleRights
	if b.castleRights != 0 {
		b.UpdateCastleRights(mu.originCoord)
		b.UpdateCastleRights(mu.dstCoord)
	}

	// Update En Passant
	if b.epCoord != A1 {
		b.HashEnPassant()
		b.epCoord = A1
	}

	if epCoord != A1 {
		b.epCoord = epCoord
		b.HashEnPassant()
	}

	// Update Move Clocks
	if b.sideToMove == BLACK {
		b.fullMove++
	}
	b.halfMove++
	if mu.originSquare == WHITE_PAWN ||
		mu.originSquare == BLACK_PAWN ||
		mu.dstSquare != EMPTY {
		b.halfMove = 0
	}

	// Move Legality Check
	kingAttacked := b.CoordAttacked(b.kingCoords[b.sideToMove], b.sideToMove)
	if kingAttacked {
		return fmt.Errorf("king is attacked")
	}

	// Update Side
	b.HashSide()
	b.sideToMove ^= 1

	return nil
}

func (b *Board) UndoMove() error {
	if b.undoIndex <= 0 {
		return fmt.Errorf("invalid undo index: %v", b.undoIndex)
	}

	epCaptureSquare := EMPTY
	epCaptureCoord := A1

	castleRookSquare := EMPTY
	castleRookClearCoord := A1
	castleRookSetCoord := A1

	// Update Undo
	b.undoIndex--
	u := b.undos[b.undoIndex]
	uu := u.Unpack()

	// Update Side
	b.sideToMove ^= 1
	b.HashSide()

	switch uu.moveKind {
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch uu.dstCoord {
		case G1:
			castleRookSquare = WHITE_ROOK
			castleRookClearCoord = F1
			castleRookSetCoord = H1
		case C1:
			castleRookSquare = WHITE_ROOK
			castleRookClearCoord = D1
			castleRookSetCoord = A1
		case G8:
			castleRookSquare = BLACK_ROOK
			castleRookClearCoord = F8
			castleRookSetCoord = H8
		case C8:
			castleRookSquare = BLACK_ROOK
			castleRookClearCoord = D8
			castleRookSetCoord = A8
		}
	case EP_CAPTURE:
		if b.sideToMove == WHITE {
			epCaptureCoord = uu.dstCoord + Coord(SHIFT_VERTICAL)
			epCaptureSquare = BLACK_PAWN
		} else if b.sideToMove == BLACK {
			epCaptureCoord = uu.dstCoord - Coord(SHIFT_VERTICAL)
			epCaptureSquare = WHITE_PAWN

		}
	}

	// Update board
	b.ClearSquare(uu.dstCoord, uu.clearSquare)
	if uu.dstSquare != EMPTY {
		b.SetSquare(uu.dstCoord, uu.dstSquare)
	}
	b.SetSquare(uu.originCoord, uu.originSquare)
	if castleRookSquare != EMPTY {
		b.ClearSquare(castleRookClearCoord, castleRookSquare)
		b.SetSquare(castleRookSetCoord, castleRookSquare)
	}
	if epCaptureSquare != EMPTY {
		b.SetSquare(epCaptureCoord, epCaptureSquare)
	}
	b.UpdateUnionBitboards()

	// Update CastleRights
	b.HashCastling()
	b.castleRights = uu.castleRights
	b.HashCastling()

	// Update En Passant
	if b.epCoord != A1 {
		b.HashEnPassant()
		b.epCoord = A1
	}

	if uu.epCoord != A1 {
		b.epCoord = uu.epCoord
		b.HashEnPassant()
	}

	// Update Move Clocks
	if b.sideToMove == BLACK {
		b.fullMove--
	}
	b.halfMove = uu.halfMove

	return nil
}
