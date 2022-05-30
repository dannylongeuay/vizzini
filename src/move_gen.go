package main

func AppendMove(moves *[]Move, originCoord Coord, dstCoord Coord, originSquare Square, dstSquare Square, moveKind MoveKind, moveOrder MoveOrder) {
	move := NewMove(originCoord, dstCoord, originSquare, dstSquare, moveKind, moveOrder)
	*moves = append(*moves, move)
}

func AppendQuietMoves(moves *[]Move, quiet Bitboard, originCoord Coord, originSquare Square) {
	for quiet > 0 {
		dstCoord := quiet.PopLSB()
		AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, QUIET, 0)
	}
}

func (b *Board) AppendCaptures(moves *[]Move, captures Bitboard, originCoord Coord, originSquare Square) {
	for captures > 0 {
		dstCoord := captures.PopLSB()
		dstSquare := b.squares[dstCoord]
		moveOrder := MoveOrder(MVV_LVA_SCORES[dstSquare][originSquare])
		AppendMove(moves, originCoord, dstCoord, originSquare, dstSquare, CAPTURE, moveOrder)
	}
}

func (b *Board) GenerateMoves(moves *[]Move, side Color) {
	b.GeneratePawnMoves(moves, side)
	b.GenerateKnightMoves(moves, side)
	b.GenerateBishopMoves(moves, side)
	b.GenerateRookMoves(moves, side)
	b.GenerateQueenMoves(moves, side)
	b.GenerateKingMoves(moves, side)
}

func (b *Board) GeneratePawnMoves(moves *[]Move, side Color) {
	bbP := b.bbWP
	originSquare := WHITE_PAWN
	bbOpponentPieces := b.bbBlackPieces
	doublePushRankMask := RANK_MASK_BITBOARDS[RANK_THREE]
	promotionRank := RANK_EIGHT
	if side == BLACK {
		bbP = b.bbBP
		originSquare = BLACK_PAWN
		doublePushRankMask = RANK_MASK_BITBOARDS[RANK_SIX]
		bbOpponentPieces = b.bbWhitePieces
		promotionRank = RANK_ONE
	}
	for bbP > 0 {
		originCoord := bbP.PopLSB()
		captures := PAWN_ATTACKS[side][originCoord] & bbOpponentPieces & RANK_CLEAR_BITBOARDS[promotionRank]
		promotionCaptures := PAWN_ATTACKS[side][originCoord] & bbOpponentPieces & RANK_MASK_BITBOARDS[promotionRank]
		var quiet Bitboard
		var promotions Bitboard
		var doublePush Bitboard
		if side == WHITE {
			quiet = (Bitboard(1<<originCoord) << SHIFT_VERTICAL) & ^b.bbAllPieces & RANK_CLEAR_BITBOARDS[promotionRank]
			promotions = (Bitboard(1<<originCoord) << SHIFT_VERTICAL) & ^b.bbAllPieces & RANK_MASK_BITBOARDS[promotionRank]
			doublePush = ((quiet & doublePushRankMask) << SHIFT_VERTICAL) & ^b.bbAllPieces
		} else if side == BLACK {
			quiet = (Bitboard(1<<originCoord) >> SHIFT_VERTICAL) & ^b.bbAllPieces & RANK_CLEAR_BITBOARDS[promotionRank]
			promotions = (Bitboard(1<<originCoord) >> SHIFT_VERTICAL) & ^b.bbAllPieces & RANK_MASK_BITBOARDS[promotionRank]
			doublePush = ((quiet & doublePushRankMask) >> SHIFT_VERTICAL) & ^b.bbAllPieces
		}

		// Normal
		AppendQuietMoves(moves, quiet, originCoord, originSquare)
		b.AppendCaptures(moves, captures, originCoord, originSquare)

		// Special
		for doublePush > 0 {
			dstCoord := doublePush.PopLSB()
			AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, DOUBLE_PAWN_PUSH, 0)
		}
		if b.epCoord != A1 {
			epAttack := PAWN_ATTACKS[side][originCoord] & Bitboard(1<<b.epCoord)
			for epAttack > 0 {
				dstCoord := epAttack.PopLSB()
				AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, EP_CAPTURE, MVV_LVA_EN_PASSANT)

			}
		}
		for promotions > 0 {
			dstCoord := promotions.PopLSB()
			AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, KNIGHT_PROMOTION, MVV_LVA_KNIGHT_PROMOTION)
			AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, BISHOP_PROMOTION, MVV_LVA_BISHOP_PROMOTION)
			AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, ROOK_PROMOTION, MVV_LVA_ROOK_PROMOTION)
			AppendMove(moves, originCoord, dstCoord, originSquare, EMPTY, QUEEN_PROMOTION, MVV_LVA_QUEEN_PROMOTION)
		}
		for promotionCaptures > 0 {
			dstCoord := promotionCaptures.PopLSB()
			dstSquare := b.squares[dstCoord]
			AppendMove(moves, originCoord, dstCoord, originSquare, dstSquare, KNIGHT_PROMOTION_CAPTURE, MVV_LVA_KNIGHT_PROMOTION_CAPTURE)
			AppendMove(moves, originCoord, dstCoord, originSquare, dstSquare, BISHOP_PROMOTION_CAPTURE, MVV_LVA_BISHOP_PROMOTION_CAPTURE)
			AppendMove(moves, originCoord, dstCoord, originSquare, dstSquare, ROOK_PROMOTION_CAPTURE, MVV_LVA_ROOK_PROMOTION_CAPTURE)
			AppendMove(moves, originCoord, dstCoord, originSquare, dstSquare, QUEEN_PROMOTION_CAPTURE, MVV_LVA_QUEEN_PROMOTION_CAPTURE)
		}
	}
}

