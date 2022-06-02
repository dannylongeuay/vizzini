package main

import (
	"fmt"
	"strings"
	"time"
)

type Search struct {
	*Board
	currentDepth int
	maxDepth     int
	nodes        int
	maxNodes     int
	bestMove     Move
	stop         bool
	stopTime     time.Time
	pvTable      []PvMove
	pvLine       []Move
	fhf          float64
	fh           float64
	quiet        bool
	killers      [][]Move
	alphaHistory [][]int
}

type PvMove struct {
	move Move
	hash Hash
}

func NewSearch(fen string, maxDepth int, maxNodes int) (*Search, error) {
	var search Search
	board, err := NewBoard(fen)
	if err != nil {
		return &search, err
	}
	search.Board = board
	search.maxDepth = maxDepth
	search.maxNodes = maxNodes
	search.Clear()
	return &search, nil
}

func (s *Search) Reset() {
	s.currentDepth = 0
	s.nodes = 0
	s.bestMove = 0
	s.fhf = 0
	s.fh = 0
	s.stop = false
	s.stopTime = time.Time{}
	s.pvTable = make([]PvMove, PV_TABLE_SIZE)
	s.killers = make([][]Move, KILLERS_SIZE)
	s.killers[0] = make([]Move, KILLERS_DEPTH)
	s.killers[1] = make([]Move, KILLERS_DEPTH)
	s.alphaHistory = make([][]int, BOARD_SQUARES)
	for i := 0; i < len(s.alphaHistory); i++ {
		s.alphaHistory[i] = make([]int, BOARD_SQUARES)
	}
}

func (s *Search) Clear() {
	s.Reset()
}

func (s *Search) TimeCheck() {
	if !s.stopTime.IsZero() && time.Now().After(s.stopTime) {
		s.stop = true
	}
}

func (s *Search) IterativeDeepening() int {
	startTime := time.Now()
	var bestScore int
	for i := 1; i <= s.maxDepth; i++ {
		bestScore = s.Negamax(i, MIN_SCORE, MAX_SCORE)
		if s.stop {
			break
		}
		s.SendInfo(i, bestScore, startTime)
	}
	return bestScore
}

func (s *Search) SendInfo(depth int, score int, startTime time.Time) {
	if s.quiet {
		return
	}
	var standard strings.Builder
	fmt.Fprintf(&standard, "info depth %v", depth)
	fmt.Fprintf(&standard, " score %v", score)
	fmt.Fprintf(&standard, " nodes %v", s.nodes)
	fmt.Fprintf(&standard, " time %v ", time.Since(startTime).Milliseconds())
	fmt.Fprint(&standard, s.GetPvLineString())
	fmt.Println(standard.String())

	var debug strings.Builder
	fmt.Fprintf(&debug, "info string fhf/fh %.2f%%", (s.fhf/s.fh)*100)
	movesToMate := MAX_SCORE - score
	if movesToMate < 10 {
		fmt.Fprintf(&debug, " mate in %v", movesToMate)
	}
	fmt.Println(debug.String())
}

func (s *Search) QSearch(alpha int, beta int) int {
	if s.nodes&2047 == 0 {
		s.TimeCheck()
	}

	s.nodes++

	if s.Repetition() || s.halfMove >= 100 {
		return DRAW
	}

	currentScore := s.Evaluate()

	if currentScore >= beta {
		return beta
	}

	if currentScore > alpha {
		alpha = currentScore
	}

	var moves []Move
	var legalMoves int
	// TODO: We want to generate capture moves only
	s.GenerateMoves(&moves, s.sideToMove)
	FilterMovesByKind(CAPTURE, &moves)

	for i := 0; i < len(moves); i++ {
		move, mu := s.PickNextMove(i, &moves, 0)
		err := s.MakeMoveUnpacked(move, mu)
		s.currentDepth++
		if err != nil {
			s.currentDepth--
			s.UndoMove()
			continue
		}
		legalMoves++
		score := -s.QSearch(-beta, -alpha)
		s.currentDepth--
		s.UndoMove()

		if s.stop {
			return DRAW
		}

		if score >= beta {
			if legalMoves == 1 {
				s.fhf++
			}
			s.fh++
			return beta
		}

		if score > alpha {
			alpha = score
			s.SetPvMove(move)
			if s.currentDepth == 0 {
				s.bestMove = move
			}
		}
	}

	return alpha
}

