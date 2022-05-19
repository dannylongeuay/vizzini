package main

var COORD_MASK_BITBOARDS [BOARD_SQUARES]Bitboard
var COORD_CLEAR_BITBOARDS [BOARD_SQUARES]Bitboard

var RANK_MASK_BITBOARDS [RANKS]Bitboard
var RANK_CLEAR_BITBOARDS [RANKS]Bitboard

var FILE_MASK_BITBOARDS [FILES]Bitboard
var FILE_CLEAR_BITBOARDS [FILES]Bitboard

var BISHOP_MASKS [BOARD_SQUARES]Bitboard
var ROOK_MASKS [BOARD_SQUARES]Bitboard

var PAWN_ATTACKS [PLAYERS][BOARD_SQUARES]Bitboard
var KNIGHT_ATTACKS [BOARD_SQUARES]Bitboard
var KING_ATTACKS [BOARD_SQUARES]Bitboard

var ROOK_SHIFTS = [BOARD_SQUARES]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var BISHOP_SHIFTS = [BOARD_SQUARES]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

func InitBitboards() {
	for i := 0; i < BOARD_SQUARES; i++ {
		bb := Bitboard(1 << i)

		COORD_MASK_BITBOARDS[i] = bb
		COORD_CLEAR_BITBOARDS[i] = ^bb

		rank := i / RANKS
		file := i % FILES

		RANK_MASK_BITBOARDS[rank] |= bb
		RANK_CLEAR_BITBOARDS[rank] = ^RANK_MASK_BITBOARDS[rank]

		FILE_MASK_BITBOARDS[file] |= bb
		FILE_CLEAR_BITBOARDS[file] = ^FILE_MASK_BITBOARDS[file]

		InitPawnAttacksBitboard(i, &bb)
		InitKnightAttacksBitboard(i, &bb)
		InitKingAttacksBitboard(i, &bb)
		InitBishopMasksBitboard(i)
		InitRookMasksBitboard(i)
	}
}

func InitPawnAttacksBitboard(i int, bb *Bitboard) {
	PAWN_ATTACKS[WHITE][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) << SHIFT_NEG_DIAG
	PAWN_ATTACKS[WHITE][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) << SHIFT_POS_DIAG

	PAWN_ATTACKS[BLACK][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) >> SHIFT_POS_DIAG
	PAWN_ATTACKS[BLACK][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) >> SHIFT_NEG_DIAG
}

func InitKnightAttacksBitboard(i int, bb *Bitboard) {
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A] &
		FILE_CLEAR_BITBOARDS[FILE_A]) << 6
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A] &
		FILE_CLEAR_BITBOARDS[FILE_A]) >> 10

	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) << 15
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) >> 17

	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) << 17
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) >> 15

	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_G] &
		FILE_CLEAR_BITBOARDS[FILE_H]) << 10
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_G] &
		FILE_CLEAR_BITBOARDS[FILE_H]) >> 6
}

func InitKingAttacksBitboard(i int, bb *Bitboard) {
	KING_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) << SHIFT_NEG_DIAG
	KING_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) >> SHIFT_HORIZONTAL
	KING_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) >> SHIFT_POS_DIAG

	KING_ATTACKS[i] |= *bb << SHIFT_VERTICAL
	KING_ATTACKS[i] |= *bb >> SHIFT_VERTICAL

	KING_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) << SHIFT_POS_DIAG
	KING_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) << SHIFT_HORIZONTAL
	KING_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) >> SHIFT_NEG_DIAG

}

func InitBishopMasksBitboard(i int) {
	r := i / RANKS
	f := i % FILES

	// POS DIAG FORWARD
	for rr, ff := r+1, f+1; rr <= 6 && ff <= 6; rr, ff = rr+1, ff+1 {
		c := Coord(rr*RANKS + ff)
		BISHOP_MASKS[i] |= Bitboard(1 << c)
	}

	// POS DIAG BACKWARD
	for rr, ff := r-1, f-1; rr >= 1 && ff >= 1; rr, ff = rr-1, ff-1 {
		c := Coord(rr*RANKS + ff)
		BISHOP_MASKS[i] |= Bitboard(1 << c)
	}

	// NEG DIAG FORWARD
	for rr, ff := r-1, f+1; rr >= 1 && ff <= 6; rr, ff = rr-1, ff+1 {
		c := Coord(rr*RANKS + ff)
		BISHOP_MASKS[i] |= Bitboard(1 << c)
	}

	// NEG DIAG BACKWARD
	for rr, ff := r+1, f-1; rr <= 6 && ff >= 1; rr, ff = rr+1, ff-1 {
		c := Coord(rr*RANKS + ff)
		BISHOP_MASKS[i] |= Bitboard(1 << c)
	}
}

func GenerateBishopAttacksBitboard(i int, blockers Bitboard) Bitboard {
	var bb Bitboard
	r := i / RANKS
	f := i % FILES

	// POS DIAG FORWARD
	for rr, ff := r+1, f+1; rr <= 7 && ff <= 7; rr, ff = rr+1, ff+1 {
		c := Coord(rr*RANKS + ff)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}

	// POS DIAG BACKWARD
	for rr, ff := r-1, f-1; rr >= 0 && ff >= 0; rr, ff = rr-1, ff-1 {
		c := Coord(rr*RANKS + ff)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}

	// NEG DIAG FORWARD
	for rr, ff := r-1, f+1; rr >= 0 && ff <= 7; rr, ff = rr-1, ff+1 {
		c := Coord(rr*RANKS + ff)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}

	// NEG DIAG BACKWARD
	for rr, ff := r+1, f-1; rr <= 7 && ff >= 0; rr, ff = rr+1, ff-1 {
		c := Coord(rr*RANKS + ff)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}
	return bb
}

func InitRookMasksBitboard(i int) {
	r := i / RANKS
	f := i % FILES

	// POS VERTICAL
	for rr := r + 1; rr <= 6; rr++ {
		c := Coord(rr*RANKS + f)
		ROOK_MASKS[i] |= Bitboard(1 << c)
	}

	// POS HORIZONTAL
	for ff := f + 1; ff <= 6; ff++ {
		c := Coord(r*RANKS + ff)
		ROOK_MASKS[i] |= Bitboard(1 << c)
	}

	// NEG VERTICAL
	for rr := r - 1; rr >= 1; rr-- {
		c := Coord(rr*RANKS + f)
		ROOK_MASKS[i] |= Bitboard(1 << c)
	}

	// NEG HORIZONTAL
	for ff := f - 1; ff >= 1; ff-- {
		c := Coord(r*RANKS + ff)
		ROOK_MASKS[i] |= Bitboard(1 << c)
	}
}

func GenerateRookAttacksBitboard(i int, blockers Bitboard) Bitboard {
	var bb Bitboard
	r := i / RANKS
	f := i % FILES

	// POS VERTICAL
	for rr := r + 1; rr <= 7; rr++ {
		c := Coord(rr*RANKS + f)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}

	// POS HORIZONTAL
	for ff := f + 1; ff <= 7; ff++ {
		c := Coord(r*RANKS + ff)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}

	// NEG VERTICAL
	for rr := r - 1; rr >= 0; rr-- {
		c := Coord(rr*RANKS + f)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}

	// NEG HORIZONTAL
	for ff := f - 1; ff >= 0; ff-- {
		c := Coord(r*RANKS + ff)
		bb |= Bitboard(1 << c)
		if blockers&Bitboard(1<<c) > 0 {
			break
		}
	}
	return bb
}
