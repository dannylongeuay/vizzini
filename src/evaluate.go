package main

// Phase weights per piece type (non-pawn, non-king).
const (
	PHASE_KNIGHT = 1
	PHASE_BISHOP = 1
	PHASE_ROOK   = 2
	PHASE_QUEEN  = 4

	// Total phase at game start: 4 knights + 4 bishops + 4 rooks + 2 queens = 24.
	PHASE_TOTAL = 24
)

// Midgame material values indexed by Square type.
var MATERIAL_MG = [SQUARE_TYPES]int{
	0, 100, 320, 330, 500, 900, 10000,
	-100, -320, -330, -500, -900, -10000,
}

// Endgame material values indexed by Square type.
var MATERIAL_EG = [SQUARE_TYPES]int{
	0, 120, 280, 300, 520, 950, 10000,
	-120, -280, -300, -520, -950, -10000,
}

// Phase weight indexed by Square type (for phase calculation).
var PHASE_WEIGHT = [SQUARE_TYPES]int{
	0, 0, PHASE_KNIGHT, PHASE_BISHOP, PHASE_ROOK, PHASE_QUEEN, 0,
	0, PHASE_KNIGHT, PHASE_BISHOP, PHASE_ROOK, PHASE_QUEEN, 0,
}

var VICTIM_SCORE = [SQUARE_TYPES]MoveOrder{
	0, 40, 50, 60, 70, 80, 90,
	40, 50, 60, 70, 80, 90,
}

var MVV_LVA_SCORES [SQUARE_TYPES][SQUARE_TYPES]MoveOrder

// Midgame piece-square tables.
var PAWN_MG = [BOARD_SQUARES]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 5, 5, -10, -10, 5, 5, 5,
	5, 0, 0, 5, 5, 0, 0, 5,
	5, 5, 5, 20, 20, 5, 5, 5,
	10, 10, 15, 25, 25, 15, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	30, 30, 30, 40, 40, 30, 30, 30,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KNIGHT_MG = [BOARD_SQUARES]int{
	-5, -10, -5, -5, -5, -5, -10, -5,
	-5, 0, 0, 5, 5, 0, 0, -5,
	-5, 5, 15, 15, 15, 15, 5, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 5, 15, 20, 20, 15, 5, -5,
	-5, 0, 0, 10, 10, 0, 0, -5,
	-5, -5, -5, -5, -5, -5, -5, -5,
}

var BISHOP_MG = [BOARD_SQUARES]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 0, 0, 10, 0, 0,
	0, 10, 0, 0, 0, 0, 10, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var ROOK_MG = [BOARD_SQUARES]int{
	0, 0, 5, 10, 10, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	20, 20, 20, 20, 20, 20, 20, 20,
	10, 10, 10, 10, 10, 10, 10, 10,
}

var QUEEN_MG = [BOARD_SQUARES]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-20, -10, -10, 0, 0, -10, -10, -20,
}

// Midgame king: prefer castled corners.
var KING_MG = [BOARD_SQUARES]int{
	20, 30, 10, 0, 0, 10, 30, 20,
	20, 20, 0, 0, 0, 0, 20, 20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
}

// Endgame piece-square tables.
var PAWN_EG = [BOARD_SQUARES]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 5, 5, 5, 5, 5, 5, 5,
	10, 10, 10, 10, 10, 10, 10, 10,
	15, 15, 15, 20, 20, 15, 15, 15,
	25, 25, 25, 30, 30, 25, 25, 25,
	40, 40, 40, 45, 45, 40, 40, 40,
	60, 60, 60, 65, 65, 60, 60, 60,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KNIGHT_EG = [BOARD_SQUARES]int{
	-20, -10, -5, -5, -5, -5, -10, -20,
	-10, -5, 0, 0, 0, 0, -5, -10,
	-5, 0, 10, 10, 10, 10, 0, -5,
	-5, 0, 10, 20, 20, 10, 0, -5,
	-5, 0, 10, 20, 20, 10, 0, -5,
	-5, 0, 10, 10, 10, 10, 0, -5,
	-10, -5, 0, 0, 0, 0, -5, -10,
	-20, -10, -5, -5, -5, -5, -10, -20,
}