func (b *Board) GenerateKnightMoves(moves *[]Move, side Color) {
	bbN := b.bbWN
	originSquare := WHITE_KNIGHT
	bbOpponentPieces := b.bbBlackPieces
	if side == BLACK {
		bbN = b.bbBN
		originSquare = BLACK_KNIGHT
		bbOpponentPieces = b.bbWhitePieces
	}
	for bbN > 0 {
		originCoord := bbN.PopLSB()
		quiet := KNIGHT_ATTACKS[originCoord] & ^b.bbAllPieces
		captures := KNIGHT_ATTACKS[originCoord] & bbOpponentPieces

		AppendQuietMoves(moves, quiet, originCoord, originSquare)
		b.AppendCaptures(moves, captures, originCoord, originSquare)
	}
}

func (b *Board) GenerateBishopMoves(moves *[]Move, side Color) {
	bbB := b.bbWB
	originSquare := WHITE_BISHOP
	bbOpponentPieces := b.bbBlackPieces
	if side == BLACK {
		bbB = b.bbBB
		originSquare = BLACK_BISHOP
		bbOpponentPieces = b.bbWhitePieces
	}
	for bbB > 0 {
		originCoord := bbB.PopLSB()
		bishopAttacks := BishopAttacks(originCoord, b.bbAllPieces)
		quiet := bishopAttacks & ^b.bbAllPieces
		captures := bishopAttacks & bbOpponentPieces

		AppendQuietMoves(moves, quiet, originCoord, originSquare)
		b.AppendCaptures(moves, captures, originCoord, originSquare)
	}
}

func (b *Board) GenerateRookMoves(moves *[]Move, side Color) {
	bbR := b.bbWR
	originSquare := WHITE_ROOK
	bbOpponentPieces := b.bbBlackPieces
	if side == BLACK {
		bbR = b.bbBR
		originSquare = BLACK_ROOK
		bbOpponentPieces = b.bbWhitePieces
	}
	for bbR > 0 {
		originCoord := bbR.PopLSB()
		rookAttacks := RookAttacks(originCoord, b.bbAllPieces)
		quiet := rookAttacks & ^b.bbAllPieces
		captures := rookAttacks & bbOpponentPieces

		AppendQuietMoves(moves, quiet, originCoord, originSquare)
		b.AppendCaptures(moves, captures, originCoord, originSquare)
	}
}

func (b *Board) GenerateQueenMoves(moves *[]Move, side Color) {
	bbQ := b.bbWQ
	originSquare := WHITE_QUEEN
	bbOpponentPieces := b.bbBlackPieces
	if side == BLACK {
		bbQ = b.bbBQ
		originSquare = BLACK_QUEEN
		bbOpponentPieces = b.bbWhitePieces
	}
	for bbQ > 0 {
		originCoord := bbQ.PopLSB()

		bishopAttacks := BishopAttacks(originCoord, b.bbAllPieces)
		bishopQuiet := bishopAttacks & ^b.bbAllPieces
		bishopCaptures := bishopAttacks & bbOpponentPieces

		rookAttacks := RookAttacks(originCoord, b.bbAllPieces)
		rookQuiet := rookAttacks & ^b.bbAllPieces
		rookCaptures := rookAttacks & bbOpponentPieces

		quiet := bishopQuiet | rookQuiet
		captures := bishopCaptures | rookCaptures

		AppendQuietMoves(moves, quiet, originCoord, originSquare)
		b.AppendCaptures(moves, captures, originCoord, originSquare)
	}
}

