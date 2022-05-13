package main

type Square uint8
type Color uint8
type File uint8
type Rank uint8
type CastleRights uint8
type MoveKind uint8
type Coord uint64

/*
	Move
	00000000       000000        000000         0000           0000     0000
	 unused		origin coord	dst coord	origin square	dst square	kind
*/
type Move uint32
type Bitboard uint64
type Hash uint64
