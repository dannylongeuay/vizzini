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
		{STARTING_FEN, 181818, 14949957711303876919},
		{STARTING_FEN, 23456789, 15652837511423572584},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			181818, 1902633635645420733},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			181818, 15174481012482708289},
	}
	for _, tt := range tests {
		InitHashKeys(tt.seed)
		BOARD_INITIALIZED = true
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
			14949957711303876919,
			2557449784889405518,
		},
		{STARTING_FEN, 181818, BLACK_ROOK, A8,
			14949957711303876919,
			8284552276674073387,
		},
	}
	for _, tt := range tests {
		InitHashKeys(tt.seed)
		BOARD_INITIALIZED = true
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
		b.HashSquare(tt.square, tt.coord)

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.HashSquare(tt.square, tt.coord)

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
			14949957711303876919,
			12064830828428678401,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1", 181818,
			12064830828428678401,
			14949957711303876919,
		},
	}
	for _, tt := range tests {
		InitHashKeys(tt.seed)
		BOARD_INITIALIZED = true
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}

		b.HashSide()

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.HashSide()

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
			14949957711303876919,
			14877183630399058325,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kkq - 0 1", 181818,
			13184813557952630284,
			14877183630399058325,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1", 181818,
			2644723196783502706,
			14877183630399058325,
		},
	}
	for _, tt := range tests {
		InitHashKeys(tt.seed)
		BOARD_INITIALIZED = true
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}

		b.HashCastling()

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.HashCastling()

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
			15174481012482708289,
			1902633635645420733,
		},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq c6 0 1", 181818,
			6985219659505869312,
			1902633635645420733,
		},
	}
	for _, tt := range tests {
		InitHashKeys(tt.seed)
		BOARD_INITIALIZED = true
		b, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}

		b.HashEnPassant()

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.HashEnPassant()

		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
	}
}
