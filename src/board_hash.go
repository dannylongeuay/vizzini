package main

import (
	"math/rand"
)

var SquareKeys [SQUARE_TYPES][BOARD_SQUARES]Hash
var SideKey Hash
var CastleKeys [CASTLING_RIGHTS_PERMUTATIONS]Hash

func InitHashKeys(seed int64) {
	rand.Seed(seed)

	for x := 0; x < SQUARE_TYPES; x++ {
		for y := 0; y < BOARD_SQUARES; y++ {
			SquareKeys[x][y] = Hash(rand.Uint64())
		}
	}

	SideKey = Hash(rand.Uint64())

	for i := 0; i < CASTLING_RIGHTS_PERMUTATIONS; i++ {
		CastleKeys[i] = Hash(rand.Uint64())
	}
}

func (b *Board) HashSquare(square Square, coord Coord) {
	b.hash ^= SquareKeys[square][coord]
}

func (b *Board) HashSide() {
	b.hash ^= SideKey
}

func (b *Board) HashEnPassant() {
	b.hash ^= SquareKeys[EMPTY][b.epCoord]
}

func (b *Board) HashCastling() {
	b.hash ^= CastleKeys[b.castleRights]
}

func (b *Board) GenerateBoardHash() {
	b.hash = 0
	for i, square := range b.squares {
		if square == EMPTY {
			continue
		}
		b.HashSquare(square, Coord(i))
	}
	if b.sideToMove == WHITE {
		b.HashSide()
	}
	if b.epCoord != A1 {
		b.HashEnPassant()
	}
	b.HashCastling()
}
