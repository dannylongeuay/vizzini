package main

var COORD_MASK_BITBOARDS [BOARD_SQUARES]Bitboard
var COORD_CLEAR_BITBOARDS [BOARD_SQUARES]Bitboard

var RANK_MASK_BITBOARDS [RANKS]Bitboard
var RANK_CLEAR_BITBOARDS [RANKS]Bitboard

var FILE_MASK_BITBOARDS [FILES]Bitboard
var FILE_CLEAR_BITBOARDS [FILES]Bitboard

var PAWN_ATTACKS [PLAYERS][BOARD_SQUARES]Bitboard
var KNIGHT_ATTACKS [BOARD_SQUARES]Bitboard
var KING_ATTACKS [BOARD_SQUARES]Bitboard
var BISHOP_ATTACKS [BOARD_SQUARES]Bitboard
var ROOK_ATTACKS [BOARD_SQUARES]Bitboard

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
		InitBishopAttacksBitboard(i)
		InitRookAttacksBitboard(i)
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

func InitBishopAttacksBitboard(i int) {
	r := i / RANKS
	f := i % FILES

	// POS DIAG FORWARD
	rr := r + 1
	ff := f + 1
	for rr <= 6 && ff <= 6 {
		c := Coord(rr*RANKS + ff)
		BISHOP_ATTACKS[i] |= Bitboard(1 << c)
		rr++
		ff++
	}

	// POS DIAG BACKWARD
	rr = r - 1
	ff = f - 1
	for rr >= 1 && ff >= 1 {
		c := Coord(rr*RANKS + ff)
		BISHOP_ATTACKS[i] |= Bitboard(1 << c)
		rr--
		ff--
	}

	// NEG DIAG FORWARD
	rr = r - 1
	ff = f + 1
	for rr >= 1 && ff <= 6 {
		c := Coord(rr*RANKS + ff)
		BISHOP_ATTACKS[i] |= Bitboard(1 << c)
		rr--
		ff++
	}

	// NEG DIAG BACKWARD
	rr = r + 1
	ff = f - 1
	for rr <= 6 && ff >= 1 {
		c := Coord(rr*RANKS + ff)
		BISHOP_ATTACKS[i] |= Bitboard(1 << c)
		rr++
		ff--
	}
}

func InitRookAttacksBitboard(i int) {
	r := i / RANKS
	f := i % FILES

	// POS VERTICAL
	rr := r + 1
	for rr <= 6 {
		c := Coord(rr*RANKS + f)
		ROOK_ATTACKS[i] |= Bitboard(1 << c)
		rr++
	}

	// POS HORIZONTAL
	ff := f + 1
	for ff <= 6 {
		c := Coord(r*RANKS + ff)
		ROOK_ATTACKS[i] |= Bitboard(1 << c)
		ff++
	}

	// NEG VERTICAL
	rr = r - 1
	for rr >= 1 {
		c := Coord(rr*RANKS + f)
		ROOK_ATTACKS[i] |= Bitboard(1 << c)
		rr--
	}

	// NEG HORIZONTAL
	ff = f - 1
	for ff >= 1 {
		c := Coord(r*RANKS + ff)
		ROOK_ATTACKS[i] |= Bitboard(1 << c)
		ff--
	}
}