var BISHOP_EG = [BOARD_SQUARES]int{
	-10, -5, -5, -5, -5, -5, -5, -10,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-5, 0, 5, 10, 10, 5, 0, -5,
	-5, 0, 5, 10, 10, 5, 0, -5,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-10, -5, -5, -5, -5, -5, -5, -10,
}

var ROOK_EG = [BOARD_SQUARES]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 5, 5, 5, 5, 5, 5, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var QUEEN_EG = [BOARD_SQUARES]int{
	-10, -5, -5, 0, 0, -5, -5, -10,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, 0,
	0, 0, 5, 5, 5, 5, 0, 0,
	-5, 0, 5, 5, 5, 5, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-10, -5, -5, 0, 0, -5, -5, -10,
}

// Endgame king: centralized and active.
var KING_EG = [BOARD_SQUARES]int{
	-50, -30, -20, -20, -20, -20, -30, -50,
	-30, -10, 0, 0, 0, 0, -10, -30,
	-20, 0, 10, 15, 15, 10, 0, -20,
	-20, 0, 15, 20, 20, 15, 0, -20,
	-20, 0, 15, 20, 20, 15, 0, -20,
	-20, 0, 10, 15, 15, 10, 0, -20,
	-30, -10, 0, 0, 0, 0, -10, -30,
	-50, -30, -20, -20, -20, -20, -30, -50,
}

// Mobility bonuses per safe square (excluding squares controlled by opponent pawns).
const (
	MOBILITY_KNIGHT_MG = 4
	MOBILITY_KNIGHT_EG = 3
	MOBILITY_BISHOP_MG = 5
	MOBILITY_BISHOP_EG = 4
	MOBILITY_ROOK_MG   = 2
	MOBILITY_ROOK_EG   = 3
	MOBILITY_QUEEN_MG  = 1
	MOBILITY_QUEEN_EG  = 2
)

// Bishop pair bonus.
const (
	BISHOP_PAIR_MG = 30
	BISHOP_PAIR_EG = 50
)

// King safety: attack weight indexed by piece kind (PAWN=1..KING=6).
// Knights=2, Bishops=2, Rooks=3, Queens=5. Others=0.
var KING_ATTACK_WEIGHT = [7]int{0, 0, 2, 2, 3, 5, 0}

// Non-linear danger table: maps cumulative attacker weight to centipawn penalty.
// Index = sum of KING_ATTACK_WEIGHT for all pieces attacking king zone. Capped at 15.
var KING_DANGER_TABLE = [16]int{
	0, 0, 3, 12, 28, 50, 78, 112,
	152, 198, 250, 308, 372, 442, 518, 600,
}

// King safety constants (MG only).
const (
	PAWN_SHIELD_BONUS_MG     = 10  // pawn on ideal shield rank
	PAWN_SHIELD_ADVANCE_MG   = 5   // pawn advanced one rank beyond ideal
	PAWN_SHIELD_MISSING_MG   = -15 // no pawn on this shield file at all
	KING_OPEN_FILE_MG        = -25 // no pawns of either color on file
	KING_SEMI_OPEN_FILE_MG   = -12 // no friendly pawns but enemy pawns present
	KING_SAFETY_MIN_MATERIAL = 4   // minimum opponent attacking material (phase weight)
)

// Pawn shield masks for king safety (initialized in InitEval).
var WHITE_PAWN_SHIELD [BOARD_SQUARES]Bitboard
var BLACK_PAWN_SHIELD [BOARD_SQUARES]Bitboard

// Advanced pawn shield masks (2 ranks ahead/behind king, initialized in InitEval).
var WHITE_PAWN_SHIELD_ADV [BOARD_SQUARES]Bitboard
var BLACK_PAWN_SHIELD_ADV [BOARD_SQUARES]Bitboard

// Pawn structure bonuses/penalties.
const PASSED_PAWN_MG = 10
const PASSED_PAWN_EG = 20
const ISOLATED_PAWN_MG = -10
const ISOLATED_PAWN_EG = -15
const DOUBLED_PAWN_MG = -10
const DOUBLED_PAWN_EG = -20

