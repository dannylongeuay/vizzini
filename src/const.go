package main

const STARTING_FEN string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

const PLAYERS int = 2
const BOARD_SQUARES int = 64
const RANKS int = 8
const FILES int = 8
const SQUARE_TYPES int = 13
const MAX_GAME_MOVES uint16 = 300
const MOVE_KINDS int = 15

const CASTLING_RIGHTS_PERMUTATIONS int = 16
const CASTLING_RIGHTS_BLACK_QUEEN_MASK CastleRights = 1
const CASTLING_RIGHTS_BLACK_KING_MASK CastleRights = 2
const CASTLING_RIGHTS_WHITE_QUEEN_MASK CastleRights = 4
const CASTLING_RIGHTS_WHITE_KING_MASK CastleRights = 8

const MOVE_ORDER_MASK Move = 4278190080
const MOVE_ORDER_SHIFT int = 24
const MOVE_ORIGIN_COORD_MASK Move = 16515072
const MOVE_ORIGIN_COORD_SHIFT int = 18
const MOVE_DST_COORD_MASK Move = 258048
const MOVE_DST_COORD_SHIFT int = 12
const MOVE_ORIGIN_SQUARE_MASK Move = 3840
const MOVE_ORIGIN_SQUARE_SHIFT int = 8
const MOVE_DST_SQUARE_MASK Move = 240
const MOVE_DST_SQUARE_SHIFT int = 4
const MOVE_KIND_MASK Move = 15

const UNDO_EP_COORD_MASK Undo = 4433230883192832
const UNDO_EP_COORD_SHIFT int = 46
const UNDO_CASTLE_RIGHTS_MASK Undo = 65970697666560
const UNDO_CASTLE_RIGHTS_SHIFT int = 42
const UNDO_HALF_MOVE_MASK Undo = 4329327034368
const UNDO_HALF_MOVE_SHIFT int = 36
const UNDO_CLEAR_SQUARE_MASK Undo = 64424509440
const UNDO_CLEAR_SQUARE_SHIFT int = 32
const UNDO_MOVE_MASK Undo = 4294967295

const KING_MOVE_RANGE int = 2
const KING_CASTLE_MOVE_DIST = 2

const SHIFT_NEG_DIAG int = 7
const SHIFT_VERTICAL int = 8
const SHIFT_POS_DIAG int = 9
const SHIFT_HORIZONTAL int = 1

const INITIAL_MOVES_CAPACITY int = 32
const MAX_PAWN_MOVES int = 12
const MAX_KNIGHT_MOVES int = 8
const MAX_BISHOP_MOVES int = 13
const MAX_ROOK_MOVES int = 14
const MAX_QUEEN_MOVES int = 27
const MAX_KING_MOVES int = 8

const MAX_SQUARE_ATTACKERS int = 16
const MAX_SQUARE_KNIGHT_ATTACKERS int = 8
const MAX_SQUARE_DIAGONAL_ATTACKERS int = 4
const MAX_SQUARE_CARDINAL_ATTACKERS int = 4

const DEFAULT_MAX_DEPTH int = 64
const DEFAULT_MAX_NODES int = 0

const SEARCH_BUFFER int64 = 50
const PV_TABLE_SIZE Hash = 1000000
const KILLERS_SIZE = 2
const KILLERS_DEPTH = 64

const DRAW int = 0
const MAX_SCORE int = 100000
const MIN_SCORE int = -100000

const MVV_LVA_EN_PASSANT MoveOrder = 15
const MVV_LVA_KNIGHT_PROMOTION MoveOrder = 101
const MVV_LVA_BISHOP_PROMOTION MoveOrder = 102
const MVV_LVA_ROOK_PROMOTION MoveOrder = 103
const MVV_LVA_QUEEN_PROMOTION MoveOrder = 104
const MVV_LVA_KNIGHT_PROMOTION_CAPTURE MoveOrder = 201
const MVV_LVA_BISHOP_PROMOTION_CAPTURE MoveOrder = 202
const MVV_LVA_ROOK_PROMOTION_CAPTURE MoveOrder = 203
const MVV_LVA_QUEEN_PROMOTION_CAPTURE MoveOrder = 204

const (
	EMPTY Square = iota
	WHITE_PAWN
	WHITE_KNIGHT
	WHITE_BISHOP
	WHITE_ROOK
	WHITE_QUEEN
	WHITE_KING
	BLACK_PAWN
	BLACK_KNIGHT
	BLACK_BISHOP
	BLACK_ROOK
	BLACK_QUEEN
	BLACK_KING
)

const (
	WHITE Color = iota
	BLACK
	COLOR_NONE
)

const (
	FILE_A File = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
	FILE_NONE
)

const (
	RANK_ONE Rank = iota
	RANK_TWO
	RANK_THREE
	RANK_FOUR
	RANK_FIVE
	RANK_SIX
	RANK_SEVEN
	RANK_EIGHT
	RANK_NONE
)

const (
	QUIET MoveKind = iota
	DOUBLE_PAWN_PUSH
	KING_CASTLE
	QUEEN_CASTLE
	CAPTURE
	EP_CAPTURE
	KNIGHT_PROMOTION
	BISHOP_PROMOTION
	ROOK_PROMOTION
	QUEEN_PROMOTION
	KNIGHT_PROMOTION_CAPTURE
	BISHOP_PROMOTION_CAPTURE
	ROOK_PROMOTION_CAPTURE
	QUEEN_PROMOTION_CAPTURE
)

const (
	A1 Coord = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)