func (s *Search) Negamax(depth int, alpha int, beta int) int {
	if depth == 0 {
		return s.QSearch(alpha, beta)
	}

	if s.nodes&2047 == 0 {
		s.TimeCheck()
	}

	if s.Repetition() || s.halfMove >= 100 {
		return DRAW
	}

	s.nodes++

	inCheck := s.CoordAttacked(s.kingCoords[s.sideToMove], s.sideToMove)
	if inCheck {
		depth++
	}

	var moves []Move
	var legalMoves int
	pvMove, _ := s.GetPvMove()
	s.GenerateMoves(&moves, s.sideToMove)

	for i := 0; i < len(moves); i++ {
		move, mu := s.PickNextMove(i, &moves, pvMove)
		err := s.MakeMoveUnpacked(move, mu)
		s.currentDepth++
		if err != nil {
			s.currentDepth--
			s.UndoMove()
			continue
		}
		legalMoves++
		score := -s.Negamax(depth-1, -beta, -alpha)
		s.currentDepth--
		s.UndoMove()

		if s.stop {
			return DRAW
		}

		if score >= beta {
			if legalMoves == 1 {
				s.fhf++
			}
			if mu.moveKind == QUIET {
				s.killers[1][s.currentDepth] = s.killers[0][s.currentDepth]
				s.killers[0][s.currentDepth] = move
			}
			s.fh++
			return beta
		}

		if score > alpha {
			alpha = score
			s.SetPvMove(move)
			if s.currentDepth == 0 {
				s.bestMove = move
			}
			if mu.moveKind == QUIET {
				s.alphaHistory[mu.originCoord][mu.dstCoord] += depth
			}
		}
	}

	if legalMoves == 0 {
		if inCheck {
			return MIN_SCORE + s.currentDepth
		} else {
			return DRAW
		}
	}

	return alpha
}

func (s *Search) Repetition() bool {
	if int(s.halfMove) > s.ply {
		return false
	}
	for i := s.ply - int(s.halfMove); i < s.ply-1; i++ {
		if s.hash == s.hashes[i] {
			return true
		}
	}
	return false
}

func (s *Search) SetPvMove(move Move) {
	index := s.hash % PV_TABLE_SIZE
	s.pvTable[index] = PvMove{move, s.hash}
}

func (s *Search) GetPvMove() (Move, error) {
	index := s.hash % PV_TABLE_SIZE
	if s.pvTable[index].hash == s.hash {
		return s.pvTable[index].move, nil
	}
	return 0, fmt.Errorf("no pv move found at:\n%v", s.ToString())
}

func (s *Search) SetPvLine() {
	s.pvLine = make([]Move, 0, 20)
	move, err := s.GetPvMove()
	var count int
	for err == nil && count < s.maxDepth {
		if s.MoveExists(move) {
			s.MakeMove(move)
			s.pvLine = append(s.pvLine, move)
			count++
		} else {
			break
		}
		move, err = s.GetPvMove()
	}
	for i := 0; i < count; i++ {
		s.UndoMove()
	}
}

func (s *Search) GetPvLineString() string {
	s.SetPvLine()
	line := "pv"
	for _, move := range s.pvLine {
		line += " " + move.ToUCIString()
	}
	return line
}

func (s *Search) PickNextMove(index int, movesPtr *[]Move, pvMove Move) (Move, MoveUnpacked) {
	moves := *movesPtr
	var bestOrder MoveOrder
	bestNum := index
	for i := index; i < len(moves); i++ {
		var mu MoveUnpacked
		moves[i].Unpack(&mu)
		order := mu.moveOrder
		if pvMove == moves[i] {
			order = 255
		} else if s.killers[0][s.currentDepth] == moves[i] {
			order = 9
		} else if s.killers[1][s.currentDepth] == moves[i] {
			order = 8
		} else if order == 0 {
			order = MoveOrder(s.alphaHistory[mu.originCoord][mu.dstCoord])
		}
		if order > bestOrder {
			bestOrder = order
			bestNum = i
		}
	}
	if index != bestNum {
		moves[index], moves[bestNum] = moves[bestNum], moves[index]
	}
	var mu MoveUnpacked
	moves[index].Unpack(&mu)
	return moves[index], mu
}
