package main

import (
	"fmt"
	"math/rand"
)

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

// TODO: Make these fancy https://www.chessprogramming.org/Magic_Bitboards#Fancy
var BISHOP_ATTACKS [BOARD_SQUARES][512]Bitboard
var ROOK_ATTACKS [BOARD_SQUARES][4096]Bitboard

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

var ROOK_MAGIC_NUMBERS = [BOARD_SQUARES]Bitboard{
	36028868594761856,
	594475452265930752,
	4755819357035569472,
	72062060812968192,
	144150442711056392,
	2377979772383921152,
	144119588404396544,
	252201891063136512,
	613193237837939232,
	45106433739456578,
	11534563350335738112,
	4758616025560268809,
	7783909083267991812,
	36169620415054336,
	36310314945612544,
	1266638485717120,
	36029071902646272,
	288653139446341632,
	2307110196901580800,
	2816949330395176,
	432628138850157585,
	7350578830039319552,
	27162885041946880,
	2199031715908,
	576531161849799308,
	869195042688901120,
	285396283949120,
	9289056484590625,
	5765733974844639232,
	2201179128832,
	13835075664648634881,
	37295717882036484,
	1152992423123624192,
	2891346213867634688,
	306614489824894976,
	72066392286826496,
	2422962989960274945,
	3377837763530888,
	9532018213697635845,
	10232179454004693124,
	576619787426234368,
	4665730898119983169,
	4507998344978432,
	35665425268744,
	9516387779802628112,
	9386064643256090648,
	18579032364023810,
	11836027213231292417,
	10412323577591234688,
	18015085746245760,
	4504836645585184,
	288318337232937216,
	18023201047314560,
	213305322930304,
	140741787910272,
	9511889471445963264,
	10152104433811489,
	36284957462657,
	13835128973917430274,
	731175500980237,
	4913990281122293762,
	658370005053997097,
	11818669187387097220,
	9223936087395926162,
}

var BISHOP_MAGIC_NUMBERS = [BOARD_SQUARES]Bitboard{
	90780112705880578,
	1130849890600960,
	4579501682000384,
	15771742887089688,
	6342198710732324992,
	9223943800149245956,
	79319525229576,
	564607945147968,
	9629269957113156616,
	125413061951616,
	2454479406292019584,
	576465236393132064,
	2094747866554695682,
	1152925387328585739,
	9259612301152487441,
	9223378636112404996,
	162164839780324352,
	166917014898639872,
	2577358368538640,
	38289410644680706,
	9800395790953349120,
	282575294170128,
	297378316150640640,
	572308964114688,
	4620992320818659584,
	290614873687044,
	2314929373910138960,
	72480390752895008,
	9232381435204096016,
	212206314619424,
	1271061498036736,
	288514325181253632,
	1190094068601062404,
	144260393959686657,
	5207289302532292681,
	5067650172715072,
	5192687754290528390,
	14089150592715400,
	73193398139373568,
	40844662780526856,
	4622949485975324744,
	576532255455840256,
	18858965442494720,
	9223653786982023680,
	288265698642236160,
	9232381608277901824,
	9609903460320384,
	9226206579542458624,
	4973663942830197248,
	4928045509115908,
	9043502600814720,
	11687685459063996417,
	9289086586325024,
	581001752779718960,
	4613941134159511552,
	586075192028200960,
	3171099570124103680,
	160535174254592,
	18018800855224320,
	4630272232445313544,
	148900280213193728,
	9641379395797505,
	9262259632844833536,
	9223979070486822979,
}

var INITIALIZED bool

func InitBitboards() {
	if INITIALIZED {
		return
	}

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
	}

	for i := 0; i < BOARD_SQUARES; i++ {
		bb := Bitboard(1 << i)

		InitPawnAttacksBitboard(i, &bb)
		InitKnightAttacksBitboard(i, &bb)
		InitKingAttacksBitboard(i, &bb)

		InitBishopMasksBitboard(i)
		InitRookMasksBitboard(i)
		InitBishopAttacksBitboard(i)
		InitRookAttacksBitboard(i)
	}

	INITIALIZED = true
}

func InitPawnAttacksBitboard(i int, bb *Bitboard) {
	PAWN_ATTACKS[WHITE][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) << SHIFT_NEG_DIAG
	PAWN_ATTACKS[WHITE][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) << SHIFT_POS_DIAG

	PAWN_ATTACKS[BLACK][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A]) >> SHIFT_POS_DIAG
	PAWN_ATTACKS[BLACK][i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_H]) >> SHIFT_NEG_DIAG
}

func InitKnightAttacksBitboard(i int, bb *Bitboard) {
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A] &
		FILE_CLEAR_BITBOARDS[FILE_B]) << 6
	KNIGHT_ATTACKS[i] |= (*bb & FILE_CLEAR_BITBOARDS[FILE_A] &
		FILE_CLEAR_BITBOARDS[FILE_B]) >> 10

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