// Passed pawn rank bonus multiplier (indexed by rank for white perspective).
var PASSED_PAWN_RANK_BONUS = [RANKS]int{0, 0, 0, 1, 2, 4, 6, 0}

// Adjacent file masks for isolated pawn detection (initialized in InitEval).
var ADJACENT_FILE_MASKS [FILES]Bitboard

// Front span masks for passed pawn detection (initialized in InitEval).
var WHITE_FRONT_SPAN [BOARD_SQUARES]Bitboard
var BLACK_FRONT_SPAN [BOARD_SQUARES]Bitboard

func InitEval() {
	// Initialize adjacent file masks.
	for f := 0; f < FILES; f++ {
		ADJACENT_FILE_MASKS[f] = 0
		if f > 0 {
			ADJACENT_FILE_MASKS[f] |= FILE_MASK_BITBOARDS[f-1]
		}
		if f < FILES-1 {
			ADJACENT_FILE_MASKS[f] |= FILE_MASK_BITBOARDS[f+1]
		}
	}

	// Initialize front span masks.
	// White front span at coord: all squares on same file + adjacent files, on ranks above.
	// Black front span at coord: all squares on same file + adjacent files, on ranks below.
	for coord := 0; coord < BOARD_SQUARES; coord++ {
		file := coord % FILES
		rank := coord / FILES

		fileMask := FILE_MASK_BITBOARDS[file] | ADJACENT_FILE_MASKS[file]

		var whiteSpan Bitboard
		for r := rank + 1; r < RANKS; r++ {
			whiteSpan |= RANK_MASK_BITBOARDS[r]
		}
		WHITE_FRONT_SPAN[coord] = fileMask & whiteSpan

		var blackSpan Bitboard
		for r := 0; r < rank; r++ {
			blackSpan |= RANK_MASK_BITBOARDS[r]
		}
		BLACK_FRONT_SPAN[coord] = fileMask & blackSpan
	}

	// Initialize pawn shield masks (immediate and advanced).
	for coord := 0; coord < BOARD_SQUARES; coord++ {
		file := coord % FILES
		rank := coord / FILES

		// White pawn shield: squares one rank ahead on king file and adjacent files.
		if rank < RANKS-1 {
			shieldRank := rank + 1
			WHITE_PAWN_SHIELD[coord] |= COORD_MASK_BITBOARDS[shieldRank*FILES+file]
			if file > 0 {
				WHITE_PAWN_SHIELD[coord] |= COORD_MASK_BITBOARDS[shieldRank*FILES+file-1]
			}
			if file < FILES-1 {
				WHITE_PAWN_SHIELD[coord] |= COORD_MASK_BITBOARDS[shieldRank*FILES+file+1]
			}
		}

		// White advanced pawn shield: two ranks ahead.
		if rank+2 < RANKS {
			advRank := rank + 2
			WHITE_PAWN_SHIELD_ADV[coord] |= COORD_MASK_BITBOARDS[advRank*FILES+file]
			if file > 0 {
				WHITE_PAWN_SHIELD_ADV[coord] |= COORD_MASK_BITBOARDS[advRank*FILES+file-1]
			}
			if file < FILES-1 {
				WHITE_PAWN_SHIELD_ADV[coord] |= COORD_MASK_BITBOARDS[advRank*FILES+file+1]
			}
		}

		// Black pawn shield: squares one rank behind on king file and adjacent files.
		if rank > 0 {
			shieldRank := rank - 1
			BLACK_PAWN_SHIELD[coord] |= COORD_MASK_BITBOARDS[shieldRank*FILES+file]
			if file > 0 {
				BLACK_PAWN_SHIELD[coord] |= COORD_MASK_BITBOARDS[shieldRank*FILES+file-1]
			}
			if file < FILES-1 {
				BLACK_PAWN_SHIELD[coord] |= COORD_MASK_BITBOARDS[shieldRank*FILES+file+1]
			}
		}

		// Black advanced pawn shield: two ranks behind.
		if rank-2 >= 0 {
			advRank := rank - 2
			BLACK_PAWN_SHIELD_ADV[coord] |= COORD_MASK_BITBOARDS[advRank*FILES+file]
			if file > 0 {
				BLACK_PAWN_SHIELD_ADV[coord] |= COORD_MASK_BITBOARDS[advRank*FILES+file-1]
			}
			if file < FILES-1 {
				BLACK_PAWN_SHIELD_ADV[coord] |= COORD_MASK_BITBOARDS[advRank*FILES+file+1]
			}
		}
	}

	InitMvvLva()
}

