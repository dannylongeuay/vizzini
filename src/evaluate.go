package main

// Phase weights per piece type (non-pawn, non-king).
const PHASE_KNIGHT = 1
const PHASE_BISHOP = 1
const PHASE_ROOK = 2
const PHASE_QUEEN = 4

// Total phase at game start: 4 knights + 4 bishops + 4 rooks + 2 queens = 24.
const PHASE_TOTAL = 24

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
