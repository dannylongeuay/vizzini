package main

import (
	"testing"
)

func TestGenerateBoardHash(t *testing.T) {
	tests := []struct {
		fen      string
		seed     int64
		expected Hash
	}{
		{STARTING_FEN, 181818, 2564669793889550577},
		{STARTING_FEN, 23456789, 8789488460745541294},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			181818, 4311714675655613366},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			181818, 4998122769579152046},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.expected {
			t.Errorf("board hash: %v != %v", b.hash, tt.expected)
		}
	}
}

func TestHashSquare(t *testing.T) {
	tests := []struct {
		fen     string
		seed    int64
		square  Square
		coord   Coord
		initial Hash
		hashed  Hash
	}{
		{STARTING_FEN, 181818, WHITE_PAWN, E2,
			2564669793889550577,
			14957159999261659528,
		},
		{STARTING_FEN, 181818, BLACK_ROOK, A8,
			2564669793889550577,
			11391792193050571501,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
		b.hashSquare(tt.square, tt.coord)

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.hashSquare(tt.square, tt.coord)

		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
	}
}

func TestHashSide(t *testing.T) {
	tests := []struct {
		fen     string
		seed    int64
		initial Hash
		hashed  Hash
	}{
		{STARTING_FEN, 181818,
			2564669793889550577,
			9718531765999254026,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1", 181818,
			9718531765999254026,
			2564669793889550577,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}

		b.hashSide()

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.hashSide()

		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
	}
}

func TestHashCastling(t *testing.T) {
	tests := []struct {
		fen     string
		seed    int64
		initial Hash
		hashed  Hash
	}{
		{STARTING_FEN, 181818,
			2564669793889550577,
			227704336027123544,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kkq - 0 1", 181818,
			3558946505462474442,
			227704336027123544,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1", 181818,
			1422002867735741380,
			227704336027123544,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}

		b.hashCastling()

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.hashCastling()

		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
	}
}

func TestHashEnPassant(t *testing.T) {
	tests := []struct {
		fen     string
		seed    int64
		initial Hash
		hashed  Hash
	}{
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", 181818,
			4998122769579152046,
			4311714675655613366,
		},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq c6 0 1", 181818,
			14262529468124065421,
			4311714675655613366,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}

		b.hashEnPassant()

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.hashEnPassant()

		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
	}
}