func InitMvvLva() {
	for attacker := WHITE_PAWN; attacker <= BLACK_KING; attacker++ {
		for victim := WHITE_PAWN; victim <= BLACK_KING; victim++ {
			MVV_LVA_SCORES[victim][attacker] = VICTIM_SCORE[victim] + 10 - (VICTIM_SCORE[attacker] / 10)
		}
	}
}

// Phase computes the game phase from piece counts (non-destructive).
// Returns 0 (endgame) to PHASE_TOTAL (midgame).
func (b *Board) Phase() int {
	var phase int
	for _, bbSquare := range b.BitboardSquares() {
		phase += PHASE_WEIGHT[bbSquare.sq] * bbSquare.bb.Count()
	}
	if phase > PHASE_TOTAL {
		phase = PHASE_TOTAL
	}
	return phase
}

// TaperedMaterial returns the phase-blended material value for a piece type.
func TaperedMaterial(sq Square, phase int) int {
	return (MATERIAL_MG[sq]*phase + MATERIAL_EG[sq]*(PHASE_TOTAL-phase)) / PHASE_TOTAL
}

func (b *Board) Evaluate() int {
	var mgScore, egScore int
	phase := b.Phase()

	for _, bbSquare := range b.BitboardSquares() {
		for bbSquare.bb > 0 {
			coord := bbSquare.bb.PopLSB()
			mirror := MIRROR_COORDS[coord]

			switch bbSquare.sq {
			case WHITE_PAWN:
				mgScore += MATERIAL_MG[WHITE_PAWN] + PAWN_MG[coord]
				egScore += MATERIAL_EG[WHITE_PAWN] + PAWN_EG[coord]
			case WHITE_KNIGHT:
				mgScore += MATERIAL_MG[WHITE_KNIGHT] + KNIGHT_MG[coord]
				egScore += MATERIAL_EG[WHITE_KNIGHT] + KNIGHT_EG[coord]
			case WHITE_BISHOP:
				mgScore += MATERIAL_MG[WHITE_BISHOP] + BISHOP_MG[coord]
				egScore += MATERIAL_EG[WHITE_BISHOP] + BISHOP_EG[coord]
			case WHITE_ROOK:
				mgScore += MATERIAL_MG[WHITE_ROOK] + ROOK_MG[coord]
				egScore += MATERIAL_EG[WHITE_ROOK] + ROOK_EG[coord]
			case WHITE_QUEEN:
				mgScore += MATERIAL_MG[WHITE_QUEEN] + QUEEN_MG[coord]
				egScore += MATERIAL_EG[WHITE_QUEEN] + QUEEN_EG[coord]
			case WHITE_KING:
				mgScore += MATERIAL_MG[WHITE_KING] + KING_MG[coord]
				egScore += MATERIAL_EG[WHITE_KING] + KING_EG[coord]
			case BLACK_PAWN:
				mgScore -= MATERIAL_MG[WHITE_PAWN] + PAWN_MG[mirror]
				egScore -= MATERIAL_EG[WHITE_PAWN] + PAWN_EG[mirror]
			case BLACK_KNIGHT:
				mgScore -= MATERIAL_MG[WHITE_KNIGHT] + KNIGHT_MG[mirror]
				egScore -= MATERIAL_EG[WHITE_KNIGHT] + KNIGHT_EG[mirror]
			case BLACK_BISHOP:
				mgScore -= MATERIAL_MG[WHITE_BISHOP] + BISHOP_MG[mirror]
				egScore -= MATERIAL_EG[WHITE_BISHOP] + BISHOP_EG[mirror]
			case BLACK_ROOK:
				mgScore -= MATERIAL_MG[WHITE_ROOK] + ROOK_MG[mirror]
				egScore -= MATERIAL_EG[WHITE_ROOK] + ROOK_EG[mirror]
			case BLACK_QUEEN:
				mgScore -= MATERIAL_MG[WHITE_QUEEN] + QUEEN_MG[mirror]
				egScore -= MATERIAL_EG[WHITE_QUEEN] + QUEEN_EG[mirror]
			case BLACK_KING:
				mgScore -= MATERIAL_MG[WHITE_KING] + KING_MG[mirror]
				egScore -= MATERIAL_EG[WHITE_KING] + KING_EG[mirror]
			}
		}
	}

	// Pawn structure evaluation.
	mgPawn, egPawn := b.evaluatePawnStructure()
	mgScore += mgPawn
	egScore += egPawn

	// Mobility evaluation.
	mgMobility, egMobility := b.evaluateMobility()
	mgScore += mgMobility
	egScore += egMobility

	// King safety evaluation (MG only).
	mgScore += b.evaluateKingSafety()

	// Bishop pair bonus.
	if b.bbWB.Count() >= 2 {
		mgScore += BISHOP_PAIR_MG
		egScore += BISHOP_PAIR_EG
	}
	if b.bbBB.Count() >= 2 {
		mgScore -= BISHOP_PAIR_MG
		egScore -= BISHOP_PAIR_EG
	}

	// Taper: phase ranges from 0 (endgame) to PHASE_TOTAL (midgame).
	score := (mgScore*phase + egScore*(PHASE_TOTAL-phase)) / PHASE_TOTAL

	if b.sideToMove == WHITE {
		return score
	}
	return -score
}