func (b *Board) GenerateKingMoves(moves *[]Move, side Color) {
	bbK := b.bbWK
	originSquare := WHITE_KING
	bbOpponentPieces := b.bbBlackPieces
	castleRightsKingMask := CASTLING_RIGHTS_WHITE_KING_MASK
	castleRightsQueenMask := CASTLING_RIGHTS_WHITE_QUEEN_MASK
	castleKingDstCoord := G1
	castleQueenDstCoord := C1
	castleKingTravel := Bitboard(96)
	castleQueenTravel := Bitboard(14)
	castleKingAttackChecks := []Coord{E1, F1, G1}
	castleQueenAttackChecks := []Coord{E1, D1, C1}
	if side == BLACK {
		bbK = b.bbBK
		originSquare = BLACK_KING
		bbOpponentPieces = b.bbWhitePieces
		castleRightsKingMask = CASTLING_RIGHTS_BLACK_KING_MASK
		castleRightsQueenMask = CASTLING_RIGHTS_BLACK_QUEEN_MASK
		castleKingDstCoord = G8
		castleQueenDstCoord = C8
		castleKingTravel = Bitboard(6917529027641081856)
		castleQueenTravel = Bitboard(1008806316530991104)
		castleKingAttackChecks = []Coord{E8, F8, G8}
		castleQueenAttackChecks = []Coord{E8, D8, C8}
	}
	for bbK > 0 {
		originCoord := bbK.PopLSB()
		quiet := KING_ATTACKS[originCoord] & ^b.bbAllPieces
		captures := KING_ATTACKS[originCoord] & bbOpponentPieces

		// Normal
		AppendQuietMoves(moves, quiet, originCoord, originSquare)
		b.AppendCaptures(moves, captures, originCoord, originSquare)

		// Special
		if (b.castleRights&castleRightsKingMask) > 0 &&
			castleKingTravel&b.bbAllPieces == 0 &&
			!b.CoordsAttacked(castleKingAttackChecks, side) {
			AppendMove(moves, originCoord, castleKingDstCoord, originSquare, 0, KING_CASTLE, 0)
		}

		if (b.castleRights&castleRightsQueenMask) > 0 &&
			castleQueenTravel&b.bbAllPieces == 0 &&
			!b.CoordsAttacked(castleQueenAttackChecks, side) {
			AppendMove(moves, originCoord, castleQueenDstCoord, originSquare, 0, QUEEN_CASTLE, 0)
		}
	}
}

func (b *Board) CoordAttacked(c Coord, side Color) bool {
	bbOpponentPawns := b.bbBP
	bbOpponentKnights := b.bbBN
	bbOpponentBishops := b.bbBB
	bbOpponentRooks := b.bbBR
	bbOpponentQueens := b.bbBQ
	bbOpponentKing := b.bbBK
	if side == BLACK {
		bbOpponentPawns = b.bbWP
		bbOpponentKnights = b.bbWN
		bbOpponentBishops = b.bbWB
		bbOpponentRooks = b.bbWR
		bbOpponentQueens = b.bbWQ
		bbOpponentKing = b.bbWK
	}

	if PAWN_ATTACKS[side][c]&bbOpponentPawns > 0 {
		return true
	}

	if KNIGHT_ATTACKS[c]&bbOpponentKnights > 0 {
		return true
	}

	bishopAttacks := BishopAttacks(c, b.bbAllPieces)
	if bishopAttacks&(bbOpponentBishops|bbOpponentQueens) > 0 {
		return true
	}

	rookAttacks := RookAttacks(c, b.bbAllPieces)
	if rookAttacks&(bbOpponentRooks|bbOpponentQueens) > 0 {
		return true
	}

	if KING_ATTACKS[c]&bbOpponentKing > 0 {
		return true
	}

	return false
}

func (b *Board) CoordsAttacked(coords []Coord, side Color) bool {
	for _, c := range coords {
		if b.CoordAttacked(c, side) {
			return true
		}
	}
	return false
}
