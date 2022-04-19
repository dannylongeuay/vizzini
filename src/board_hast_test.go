package main

import (
	"testing"
)

func TestGenerateBoardHash(t *testing.T) {
	tests := []struct {
		fen      string
		seed     int64
		expected uint64
	}{
		{STARTING_FEN, 181818, 2353279852124060199},
		{STARTING_FEN, 23456789, 11159194169277155965},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			181818, 6612989906737374792},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			181818, 11425164701060763039},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
			181818, 10669042118094424949},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			181818, 5470054360002429070},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR b KQkq c6 0 2",
			181818, 12401481765782589279},
		{"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w Qkq c6 0 2",
			181818, 17721979826370444963},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := newBoard(tt.fen)
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
		fen         string
		seed        int64
		square      Square
		squareCoord SquareCoord
		initial     uint64
		hashed      uint64
	}{
		{STARTING_FEN, 181818, WHITE_PAWN, "e2",
			2353279852124060199,
			11656825373806580966,
		},
		{STARTING_FEN, 181818, BLACK_ROOK, "a8",
			2353279852124060199,
			1042891947479458529,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		b.hashSquare(tt.square, squareIndex)

		if b.hash != tt.hashed {
			t.Errorf("board hash: %v != %v", b.hash, tt.hashed)
		}

		b.hashSquare(tt.square, squareIndex)

		if b.hash != tt.initial {
			t.Errorf("board hash: %v != %v", b.hash, tt.initial)
		}
	}
}

func TestHashSide(t *testing.T) {
	tests := []struct {
		fen     string
		seed    int64
		initial uint64
		hashed  uint64
	}{
		{STARTING_FEN, 181818,
			2353279852124060199,
			1775090316340608525,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1", 181818,
			1775090316340608525,
			2353279852124060199,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := newBoard(tt.fen)
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
		initial uint64
		hashed  uint64
	}{
		{STARTING_FEN, 181818,
			2353279852124060199,
			4201107985908305079,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kkq - 0 1", 181818,
			5286067747604221733,
			4201107985908305079,
		},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1", 181818,
			10616819202386075795,
			4201107985908305079,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := newBoard(tt.fen)
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
		initial uint64
		hashed  uint64
	}{
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", 181818,
			11425164701060763039,
			6612989906737374792,
		},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq c6 0 1", 181818,
			9529488694904729011,
			6612989906737374792,
		},
	}
	for _, tt := range tests {
		seedKeys(tt.seed)
		b, err := newBoard(tt.fen)
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