func (b *Board) evaluatePawnStructure() (int, int) {
	var mgScore, egScore int

	// White pawns.
	wp := b.bbWP
	for wp > 0 {
		coord := wp.PopLSB()
		file := int(coord) % FILES
		rank := int(coord) / FILES

		// Passed pawn: no enemy pawns on same or adjacent files ahead.
		if WHITE_FRONT_SPAN[coord]&b.bbBP == 0 {
			bonus := PASSED_PAWN_RANK_BONUS[rank]
			mgScore += PASSED_PAWN_MG * bonus
			egScore += PASSED_PAWN_EG * bonus
		}

		// Isolated pawn: no friendly pawns on adjacent files.
		if ADJACENT_FILE_MASKS[file]&b.bbWP == 0 {
			mgScore += ISOLATED_PAWN_MG
			egScore += ISOLATED_PAWN_EG
		}
	}

	// White doubled pawns (per file).
	for f := 0; f < FILES; f++ {
		count := (FILE_MASK_BITBOARDS[f] & b.bbWP).Count()
		if count > 1 {
			mgScore += DOUBLED_PAWN_MG * (count - 1)
			egScore += DOUBLED_PAWN_EG * (count - 1)
		}
	}

	// Black pawns.
	bp := b.bbBP
	for bp > 0 {
		coord := bp.PopLSB()
		file := int(coord) % FILES
		rank := int(coord) / FILES

		// Passed pawn: no enemy pawns on same or adjacent files ahead (toward rank 0 for black).
		if BLACK_FRONT_SPAN[coord]&b.bbWP == 0 {
			bonus := PASSED_PAWN_RANK_BONUS[RANKS-1-rank]
			mgScore -= PASSED_PAWN_MG * bonus
			egScore -= PASSED_PAWN_EG * bonus
		}

		// Isolated pawn.
		if ADJACENT_FILE_MASKS[file]&b.bbBP == 0 {
			mgScore -= ISOLATED_PAWN_MG
			egScore -= ISOLATED_PAWN_EG
		}
	}

	// Black doubled pawns.
	for f := 0; f < FILES; f++ {
		count := (FILE_MASK_BITBOARDS[f] & b.bbBP).Count()
		if count > 1 {
			mgScore -= DOUBLED_PAWN_MG * (count - 1)
			egScore -= DOUBLED_PAWN_EG * (count - 1)
		}
	}

	return mgScore, egScore
}

