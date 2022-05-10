package main

import (
	"testing"
)

func TestSquareKnightAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2", BLACK, "e5", 1,
			[]SquareCoord{"f3"},
		},
		{"k7/8/2n1n3/1n3n2/3Q4/1n3n2/2n1n3/7K w - - 0 1", WHITE, "d4", MAX_SQUARE_KNIGHT_ATTACKERS,
			[]SquareCoord{"b3", "b5", "c2", "c6", "e2", "e6", "f3", "f5"},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		attackers := make(map[SquareIndex]bool, MAX_SQUARE_KNIGHT_ATTACKERS)
		b.squareKnightAttackers(&attackers, tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			_, present := attackers[attackerIndex]
			if !present {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}

func TestSquareDiagonalAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3", BLACK, "c6", 1,
			[]SquareCoord{"b5"},
		},
		{"k7/8/5Q2/2K5/3n4/8/5P2/B7 b - - 0 1", BLACK, "d4", 3,
			[]SquareCoord{"a1", "c5", "f6"},
		},
		{"k6B/8/8/2P5/3q4/4B3/1Q6/7K b - - 0 1", BLACK, "d4", 3,
			[]SquareCoord{"b2", "e3", "h8"},
		},
		{"k6B/8/8/2B5/3q4/4P3/1Q6/7K b - - 0 1", BLACK, "d4", MAX_SQUARE_DIAGONAL_ATTACKERS,
			[]SquareCoord{"b2", "c5", "e3", "h8"},
		},
		{"rn1q2nr/p2kb1pp/1p6/4p1p1/4K3/8/PPPP1P1P/RNBQ1BNR w - - 0 10", WHITE, "d4", 1,
			[]SquareCoord{"e5"},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		attackers := make(map[SquareIndex]bool, MAX_SQUARE_DIAGONAL_ATTACKERS)
		b.squareDiagonalAttackers(&attackers, tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			_, present := attackers[attackerIndex]
			if !present {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}

func TestSquareCardinalAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"k7/3Q4/8/8/1R1p1K2/8/3R4/8 b - - 0 1", BLACK, "d4", 3,
			[]SquareCoord{"b4", "d2", "d7"},
		},
		{"k7/3Q4/8/8/1R1pK3/8/3R4/8 b - - 0 1", BLACK, "d4", MAX_SQUARE_CARDINAL_ATTACKERS,
			[]SquareCoord{"b4", "d2", "d7", "e4"},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		attackers := make(map[SquareIndex]bool, MAX_SQUARE_CARDINAL_ATTACKERS)
		b.squareCardinalAttackers(&attackers, tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			_, present := attackers[attackerIndex]
			if !present {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}

func TestSquareAttackers(t *testing.T) {
	tests := []struct {
		fen             string
		side            Color
		squareCoord     SquareCoord
		attackersLength int
		attackerCoords  []SquareCoord
	}{
		{STARTING_FEN, WHITE, "a1", 0,
			[]SquareCoord{},
		},
		{"k7/8/5Q2/2K5/3n2R1/4P3/2N2P2/B7 b - - 0 1", BLACK, "d4", 6,
			[]SquareCoord{"a1", "c2", "c5", "e3", "f6", "g4"},
		},
		{"k7/8/2N1NQ2/1NBK1N2/1R1n2R1/1N1QPN2/2N1NP2/B7 b - - 0 1", BLACK, "d4", MAX_SQUARE_ATTACKERS,
			[]SquareCoord{
				"a1", "b3", "b4", "b5", "c2", "c5", "c6", "d3",
				"d5", "e2", "e3", "e6", "f3", "f5", "f6", "g4",
			},
		},
	}
	for _, tt := range tests {
		b, err := newBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		squareIndex, err := squareIndexByCoord(tt.squareCoord)
		if err != nil {
			t.Error(err)
		}
		attackers := make(map[SquareIndex]bool, MAX_SQUARE_ATTACKERS)
		b.squareAttackers(&attackers, tt.side, squareIndex)
		if len(attackers) != tt.attackersLength {
			t.Errorf("attackers length: %v != %v", len(attackers), tt.attackersLength)
		}
		for _, attackerCoord := range tt.attackerCoords {
			attackerIndex, err := squareIndexByCoord(attackerCoord)
			if err != nil {
				t.Error(err)
			}
			_, present := attackers[attackerIndex]
			if !present {
				t.Errorf("attacker index %v not found in %v", attackerIndex, attackers)
			}
		}
	}
}
