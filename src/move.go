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

func NewMove(
	originCoord Coord,
	dstCoord Coord,
	originSquare Square,
	dstSquare Square,
	moveKind MoveKind,
) Move {
	moveOriginCoord := Move(originCoord) << MOVE_ORIGIN_COORD_SHIFT
	moveDstCoord := Move(dstCoord) << MOVE_DST_COORD_SHIFT
	moveOriginSquare := Move(originSquare) << MOVE_ORIGIN_SQUARE_SHIFT
	moveDstSquare := Move(dstSquare) << MOVE_DST_SQUARE_SHIFT
	return moveOriginCoord | moveDstCoord | moveOriginSquare | moveDstSquare | Move(moveKind)
}

func NewMoveFromMoveUnpacked(mu MoveUnpacked) Move {
	return NewMove(mu.originCoord, mu.dstCoord, mu.originSquare, mu.dstSquare, mu.moveKind)
}

func NewUndo(
	move Move,
	clearSquare Square,
	halfMove HalfMove,
	castleRights CastleRights,
	epCoord Coord,
) Undo {
	undoClearSquare := Undo(clearSquare) << UNDO_CLEAR_SQUARE_SHIFT
	undoHalfMove := Undo(halfMove) << UNDO_HALF_MOVE_SHIFT
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
	clearSquare  Square
	halfMove     HalfMove
	castleRights CastleRights
	epCoord      Coord
}

func (m *Move) ToString() string {
	var mu MoveUnpacked
	m.Unpack(&mu)
	mk := MOVE_KIND_MAP[mu.moveKind]
	oc := COORD_STRINGS[mu.originCoord]
	os := SQUARES[mu.originSquare]
	dc := COORD_STRINGS[mu.dstCoord]
	ds := SQUARES[mu.dstSquare]
	s := fmt.Sprintf("Move{%v: %v(%v) -> %v(%v)}", mk, oc, os, dc, ds)
	return s
}

func (u *Undo) ToString() string {
	var mu MoveUnpacked
	var uu UndoUnpacked
	move := Move(*u & UNDO_MOVE_MASK)
	u.Unpack(&mu, &uu)
	cs := SQUARES[uu.clearSquare]
	hm := uu.halfMove
	cr := uu.castleRights
	ec := COORD_STRINGS[uu.epCoord]
	s := fmt.Sprintf("Undo{ClearSquare: %v | HalfMove: %v | CastleRights: %v | EP Coord: %v | %v}", cs, hm, cr, ec, move.ToString())
	return s
}

func (m *Move) Unpack(mu *MoveUnpacked) {
	mu.originCoord = Coord((*m & MOVE_ORIGIN_COORD_MASK) >> MOVE_ORIGIN_COORD_SHIFT)
	mu.dstCoord = Coord((*m & MOVE_DST_COORD_MASK) >> MOVE_DST_COORD_SHIFT)
	mu.originSquare = Square((*m & MOVE_ORIGIN_SQUARE_MASK) >> MOVE_ORIGIN_SQUARE_SHIFT)
	mu.dstSquare = Square((*m & MOVE_DST_SQUARE_MASK) >> MOVE_DST_SQUARE_SHIFT)
	mu.moveKind = MoveKind(*m & MOVE_KIND_MASK)
}

func (u *Undo) Unpack(mu *MoveUnpacked, uu *UndoUnpacked) {
	move := Move(*u & UNDO_MOVE_MASK)
	move.Unpack(mu)
	uu.clearSquare = Square((*u & UNDO_CLEAR_SQUARE_MASK) >> UNDO_CLEAR_SQUARE_SHIFT)
	uu.halfMove = HalfMove((*u & UNDO_HALF_MOVE_MASK) >> UNDO_HALF_MOVE_SHIFT)
	uu.castleRights = CastleRights((*u & UNDO_CASTLE_RIGHTS_MASK) >> UNDO_CASTLE_RIGHTS_SHIFT)
	uu.epCoord = Coord((*u & UNDO_EP_COORD_MASK) >> UNDO_EP_COORD_SHIFT)
}