func (b *Board) evaluateMobility() (int, int) {
	var mgScore, egScore int

	// Compute pawn attack masks for safe mobility (exclude squares attacked by enemy pawns).
	var whitePawnAttacks, blackPawnAttacks Bitboard
	wp := b.bbWP
	for wp > 0 {
		coord := wp.PopLSB()
		whitePawnAttacks |= PAWN_ATTACKS[WHITE][coord]
	}
	bp := b.bbBP
	for bp > 0 {
		coord := bp.PopLSB()
		blackPawnAttacks |= PAWN_ATTACKS[BLACK][coord]
	}

	// White knights.
	wn := b.bbWN
	for wn > 0 {
		coord := wn.PopLSB()
		mobility := (KNIGHT_ATTACKS[coord] & ^blackPawnAttacks & ^b.bbWhitePieces).Count()
		mgScore += MOBILITY_KNIGHT_MG * mobility
		egScore += MOBILITY_KNIGHT_EG * mobility
	}

	// White bishops.
	wb := b.bbWB
	for wb > 0 {
		coord := wb.PopLSB()
		mobility := (BishopAttacks(coord, b.bbAllPieces) & ^blackPawnAttacks & ^b.bbWhitePieces).Count()
		mgScore += MOBILITY_BISHOP_MG * mobility
		egScore += MOBILITY_BISHOP_EG * mobility
	}

	// White rooks.
	wr := b.bbWR
	for wr > 0 {
		coord := wr.PopLSB()
		mobility := (RookAttacks(coord, b.bbAllPieces) & ^blackPawnAttacks & ^b.bbWhitePieces).Count()
		mgScore += MOBILITY_ROOK_MG * mobility
		egScore += MOBILITY_ROOK_EG * mobility
	}

	// White queens.
	wq := b.bbWQ
	for wq > 0 {
		coord := wq.PopLSB()
		mobility := ((BishopAttacks(coord, b.bbAllPieces) | RookAttacks(coord, b.bbAllPieces)) & ^blackPawnAttacks & ^b.bbWhitePieces).Count()
		mgScore += MOBILITY_QUEEN_MG * mobility
		egScore += MOBILITY_QUEEN_EG * mobility
	}

	// Black knights.
	bn := b.bbBN
	for bn > 0 {
		coord := bn.PopLSB()
		mobility := (KNIGHT_ATTACKS[coord] & ^whitePawnAttacks & ^b.bbBlackPieces).Count()
		mgScore -= MOBILITY_KNIGHT_MG * mobility
		egScore -= MOBILITY_KNIGHT_EG * mobility
	}

	// Black bishops.
	bb := b.bbBB
	for bb > 0 {
		coord := bb.PopLSB()
		mobility := (BishopAttacks(coord, b.bbAllPieces) & ^whitePawnAttacks & ^b.bbBlackPieces).Count()
		mgScore -= MOBILITY_BISHOP_MG * mobility
		egScore -= MOBILITY_BISHOP_EG * mobility
	}

	// Black rooks.
	br := b.bbBR
	for br > 0 {
		coord := br.PopLSB()
		mobility := (RookAttacks(coord, b.bbAllPieces) & ^whitePawnAttacks & ^b.bbBlackPieces).Count()
		mgScore -= MOBILITY_ROOK_MG * mobility
		egScore -= MOBILITY_ROOK_EG * mobility
	}

	// Black queens.
	bq := b.bbBQ
	for bq > 0 {
		coord := bq.PopLSB()
		mobility := ((BishopAttacks(coord, b.bbAllPieces) | RookAttacks(coord, b.bbAllPieces)) & ^whitePawnAttacks & ^b.bbBlackPieces).Count()
		mgScore -= MOBILITY_QUEEN_MG * mobility
		egScore -= MOBILITY_QUEEN_EG * mobility
	}

	return mgScore, egScore
}

