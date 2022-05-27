package main

type Search struct {
	*Board
	ply      int
	nodes    int
	bestMove Move
}

func (s *Search) Negamax(depth int, alpha int, beta int) int {

	if depth == 0 {
		return s.Evaluate()
	}

	s.nodes++

	var moves []Move
	s.GenerateMoves(&moves, s.sideToMove)

	for _, move := range moves {
		err := s.MakeMove(move)
		s.ply++
		if err != nil {
			s.ply--
			s.UndoMove()
			continue
		}
		score := -s.Negamax(depth-1, -beta, -alpha)
		s.ply--
		s.UndoMove()

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
			if s.ply == 0 {
				s.bestMove = move
				// fmt.Println("best move", depth, alpha, move.ToString())
			}
		}
	}

	return alpha
}