func (b *Board) pieceBitboards(sq Square) (*Bitboard, *Bitboard) {
	switch sq {
	case WHITE_PAWN:
		return &b.bbWP, &b.bbWhitePieces
	case WHITE_KNIGHT:
		return &b.bbWN, &b.bbWhitePieces
	case WHITE_BISHOP:
		return &b.bbWB, &b.bbWhitePieces
	case WHITE_ROOK:
		return &b.bbWR, &b.bbWhitePieces
	case WHITE_QUEEN:
		return &b.bbWQ, &b.bbWhitePieces
	case WHITE_KING:
		return &b.bbWK, &b.bbWhitePieces
	case BLACK_PAWN:
		return &b.bbBP, &b.bbBlackPieces
	case BLACK_KNIGHT:
		return &b.bbBN, &b.bbBlackPieces
	case BLACK_BISHOP:
		return &b.bbBB, &b.bbBlackPieces
	case BLACK_ROOK:
		return &b.bbBR, &b.bbBlackPieces
	case BLACK_QUEEN:
		return &b.bbBQ, &b.bbBlackPieces
	case BLACK_KING:
		return &b.bbBK, &b.bbBlackPieces
	default:
		panic(fmt.Errorf("invalid square: %v", SQUARES[sq]))
	}
}

func (b *Board) SetSquare(c Coord, sq Square) {
	b.squares[c] = sq
	b.HashSquare(sq, c)
	b.bbAllPieces.SetBit(c)
	piece, side := b.pieceBitboards(sq)
	piece.SetBit(c)
	side.SetBit(c)
	if sq == WHITE_KING || sq == BLACK_KING {
		b.kingCoords[b.sideToMove] = c
	}
}

func (b *Board) ClearSquare(c Coord, sq Square) {
	if sq != b.squares[c] {
		panic(fmt.Errorf("clearing square mismatch %v != %v at coord %v\n\n%v", SQUARES[sq], SQUARES[b.squares[c]], COORD_STRINGS[c], b.ToString()))
	}
	b.squares[c] = EMPTY
	b.HashSquare(sq, c)
	b.bbAllPieces.ClearBit(c)
	piece, side := b.pieceBitboards(sq)
	piece.ClearBit(c)
	side.ClearBit(c)
}

func (b *Board) UpdateCastleRights(c Coord) {
	var castleBits CastleRights
	switch c {
	case A8:
		castleBits = 14
	case H8:
		castleBits = 13
	case E8:
		castleBits = 12
	case A1:
		castleBits = 11
	case H1:
		castleBits = 7
	case E1:
		castleBits = 3
	}
	if castleBits > 0 {
		b.HashCastling()
		b.castleRights &= castleBits
		b.HashCastling()
	}

}

func (b *Board) MakeMove(m Move) error {
	var mu MoveUnpacked
	m.Unpack(&mu)
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
	b.undos[b.ply] = u
	b.hashes[b.ply] = b.hash
	b.ply++

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

	// Update Side
	b.HashSide()
	b.sideToMove ^= 1

	if kingAttacked {
		return fmt.Errorf("king is attacked")
	}

	return nil
}

func (b *Board) UndoMove() error {
	if b.ply <= 0 {
		return fmt.Errorf("invalid undo index: %v", b.ply)
	}

	epCaptureSquare := EMPTY
	epCaptureCoord := A1

	castleRookSquare := EMPTY
	castleRookClearCoord := A1
	castleRookSetCoord := A1

	// Update Undo
	b.ply--
	u := b.undos[b.ply]
	var mu MoveUnpacked
	var uu UndoUnpacked
	u.Unpack(&mu, &uu)

	// Update Side
	b.sideToMove ^= 1
	b.HashSide()

	switch mu.moveKind {
	case KING_CASTLE:
		fallthrough
	case QUEEN_CASTLE:
		switch mu.dstCoord {
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
			epCaptureCoord = uu.epCoord - Coord(SHIFT_VERTICAL)
			epCaptureSquare = BLACK_PAWN
		} else if b.sideToMove == BLACK {
			epCaptureCoord = uu.epCoord + Coord(SHIFT_VERTICAL)
			epCaptureSquare = WHITE_PAWN

		}
	}

	// Update board
	b.ClearSquare(mu.dstCoord, uu.clearSquare)
	if mu.dstSquare != EMPTY {
		b.SetSquare(mu.dstCoord, mu.dstSquare)
	}
	b.SetSquare(mu.originCoord, mu.originSquare)
	if castleRookSquare != EMPTY {
		b.ClearSquare(castleRookClearCoord, castleRookSquare)
		b.SetSquare(castleRookSetCoord, castleRookSquare)
	}
	if epCaptureSquare != EMPTY {
		b.SetSquare(epCaptureCoord, epCaptureSquare)
	}

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

func (b *Board) MoveExists(move Move) bool {
	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	b.GenerateMoves(&moves, b.sideToMove, false)
	for _, m := range moves {
		if m == move {
			err := b.MakeMove(m)
			b.UndoMove()
			return err == nil
		}
	}
	return false
}
