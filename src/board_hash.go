package main

import (
	"math/rand"
)

var SquareKeys [14][120]uint64
var SideKey uint64
var CastleKeys [16]uint64

func seedKeys(seed int64) {
	rand.Seed(seed)

	for x := 0; x < 14; x++ {
		for y := 0; y < 120; y++ {
			SquareKeys[x][y] = rand.Uint64()
		}
	}

	SideKey = rand.Uint64()

	for i := 0; i < 16; i++ {
		CastleKeys[i] = rand.Uint64()
	}
}

func (b *board) hashSquare(square Square, squareIndex SquareIndex) {
	b.hash ^= SquareKeys[square][squareIndex]
}

func (b *board) hashSide() {
	b.hash ^= SideKey
}

func (b *board) hashEnPassant() {
	b.hash ^= SquareKeys[EMPTY][b.epIndex]
}

func (b *board) hashCastling() {
	b.hash ^= CastleKeys[b.castleRights]
}

func (b *board) generateBoardHash() {
	b.hash = 0
	for squareIndex, square := range b.squares {
		if square == INVALID || square == EMPTY {
			continue
		}
		b.hashSquare(square, SquareIndex(squareIndex))
	}
	if b.sideToMove == WHITE {
		b.hashSide()
	}
	if b.epIndex != 0 {
		b.hashEnPassant()
	}
	b.hashCastling()
}
