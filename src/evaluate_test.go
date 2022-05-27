package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		fen      string
		expected int
	}{
		{
			STARTING_FEN,
			0,
		},
		{
			"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			-30,
		},
		{
			"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			30,
		},
		{
			"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			0,
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 4 4",
			0,
		},
		{
			"r1bqk1nr/ppppbppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
			-15,
		},
		{
			"r1bq1rk1/ppppbppp/2n2n2/1B2p3/4P3/3P1N2/PPP2PPP/RNBQ1RK1 w - - 1 6",
			-10,
		},
	}
	for _, tt := range tests {
		board, err := NewBoard(tt.fen)
		if err != nil {
			t.Error(err)
		}
		actual := board.Evaluate()

		if actual != tt.expected {
			t.Errorf("score: %v != %v\n\n%v", actual, tt.expected, board.ToString())
		}
	}
}
