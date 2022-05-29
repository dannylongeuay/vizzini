package main

import (
	"fmt"
	"math"
	"time"
)

type Search struct {
	*Board
	ply      int
	maxDepth int
	nodes    int
	maxNodes int
	bestMove Move
	stop     bool
	stopTime time.Time
	pvTable  []PvMove
	pvLine   []Move
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
	search.Reset()
	return &search, nil
}

func (s *Search) Clear() {
	s.bestMove = 0
	s.stop = false
	s.stopTime = time.Time{}
}

func (s *Search) Reset() {
	s.Clear()
	s.pvTable = make([]PvMove, PV_TABLE_SIZE)
}

func (s *Search) IterativeDeepening() int {
	startTime := time.Now()
	var bestScore int
	for i := 1; i <= s.maxDepth; i++ {
		bestScore = s.Negamax(i, math.MinInt+1, math.MaxInt)
		fmt.Printf("info depth %v score %v nodes %v time %v %v\n",
			i, bestScore, s.nodes, time.Since(startTime).Milliseconds(), s.GetPvLineString())
		if time.Now().After(s.stopTime) {
			break
		}
	}
	return bestScore
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
			s.SetPvMove(move)
			if s.ply == 0 {
				s.bestMove = move
			}
		}
	}

	return alpha
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
