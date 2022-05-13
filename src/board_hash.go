package main

import (
	"math/rand"
)

var SquareKeys [SQUARE_TYPES][BOARD_SQUARES]Hash
var SideKey Hash
var CastleKeys [CASTLING_RIGHTS_PERMUTATIONS]Hash

func seedKeys(seed int64) {
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

func (b *Board) hashSquare(square Square, coord Coord) {
	b.hash ^= SquareKeys[square][coord]
}

func (b *Board) hashSide() {
	b.hash ^= SideKey
}

func (b *Board) hashEnPassant() {
	b.hash ^= SquareKeys[EMPTY][b.epCoord]
}

func (b *Board) hashCastling() {
	b.hash ^= CastleKeys[b.castleRights]
}

func (b *Board) generateBoardHash() {
	b.hash = 0
	for i, square := range b.squares {
		if square == INVALID || square == EMPTY {
			continue
		}
		b.hashSquare(square, Coord(i))
	}
	if b.sideToMove == WHITE {
		b.hashSide()
	}
	if b.epCoord != 0 {
		b.hashEnPassant()
	}
	b.hashCastling()
}
