package main

var COORD_MASK_BITBOARDS [BOARD_SQUARES]Bitboard
var COORD_CLEAR_BITBOARDS [BOARD_SQUARES]Bitboard

var RANK_MASK_BITBOARDS [RANKS]Bitboard
var RANK_CLEAR_BITBOARDS [RANKS]Bitboard

var FILE_MASK_BITBOARDS [FILES]Bitboard
var FILE_CLEAR_BITBOARDS [FILES]Bitboard

var PAWN_ATTACKS [PLAYERS][BOARD_SQUARES]Bitboard

func InitBitboards() {
	for i := 0; i < BOARD_SQUARES; i++ {
		currentBitboard := Bitboard(1 << i)

		COORD_MASK_BITBOARDS[i] = currentBitboard
		COORD_CLEAR_BITBOARDS[i] = ^currentBitboard

		rank := i / RANKS
		file := i % FILES

		RANK_MASK_BITBOARDS[rank] |= currentBitboard
		RANK_CLEAR_BITBOARDS[rank] = ^RANK_MASK_BITBOARDS[rank]

		FILE_MASK_BITBOARDS[file] |= currentBitboard
		FILE_CLEAR_BITBOARDS[file] = ^FILE_MASK_BITBOARDS[file]

		InitPawnAttacksBitboard(i, &currentBitboard)
	}
}

func InitPawnAttacksBitboard(i int, bb *Bitboard) {
	PAWN_ATTACKS[WHITE][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) << SHIFT_NEG_DIAG
	PAWN_ATTACKS[WHITE][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) << SHIFT_POS_DIAG

	PAWN_ATTACKS[BLACK][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) >> SHIFT_POS_DIAG
	PAWN_ATTACKS[BLACK][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) >> SHIFT_NEG_DIAG
}
