package main

import (
	"fmt"
	"sort"
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
}

func (s *Search) Clear() {
	s.Reset()
	s.pvTable = make([]PvMove, PV_TABLE_SIZE)
}

func (s *Search) IterativeDeepening() int {
	startTime := time.Now()
	var bestScore int
	for i := 1; i <= s.maxDepth; i++ {
		bestScore = s.Negamax(i, MIN_SCORE, MAX_SCORE)
		s.SendInfo(i, bestScore, startTime)
		if !s.stopTime.IsZero() && time.Now().After(s.stopTime) {
			break
		}
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

func (s *Search) Negamax(depth int, alpha int, beta int) int {
	if depth == 0 {
		s.nodes++
		return s.Evaluate()
	}

	if s.Repetition() || s.halfMove >= 100 {
		return DRAW
	}

	s.nodes++

	var moves []Move
	var legalMoves int
	s.GenerateMoves(&moves, s.sideToMove)
	// TODO: Sorting entire list is inefficient
	sort.Slice(moves, func(i, j int) bool {
		moveOrderA := (moves[i] & MOVE_ORDER_MASK) >> MOVE_ORDER_SHIFT
		moveOrderB := (moves[j] & MOVE_ORDER_MASK) >> MOVE_ORDER_SHIFT
		return moveOrderA > moveOrderB
	})

	for _, move := range moves {
		err := s.MakeMove(move)
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

	if legalMoves == 0 {
		if s.CoordAttacked(s.kingCoords[s.sideToMove], s.sideToMove) {
			return MIN_SCORE + s.ply
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
