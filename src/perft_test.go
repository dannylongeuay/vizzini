package main

import (
	"testing"
)

type DepthCount struct {
	depth int
	count int
}

func Perft(b *Board, depth int) int {
	var nodes int

	if depth == 0 {
		return 1
	}

	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	b.GenerateMoves(&moves, b.sideToMove)
	for _, m := range moves {
		err := b.MakeMove(m)
		if err == nil {
			nodes += Perft(b, depth-1)
		}
		b.UndoMove()
	}

	return nodes
}

func TestPerft(t *testing.T) {
	tests := []struct {
		fen         string
		depthCounts []DepthCount
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			[]DepthCount{{1, 20}, {2, 400}, {3, 8902}, {4, 197281}, {5, 4865609}, {6, 119060324}},
		},
		{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			[]DepthCount{{1, 48}, {2, 2039}, {3, 97862}, {4, 4085603}, {5, 193690690}},
		},
		{"4k3/8/8/8/8/8/8/4K2R w K - 0 1",
			[]DepthCount{{1, 15}, {2, 66}, {3, 1197}, {4, 7059}, {5, 133987}, {6, 764643}},
		},
		{"4k3/8/8/8/8/8/8/R3K3 w Q - 0 1",
			[]DepthCount{{1, 16}, {2, 71}, {3, 1287}, {4, 7626}, {5, 145232}, {6, 846648}},
		},
		{"4k2r/8/8/8/8/8/8/4K3 w k - 0 1",
			[]DepthCount{{1, 5}, {2, 75}, {3, 459}, {4, 8290}, {5, 47635}, {6, 899442}},
		},
		{"r3k3/8/8/8/8/8/8/4K3 w q - 0 1",
			[]DepthCount{{1, 5}, {2, 80}, {3, 493}, {4, 8897}, {5, 52710}, {6, 1001523}},
		},
		{"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			[]DepthCount{{1, 26}, {2, 112}, {3, 3189}, {4, 17945}, {5, 532933}, {6, 2788982}},
		},
		{"r3k2r/8/8/8/8/8/8/4K3 w kq - 0 1",
			[]DepthCount{{1, 5}, {2, 130}, {3, 782}, {4, 22180}, {5, 118882}, {6, 3517770}},
		},
		{"8/8/8/8/8/8/6k1/4K2R w K - 0 1",
			[]DepthCount{{1, 12}, {2, 38}, {3, 564}, {4, 2219}, {5, 37735}, {6, 185867}},
		},
		{"8/8/8/8/8/8/1k6/R3K3 w Q - 0 1",
			[]DepthCount{{1, 15}, {2, 65}, {3, 1018}, {4, 4573}, {5, 80619}, {6, 413018}},
		},
		{"4k2r/6K1/8/8/8/8/8/8 w k - 0 1",
			[]DepthCount{{1, 3}, {2, 32}, {3, 134}, {4, 2073}, {5, 10485}, {6, 179869}},
		},
		{"r3k3/1K6/8/8/8/8/8/8 w q - 0 1",
			[]DepthCount{{1, 4}, {2, 49}, {3, 243}, {4, 3991}, {5, 20780}, {6, 367724}},
		},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			[]DepthCount{{1, 26}, {2, 568}, {3, 13744}, {4, 314346}, {5, 7594526}, {6, 179862938}},
		},
		{"r3k2r/8/8/8/8/8/8/1R2K2R w Kkq - 0 1",
			[]DepthCount{{1, 25}, {2, 567}, {3, 14095}, {4, 328965}, {5, 8153719}, {6, 195629489}},
		},
		{"r3k2r/8/8/8/8/8/8/2R1K2R w Kkq - 0 1",
			[]DepthCount{{1, 25}, {2, 548}, {3, 13502}, {4, 312835}, {5, 7736373}, {6, 184411439}},
		},
		{"r3k2r/8/8/8/8/8/8/R3K1R1 w Qkq - 0 1",
			[]DepthCount{{1, 25}, {2, 547}, {3, 13579}, {4, 316214}, {5, 7878456}, {6, 189224276}},
		},
		{"1r2k2r/8/8/8/8/8/8/R3K2R w KQk - 0 1",
			[]DepthCount{{1, 26}, {2, 583}, {3, 14252}, {4, 334705}, {5, 8198901}, {6, 198328929}},
		},
		{"2r1k2r/8/8/8/8/8/8/R3K2R w KQk - 0 1",
			[]DepthCount{{1, 25}, {2, 560}, {3, 13592}, {4, 317324}, {5, 7710115}, {6, 185959088}},
		},
		{"r3k1r1/8/8/8/8/8/8/R3K2R w KQq - 0 1",
			[]DepthCount{{1, 25}, {2, 560}, {3, 13607}, {4, 320792}, {5, 7848606}, {6, 190755813}},
		},
		{"4k3/8/8/8/8/8/8/4K2R b K - 0 1",
			[]DepthCount{{1, 5}, {2, 75}, {3, 459}, {4, 8290}, {5, 47635}, {6, 899442}},
		},
		{"4k3/8/8/8/8/8/8/R3K3 b Q - 0 1",
			[]DepthCount{{1, 5}, {2, 80}, {3, 493}, {4, 8897}, {5, 52710}, {6, 1001523}},
		},
		{"4k2r/8/8/8/8/8/8/4K3 b k - 0 1",
			[]DepthCount{{1, 15}, {2, 66}, {3, 1197}, {4, 7059}, {5, 133987}, {6, 764643}},
		},
		{"r3k3/8/8/8/8/8/8/4K3 b q - 0 1",
			[]DepthCount{{1, 16}, {2, 71}, {3, 1287}, {4, 7626}, {5, 145232}, {6, 846648}},
		},
		{"4k3/8/8/8/8/8/8/R3K2R b KQ - 0 1",
			[]DepthCount{{1, 5}, {2, 130}, {3, 782}, {4, 22180}, {5, 118882}, {6, 3517770}},
		},
		{"r3k2r/8/8/8/8/8/8/4K3 b kq - 0 1",
			[]DepthCount{{1, 26}, {2, 112}, {3, 3189}, {4, 17945}, {5, 532933}, {6, 2788982}},
		},
		{"8/8/8/8/8/8/6k1/4K2R b K - 0 1",
			[]DepthCount{{1, 3}, {2, 32}, {3, 134}, {4, 2073}, {5, 10485}, {6, 179869}},
		},
		{"8/8/8/8/8/8/1k6/R3K3 b Q - 0 1",
			[]DepthCount{{1, 4}, {2, 49}, {3, 243}, {4, 3991}, {5, 20780}, {6, 367724}},
		},
		{"4k2r/6K1/8/8/8/8/8/8 b k - 0 1",
			[]DepthCount{{1, 12}, {2, 38}, {3, 564}, {4, 2219}, {5, 37735}, {6, 185867}},
		},
		{"r3k3/1K6/8/8/8/8/8/8 b q - 0 1",
			[]DepthCount{{1, 15}, {2, 65}, {3, 1018}, {4, 4573}, {5, 80619}, {6, 413018}},
		},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			[]DepthCount{{1, 26}, {2, 568}, {3, 13744}, {4, 314346}, {5, 7594526}, {6, 179862938}},
		},
		{"r3k2r/8/8/8/8/8/8/1R2K2R b Kkq - 0 1",
			[]DepthCount{{1, 26}, {2, 583}, {3, 14252}, {4, 334705}, {5, 8198901}, {6, 198328929}},
		},
		{"r3k2r/8/8/8/8/8/8/2R1K2R b Kkq - 0 1",
			[]DepthCount{{1, 25}, {2, 560}, {3, 13592}, {4, 317324}, {5, 7710115}, {6, 185959088}},
		},
		{"r3k2r/8/8/8/8/8/8/R3K1R1 b Qkq - 0 1",
			[]DepthCount{{1, 25}, {2, 560}, {3, 13607}, {4, 320792}, {5, 7848606}, {6, 190755813}},
		},
		{"1r2k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1",
			[]DepthCount{{1, 25}, {2, 567}, {3, 14095}, {4, 328965}, {5, 8153719}, {6, 195629489}},
		},
		{"2r1k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1",
			[]DepthCount{{1, 25}, {2, 548}, {3, 13502}, {4, 312835}, {5, 7736373}, {6, 184411439}},
		},
		{"r3k1r1/8/8/8/8/8/8/R3K2R b KQq - 0 1",
			[]DepthCount{{1, 25}, {2, 547}, {3, 13579}, {4, 316214}, {5, 7878456}, {6, 189224276}},
		},
		{"8/1n4N1/2k5/8/8/5K2/1N4n1/8 w - - 0 1",
			[]DepthCount{{1, 14}, {2, 195}, {3, 2760}, {4, 38675}, {5, 570726}, {6, 8107539}},
		},
		{"8/1k6/8/5N2/8/4n3/8/2K5 w - - 0 1",
			[]DepthCount{{1, 11}, {2, 156}, {3, 1636}, {4, 20534}, {5, 223507}, {6, 2594412}},
		},
		{"8/8/4k3/3Nn3/3nN3/4K3/8/8 w - - 0 1",
			[]DepthCount{{1, 19}, {2, 289}, {3, 4442}, {4, 73584}, {5, 1198299}, {6, 19870403}},
		},
		{"K7/8/2n5/1n6/8/8/8/k6N w - - 0 1",
			[]DepthCount{{1, 3}, {2, 51}, {3, 345}, {4, 5301}, {5, 38348}, {6, 588695}},
		},
		{"k7/8/2N5/1N6/8/8/8/K6n w - - 0 1",
			[]DepthCount{{1, 17}, {2, 54}, {3, 835}, {4, 5910}, {5, 92250}, {6, 688780}},
		},
		{"8/1n4N1/2k5/8/8/5K2/1N4n1/8 b - - 0 1",
			[]DepthCount{{1, 15}, {2, 193}, {3, 2816}, {4, 40039}, {5, 582642}, {6, 8503277}},
		},
		{"8/1k6/8/5N2/8/4n3/8/2K5 b - - 0 1",
			[]DepthCount{{1, 16}, {2, 180}, {3, 2290}, {4, 24640}, {5, 288141}, {6, 3147566}},
		},
		{"8/8/3K4/3Nn3/3nN3/4k3/8/8 b - - 0 1",
			[]DepthCount{{1, 4}, {2, 68}, {3, 1118}, {4, 16199}, {5, 281190}, {6, 4405103}},
		},
		{"K7/8/2n5/1n6/8/8/8/k6N b - - 0 1",
			[]DepthCount{{1, 17}, {2, 54}, {3, 835}, {4, 5910}, {5, 92250}, {6, 688780}},
		},
		{"k7/8/2N5/1N6/8/8/8/K6n b - - 0 1",
			[]DepthCount{{1, 3}, {2, 51}, {3, 345}, {4, 5301}, {5, 38348}, {6, 588695}},
		},
		{"B6b/8/8/8/2K5/4k3/8/b6B w - - 0 1",
			[]DepthCount{{1, 17}, {2, 278}, {3, 4607}, {4, 76778}, {5, 1320507}, {6, 22823890}},
		},
		{"8/8/1B6/7b/7k/8/2B1b3/7K w - - 0 1",
			[]DepthCount{{1, 21}, {2, 316}, {3, 5744}, {4, 93338}, {5, 1713368}, {6, 28861171}},
		},
		{"k7/B7/1B6/1B6/8/8/8/K6b w - - 0 1",
			[]DepthCount{{1, 21}, {2, 144}, {3, 3242}, {4, 32955}, {5, 787524}, {6, 7881673}},
		},
		{"K7/b7/1b6/1b6/8/8/8/k6B w - - 0 1",
			[]DepthCount{{1, 7}, {2, 143}, {3, 1416}, {4, 31787}, {5, 310862}, {6, 7382896}},
		},
		{"B6b/8/8/8/2K5/5k2/8/b6B b - - 0 1",
			[]DepthCount{{1, 6}, {2, 106}, {3, 1829}, {4, 31151}, {5, 530585}, {6, 9250746}},
		},
		{"8/8/1B6/7b/7k/8/2B1b3/7K b - - 0 1",
			[]DepthCount{{1, 17}, {2, 309}, {3, 5133}, {4, 93603}, {5, 1591064}, {6, 29027891}},
		},
		{"k7/B7/1B6/1B6/8/8/8/K6b b - - 0 1",
			[]DepthCount{{1, 7}, {2, 143}, {3, 1416}, {4, 31787}, {5, 310862}, {6, 7382896}},
		},
		{"K7/b7/1b6/1b6/8/8/8/k6B b - - 0 1",
			[]DepthCount{{1, 21}, {2, 144}, {3, 3242}, {4, 32955}, {5, 787524}, {6, 7881673}},
		},
		{"7k/RR6/8/8/8/8/rr6/7K w - - 0 1",
			[]DepthCount{{1, 19}, {2, 275}, {3, 5300}, {4, 104342}, {5, 2161211}, {6, 44956585}},
		},
		{"R6r/8/8/2K5/5k2/8/8/r6R w - - 0 1",
			[]DepthCount{{1, 36}, {2, 1027}, {3, 29215}, {4, 771461}, {5, 20506480}, {6, 525169084}},
		},
		{"7k/RR6/8/8/8/8/rr6/7K b - - 0 1",
			[]DepthCount{{1, 19}, {2, 275}, {3, 5300}, {4, 104342}, {5, 2161211}, {6, 44956585}},
		},
		{"R6r/8/8/2K5/5k2/8/8/r6R b - - 0 1",
			[]DepthCount{{1, 36}, {2, 1027}, {3, 29227}, {4, 771368}, {5, 20521342}, {6, 524966748}},
		},
		{"6kq/8/8/8/8/8/8/7K w - - 0 1",
			[]DepthCount{{1, 2}, {2, 36}, {3, 143}, {4, 3637}, {5, 14893}, {6, 391507}},
		},
		{"6KQ/8/8/8/8/8/8/7k b - - 0 1",
			[]DepthCount{{1, 2}, {2, 36}, {3, 143}, {4, 3637}, {5, 14893}, {6, 391507}},
		},
		{"K7/8/8/3Q4/4q3/8/8/7k w - - 0 1",
			[]DepthCount{{1, 6}, {2, 35}, {3, 495}, {4, 8349}, {5, 166741}, {6, 3370175}},
		},
		{"6qk/8/8/8/8/8/8/7K b - - 0 1",
			[]DepthCount{{1, 22}, {2, 43}, {3, 1015}, {4, 4167}, {5, 105749}, {6, 419369}},
		},
		{"6KQ/8/8/8/8/8/8/7k b - - 0 1",
			[]DepthCount{{1, 2}, {2, 36}, {3, 143}, {4, 3637}, {5, 14893}, {6, 391507}},
		},
		{"K7/8/8/3Q4/4q3/8/8/7k b - - 0 1",
			[]DepthCount{{1, 6}, {2, 35}, {3, 495}, {4, 8349}, {5, 166741}, {6, 3370175}},
		},
		{"8/8/8/8/8/K7/P7/k7 w - - 0 1",
			[]DepthCount{{1, 3}, {2, 7}, {3, 43}, {4, 199}, {5, 1347}, {6, 6249}},
		},
		{"8/8/8/8/8/7K/7P/7k w - - 0 1",
			[]DepthCount{{1, 3}, {2, 7}, {3, 43}, {4, 199}, {5, 1347}, {6, 6249}},
		},
		{"K7/p7/k7/8/8/8/8/8 w - - 0 1",
			[]DepthCount{{1, 1}, {2, 3}, {3, 12}, {4, 80}, {5, 342}, {6, 2343}},
		},
		{"7K/7p/7k/8/8/8/8/8 w - - 0 1",
			[]DepthCount{{1, 1}, {2, 3}, {3, 12}, {4, 80}, {5, 342}, {6, 2343}},
		},
		{"8/2k1p3/3pP3/3P2K1/8/8/8/8 w - - 0 1",
			[]DepthCount{{1, 7}, {2, 35}, {3, 210}, {4, 1091}, {5, 7028}, {6, 34834}},
		},
		{"8/8/8/8/8/K7/P7/k7 b - - 0 1",
			[]DepthCount{{1, 1}, {2, 3}, {3, 12}, {4, 80}, {5, 342}, {6, 2343}},
		},
		{"8/8/8/8/8/7K/7P/7k b - - 0 1",
			[]DepthCount{{1, 1}, {2, 3}, {3, 12}, {4, 80}, {5, 342}, {6, 2343}},
		},
		{"K7/p7/k7/8/8/8/8/8 b - - 0 1",
			[]DepthCount{{1, 3}, {2, 7}, {3, 43}, {4, 199}, {5, 1347}, {6, 6249}},
		},
		{"7K/7p/7k/8/8/8/8/8 b - - 0 1",
			[]DepthCount{{1, 3}, {2, 7}, {3, 43}, {4, 199}, {5, 1347}, {6, 6249}},
		},
		{"8/2k1p3/3pP3/3P2K1/8/8/8/8 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 35}, {3, 182}, {4, 1091}, {5, 5408}, {6, 34822}},
		},
		{"8/8/8/8/8/4k3/4P3/4K3 w - - 0 1",
			[]DepthCount{{1, 2}, {2, 8}, {3, 44}, {4, 282}, {5, 1814}, {6, 11848}},
		},
		{"4k3/4p3/4K3/8/8/8/8/8 b - - 0 1",
			[]DepthCount{{1, 2}, {2, 8}, {3, 44}, {4, 282}, {5, 1814}, {6, 11848}},
		},
		{"8/8/7k/7p/7P/7K/8/8 w - - 0 1",
			[]DepthCount{{1, 3}, {2, 9}, {3, 57}, {4, 360}, {5, 1969}, {6, 10724}},
		},
		{"8/8/k7/p7/P7/K7/8/8 w - - 0 1",
			[]DepthCount{{1, 3}, {2, 9}, {3, 57}, {4, 360}, {5, 1969}, {6, 10724}},
		},
		{"8/8/3k4/3p4/3P4/3K4/8/8 w - - 0 1",
			[]DepthCount{{1, 5}, {2, 25}, {3, 180}, {4, 1294}, {5, 8296}, {6, 53138}},
		},
		{"8/3k4/3p4/8/3P4/3K4/8/8 w - - 0 1",
			[]DepthCount{{1, 8}, {2, 61}, {3, 483}, {4, 3213}, {5, 23599}, {6, 157093}},
		},
		{"8/8/3k4/3p4/8/3P4/3K4/8 w - - 0 1",
			[]DepthCount{{1, 8}, {2, 61}, {3, 411}, {4, 3213}, {5, 21637}, {6, 158065}},
		},
		{"k7/8/3p4/8/3P4/8/8/7K w - - 0 1",
			[]DepthCount{{1, 4}, {2, 15}, {3, 90}, {4, 534}, {5, 3450}, {6, 20960}},
		},
		{"8/8/7k/7p/7P/7K/8/8 b - - 0 1",
			[]DepthCount{{1, 3}, {2, 9}, {3, 57}, {4, 360}, {5, 1969}, {6, 10724}},
		},
		{"8/8/k7/p7/P7/K7/8/8 b - - 0 1",
			[]DepthCount{{1, 3}, {2, 9}, {3, 57}, {4, 360}, {5, 1969}, {6, 10724}},
		},
		{"8/8/3k4/3p4/3P4/3K4/8/8 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 25}, {3, 180}, {4, 1294}, {5, 8296}, {6, 53138}},
		},
		{"8/3k4/3p4/8/3P4/3K4/8/8 b - - 0 1",
			[]DepthCount{{1, 8}, {2, 61}, {3, 411}, {4, 3213}, {5, 21637}, {6, 158065}},
		},
		{"8/8/3k4/3p4/8/3P4/3K4/8 b - - 0 1",
			[]DepthCount{{1, 8}, {2, 61}, {3, 483}, {4, 3213}, {5, 23599}, {6, 157093}},
		},
		{"k7/8/3p4/8/3P4/8/8/7K b - - 0 1",
			[]DepthCount{{1, 4}, {2, 15}, {3, 89}, {4, 537}, {5, 3309}, {6, 21104}},
		},
		{"7k/3p4/8/8/3P4/8/8/K7 w - - 0 1",
			[]DepthCount{{1, 4}, {2, 19}, {3, 117}, {4, 720}, {5, 4661}, {6, 32191}},
		},
		{"7k/8/8/3p4/8/8/3P4/K7 w - - 0 1",
			[]DepthCount{{1, 5}, {2, 19}, {3, 116}, {4, 716}, {5, 4786}, {6, 30980}},
		},
		{"k7/8/8/7p/6P1/8/8/K7 w - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"k7/8/7p/8/8/6P1/8/K7 w - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"k7/8/8/6p1/7P/8/8/K7 w - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"k7/8/6p1/8/8/7P/8/K7 w - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"k7/8/8/3p4/4p3/8/8/7K w - - 0 1",
			[]DepthCount{{1, 3}, {2, 15}, {3, 84}, {4, 573}, {5, 3013}, {6, 22886}},
		},
		{"k7/8/3p4/8/8/4P3/8/7K w - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4271}, {6, 28662}},
		},
		{"7k/3p4/8/8/3P4/8/8/K7 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 19}, {3, 117}, {4, 720}, {5, 5014}, {6, 32167}},
		},
		{"7k/8/8/3p4/8/8/3P4/K7 b - - 0 1",
			[]DepthCount{{1, 4}, {2, 19}, {3, 117}, {4, 712}, {5, 4658}, {6, 30749}},
		},
		{"k7/8/8/7p/6P1/8/8/K7 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"k7/8/7p/8/8/6P1/8/K7 b - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"k7/8/8/6p1/7P/8/8/K7 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"k7/8/6p1/8/8/7P/8/K7 b - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"k7/8/8/3p4/4p3/8/8/7K b - - 0 1",
			[]DepthCount{{1, 5}, {2, 15}, {3, 102}, {4, 569}, {5, 4337}, {6, 22579}},
		},
		{"k7/8/3p4/8/8/4P3/8/7K b - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4271}, {6, 28662}},
		},
		{"7k/8/8/p7/1P6/8/8/7K w - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"7k/8/p7/8/8/1P6/8/7K w - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"7k/8/8/1p6/P7/8/8/7K w - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"7k/8/1p6/8/8/P7/8/7K w - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"k7/7p/8/8/8/8/6P1/K7 w - - 0 1",
			[]DepthCount{{1, 5}, {2, 25}, {3, 161}, {4, 1035}, {5, 7574}, {6, 55338}},
		},
		{"k7/6p1/8/8/8/8/7P/K7 w - - 0 1",
			[]DepthCount{{1, 5}, {2, 25}, {3, 161}, {4, 1035}, {5, 7574}, {6, 55338}},
		},
		{"3k4/3pp3/8/8/8/8/3PP3/3K4 w - - 0 1",
			[]DepthCount{{1, 7}, {2, 49}, {3, 378}, {4, 2902}, {5, 24122}, {6, 199002}},
		},
		{"7k/8/8/p7/1P6/8/8/7K b - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"7k/8/p7/8/8/1P6/8/7K b - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"7k/8/8/1p6/P7/8/8/7K b - - 0 1",
			[]DepthCount{{1, 5}, {2, 22}, {3, 139}, {4, 877}, {5, 6112}, {6, 41874}},
		},
		{"7k/8/1p6/8/8/P7/8/7K b - - 0 1",
			[]DepthCount{{1, 4}, {2, 16}, {3, 101}, {4, 637}, {5, 4354}, {6, 29679}},
		},
		{"k7/7p/8/8/8/8/6P1/K7 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 25}, {3, 161}, {4, 1035}, {5, 7574}, {6, 55338}},
		},
		{"k7/6p1/8/8/8/8/7P/K7 b - - 0 1",
			[]DepthCount{{1, 5}, {2, 25}, {3, 161}, {4, 1035}, {5, 7574}, {6, 55338}},
		},
		{"3k4/3pp3/8/8/8/8/3PP3/3K4 b - - 0 1",
			[]DepthCount{{1, 7}, {2, 49}, {3, 378}, {4, 2902}, {5, 24122}, {6, 199002}},
		},
		{"8/Pk6/8/8/8/8/6Kp/8 w - - 0 1",
			[]DepthCount{{1, 11}, {2, 97}, {3, 887}, {4, 8048}, {5, 90606}, {6, 1030499}},
		},
		{"n1n5/1Pk5/8/8/8/8/5Kp1/5N1N w - - 0 1",
			[]DepthCount{{1, 24}, {2, 421}, {3, 7421}, {4, 124608}, {5, 2193768}, {6, 37665329}},
		},
		{"8/PPPk4/8/8/8/8/4Kppp/8 w - - 0 1",
			[]DepthCount{{1, 18}, {2, 270}, {3, 4699}, {4, 79355}, {5, 1533145}, {6, 28859283}},
		},
		{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N w - - 0 1",
			[]DepthCount{{1, 24}, {2, 496}, {3, 9483}, {4, 182838}, {5, 3605103}, {6, 71179139}},
		},
		{"8/Pk6/8/8/8/8/6Kp/8 b - - 0 1",
			[]DepthCount{{1, 11}, {2, 97}, {3, 887}, {4, 8048}, {5, 90606}, {6, 1030499}},
		},
		{"n1n5/1Pk5/8/8/8/8/5Kp1/5N1N b - - 0 1",
			[]DepthCount{{1, 24}, {2, 421}, {3, 7421}, {4, 124608}, {5, 2193768}, {6, 37665329}},
		},
		{"8/PPPk4/8/8/8/8/4Kppp/8 b - - 0 1",
			[]DepthCount{{1, 18}, {2, 270}, {3, 4699}, {4, 79355}, {5, 1533145}, {6, 28859283}},
		},
		{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1",
			[]DepthCount{{1, 24}, {2, 496}, {3, 9483}, {4, 182838}, {5, 3605103}, {6, 71179139}},
		},
	}
	InitHashKeys(181818)
	for _, tt := range tests {
		for _, dc := range tt.depthCounts {
			if testing.Short() == true && dc.count >= 1000000 {
				continue
			}
			b, err := NewBoard(tt.fen)
			if err != nil {
				t.Error(err)
			}
			actual := Perft(b, dc.depth)
			if actual != dc.count {
				t.Errorf("nodes %v != %v", actual, dc.count)
			}
		}
	}
}