func InitBishopAttacksBitboard(i int) {
	mask := BISHOP_MASKS[i]
	shifts := BISHOP_SHIFTS[i]
	for j := 0; j < (1 << shifts); j++ {
		blockers := GetBlockersPermutation(j, shifts, mask)
		magicIndex := (blockers * BISHOP_MAGIC_NUMBERS[i]) >> (BOARD_SQUARES - shifts)
		BISHOP_ATTACKS[i][magicIndex] = GenerateBishopAttacksBitboard(i, blockers)
	}
}

func InitRookAttacksBitboard(i int) {
	mask := ROOK_MASKS[i]
	shifts := ROOK_SHIFTS[i]
	for j := 0; j < (1 << shifts); j++ {
		blockers := GetBlockersPermutation(j, shifts, mask)
		magicIndex := (blockers * ROOK_MAGIC_NUMBERS[i]) >> (BOARD_SQUARES - shifts)
		ROOK_ATTACKS[i][magicIndex] = GenerateRookAttacksBitboard(i, blockers)
	}
}

func BishopAttacks(c Coord, blockers Bitboard) Bitboard {
	magicNumber := ((BISHOP_MASKS[c] & blockers) * BISHOP_MAGIC_NUMBERS[c]) >> (BOARD_SQUARES - BISHOP_SHIFTS[c])
	return BISHOP_ATTACKS[c][magicNumber]
}

func RookAttacks(c Coord, blockers Bitboard) Bitboard {
	magicNumber := ((ROOK_MASKS[c] & blockers) * ROOK_MAGIC_NUMBERS[c]) >> (BOARD_SQUARES - ROOK_SHIFTS[c])
	return ROOK_ATTACKS[c][magicNumber]
}

func FindMagicNumbers() error {
	rand.Seed(181818)

	fmt.Println("Rook Magic Numbers:")
	for i := 0; i < BOARD_SQUARES; i++ {
		mask := ROOK_MASKS[i]
		shifts := ROOK_SHIFTS[i]
		blockers := make([]Bitboard, 1<<shifts)
		attacks := make([]Bitboard, 1<<shifts)
		for j := 0; j < (1 << shifts); j++ {
			blockers[j] = GetBlockersPermutation(j, shifts, mask)
			attacks[j] = GenerateRookAttacksBitboard(i, blockers[j])
		}
		foundMagicNumber := false
		for k := 0; k < 100000000 && !foundMagicNumber; k++ {
			magicNumber := rand.Uint64() & rand.Uint64() & rand.Uint64()
			magicIndexes := make([]Bitboard, 1<<shifts)
			invalid := false
			for j := 0; j < (1<<shifts) && !invalid; j++ {
				magicIndex := (blockers[j] * Bitboard(magicNumber)) >> (BOARD_SQUARES - shifts)
				if magicIndexes[magicIndex] == 0 {
					magicIndexes[magicIndex] = attacks[j]
				} else if magicIndexes[magicIndex] != attacks[j] {
					invalid = true
				}
			}
			if !invalid {
				fmt.Printf("%v,\n", magicNumber)
				foundMagicNumber = true
			}
		}
		if !foundMagicNumber {
			return fmt.Errorf("unable to find magic number for rook at coord %v", COORD_MAP[i])
		}
	}

	fmt.Println("Bishop Magic Numbers:")
	for i := 0; i < BOARD_SQUARES; i++ {
		mask := BISHOP_MASKS[i]
		shifts := BISHOP_SHIFTS[i]
		blockers := make([]Bitboard, 1<<shifts)
		attacks := make([]Bitboard, 1<<shifts)
		for j := 0; j < (1 << shifts); j++ {
			blockers[j] = GetBlockersPermutation(j, shifts, mask)
			attacks[j] = GenerateBishopAttacksBitboard(i, blockers[j])
		}
		foundMagicNumber := false
		for k := 0; k < 100000000 && !foundMagicNumber; k++ {
			magicNumber := rand.Uint64() & rand.Uint64() & rand.Uint64()
			magicIndexes := make([]Bitboard, 1<<shifts)
			invalid := false
			for j := 0; j < (1<<shifts) && !invalid; j++ {
				magicIndex := (blockers[j] * Bitboard(magicNumber)) >> (BOARD_SQUARES - shifts)
				if magicIndexes[magicIndex] == 0 {
					magicIndexes[magicIndex] = attacks[j]
				} else if magicIndexes[magicIndex] != attacks[j] {
					invalid = true
				}
			}
			if !invalid {
				fmt.Printf("%v,\n", magicNumber)
				foundMagicNumber = true
			}
		}
		if !foundMagicNumber {
			return fmt.Errorf("unable to find magic number for bishop at coord %v", COORD_MAP[i])
		}
	}

	return nil
}

func GetBlockersPermutation(index int, count int, mask Bitboard) Bitboard {
	var blockers Bitboard

	for i := 0; i < count; i++ {
		lsbIndex := mask.PopLSB()
		if index&(1<<i) > 0 {
			blockers |= Bitboard(1 << lsbIndex)
		}
	}

	return blockers
}
