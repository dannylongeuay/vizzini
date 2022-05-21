package main

func (b *Board) GenerateMoves(side Color) []Move {
	moves := make([]Move, 0, MAX_GENERATED_MOVES)
	return moves
}

func (b *Board) GeneratePawnMoves(moves *[]Move, side Color) {
	bbP := b.bbWP
	originSquare := WHITE_PAWN
	bbOpponentPieces := b.bbBlackPieces
	doublePushRankMask := RANK_MASK_BITBOARDS[RANK_THREE]
	epSquare := BLACK_PAWN
	if side == BLACK {
		bbP = b.bbBP
		originSquare = BLACK_PAWN
		doublePushRankMask = RANK_MASK_BITBOARDS[RANK_SIX]
		bbOpponentPieces = b.bbWhitePieces
		epSquare = WHITE_PAWN
	}
	for bbP > 0 {
		originCoord := bbP.PopLSB()
		// TODO: Handle enPassant
		pawnAttacks := PAWN_ATTACKS[side][originCoord] & bbOpponentPieces
		var quiet Bitboard
		var doublePush Bitboard
		if side == WHITE {
			quiet = (Bitboard(1<<originCoord) << SHIFT_VERTICAL) & ^b.bbAllPieces
			doublePush = ((quiet & doublePushRankMask) << SHIFT_VERTICAL) & ^b.bbAllPieces
		} else if side == BLACK {
			quiet = (Bitboard(1<<originCoord) >> SHIFT_VERTICAL) & ^b.bbAllPieces
			doublePush = ((quiet & doublePushRankMask) >> SHIFT_VERTICAL) & ^b.bbAllPieces
		}
		for pawnAttacks > 0 {
			dstCoord := pawnAttacks.PopLSB()
			dstSquare := b.squares[dstCoord]
			AppendMove(moves, originCoord, dstCoord, originSquare, dstSquare, CAPTURE)
		}
		for quiet > 0 {
			dstCoord := quiet.PopLSB()
			AppendMove(moves, originCoord, dstCoord, originSquare, 0, QUIET)
		}
		for doublePush > 0 {
			dstCoord := doublePush.PopLSB()
			AppendMove(moves, originCoord, dstCoord, originSquare, 0, DOUBLE_PAWN_PUSH)
		}
		if b.epCoord != A1 {
			epAttack := PAWN_ATTACKS[side][originCoord] & Bitboard(1<<b.epCoord)
			for epAttack > 0 {
				dstCoord := epAttack.PopLSB()
				AppendMove(moves, originCoord, dstCoord, originSquare, epSquare, EP_CAPTURE)

			}
		}
	}
}

func (b *Board) GenerateKnightMoves(moves *[]Move, side Color) {
}

func (b *Board) GenerateBishopMoves(moves *[]Move, side Color) {
}

func (b *Board) GenerateRookMoves(moves *[]Move, side Color) {
}

func (b *Board) GenerateQueenMoves(moves *[]Move, side Color) {
}

func (b *Board) GenerateKingMoves(moves *[]Move, side Color) {
}