func (b *Board) evaluateKingSafety() int {
	var mgScore int

	// White king safety: penalty from black attackers.
	blackAttackMaterial := b.bbBN.Count()*PHASE_KNIGHT +
		b.bbBB.Count()*PHASE_BISHOP +
		b.bbBR.Count()*PHASE_ROOK +
		b.bbBQ.Count()*PHASE_QUEEN
	if blackAttackMaterial >= KING_SAFETY_MIN_MATERIAL {
		mgScore -= b.kingSafetyForSide(
			b.kingCoords[WHITE],
			b.bbBN, b.bbBB, b.bbBR, b.bbBQ,
			b.bbWP, b.bbBP,
			&WHITE_PAWN_SHIELD, &WHITE_PAWN_SHIELD_ADV,
		)
	}

	// Black king safety: penalty from white attackers.
	whiteAttackMaterial := b.bbWN.Count()*PHASE_KNIGHT +
		b.bbWB.Count()*PHASE_BISHOP +
		b.bbWR.Count()*PHASE_ROOK +
		b.bbWQ.Count()*PHASE_QUEEN
	if whiteAttackMaterial >= KING_SAFETY_MIN_MATERIAL {
		mgScore += b.kingSafetyForSide(
			b.kingCoords[BLACK],
			b.bbWN, b.bbWB, b.bbWR, b.bbWQ,
			b.bbBP, b.bbWP,
			&BLACK_PAWN_SHIELD, &BLACK_PAWN_SHIELD_ADV,
		)
	}

	return mgScore
}

// kingSafetyForSide computes the king safety penalty for the defending side.
// Returns a positive value representing how unsafe the king is.
func (b *Board) kingSafetyForSide(
	kingCoord Coord,
	attackerKnights, attackerBishops, attackerRooks, attackerQueens Bitboard,
	friendlyPawns, enemyPawns Bitboard,
	shieldTable, advShieldTable *[BOARD_SQUARES]Bitboard,
) int {
	var penalty int
	kingZone := KING_ATTACKS[kingCoord] | COORD_MASK_BITBOARDS[kingCoord]

	// Attacker weight accumulation.
	var attackWeight int
	knights := attackerKnights
	for knights > 0 {
		coord := knights.PopLSB()
		if KNIGHT_ATTACKS[coord]&kingZone != 0 {
			attackWeight += KING_ATTACK_WEIGHT[2] // KNIGHT
		}
	}
	bishops := attackerBishops
	for bishops > 0 {
		coord := bishops.PopLSB()
		if BishopAttacks(coord, b.bbAllPieces)&kingZone != 0 {
			attackWeight += KING_ATTACK_WEIGHT[3] // BISHOP
		}
	}
	rooks := attackerRooks
	for rooks > 0 {
		coord := rooks.PopLSB()
		if RookAttacks(coord, b.bbAllPieces)&kingZone != 0 {
			attackWeight += KING_ATTACK_WEIGHT[4] // ROOK
		}
	}
	queens := attackerQueens
	for queens > 0 {
		coord := queens.PopLSB()
		if (BishopAttacks(coord, b.bbAllPieces)|RookAttacks(coord, b.bbAllPieces))&kingZone != 0 {
			attackWeight += KING_ATTACK_WEIGHT[5] // QUEEN
		}
	}
	if attackWeight > 15 {
		attackWeight = 15
	}
	penalty += KING_DANGER_TABLE[attackWeight]

	// Pawn shield evaluation.
	kingFile := int(kingCoord) % FILES
	shieldMask := shieldTable[kingCoord]
	advancedShieldMask := advShieldTable[kingCoord]

	startFile := kingFile - 1
	if startFile < 0 {
		startFile = 0
	}
	endFile := kingFile + 1
	if endFile > FILES-1 {
		endFile = FILES - 1
	}
	for f := startFile; f <= endFile; f++ {
		filePawns := FILE_MASK_BITBOARDS[f] & friendlyPawns
		if filePawns&shieldMask != 0 {
			penalty -= PAWN_SHIELD_BONUS_MG
		} else if filePawns&advancedShieldMask != 0 {
			penalty -= PAWN_SHIELD_ADVANCE_MG
		} else if filePawns == 0 {
			penalty -= PAWN_SHIELD_MISSING_MG
		}
	}

	// Open/semi-open files near king (stacks with pawn shield missing penalty by design).
	for f := startFile; f <= endFile; f++ {
		fileMask := FILE_MASK_BITBOARDS[f]
		if fileMask&friendlyPawns == 0 && fileMask&enemyPawns == 0 {
			penalty -= KING_OPEN_FILE_MG
		} else if fileMask&friendlyPawns == 0 {
			penalty -= KING_SEMI_OPEN_FILE_MG
		}
	}

	return penalty
}
