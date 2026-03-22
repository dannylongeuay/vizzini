package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type Search struct {
	*Board
	currentDepth     int
	maxDepth         int
	completedDepth   int
	nodes            int
	maxNodes         int
	bestMove         Move
	stop             bool
	stopTime         time.Time
	tt               []TTEntry
	pvLine           []Move
	fhf              float64
	fh               float64
	quiet            bool
	killers          [][]Move
	alphaHistory     [][]int
	nullMoveAllowed  bool
	evalNoise        int
	rng              *rand.Rand
	temperature      float64
	tempMoveLimit    int
	rootMoves        []ScoredMove
}

type ScoredMove struct {
	move  Move
	score int
}

const (
	TT_NONE  uint8 = iota
	TT_EXACT
	TT_ALPHA
	TT_BETA
)

// TTEntry stores search results for a position.
type TTEntry struct {
	hash  Hash
	move  Move
	score int32
	depth int8
	flag  uint8
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
	search.tt = make([]TTEntry, TT_SIZE)
	search.killers = make([][]Move, KILLERS_SIZE)
	search.killers[0] = make([]Move, KILLERS_DEPTH)
	search.killers[1] = make([]Move, KILLERS_DEPTH)
	search.alphaHistory = make([][]int, BOARD_SQUARES)
	for i := range search.alphaHistory {
		search.alphaHistory[i] = make([]int, BOARD_SQUARES)
	}
	search.evalNoise = DEFAULT_EVAL_NOISE
	search.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	search.temperature = DEFAULT_TEMPERATURE
	search.tempMoveLimit = DEFAULT_TEMP_MOVE_LIMIT
	search.Clear()
	return &search, nil
}

func (s *Search) ResetKeepTT() {
	s.currentDepth = 0
	s.completedDepth = 0
	s.nodes = 0
	s.bestMove = 0
	s.fhf = 0
	s.fh = 0
	s.stop = false
	s.stopTime = time.Time{}
	s.nullMoveAllowed = true
	s.rootMoves = s.rootMoves[:0]
	clear(s.killers[0])
	clear(s.killers[1])
	for i := range s.alphaHistory {
		clear(s.alphaHistory[i])
	}
}

func (s *Search) Reset() {
	s.ResetKeepTT()
	clear(s.tt)
}

func (s *Search) NoisyEvaluate() int {
	score := s.Evaluate()
	if s.evalNoise > 0 && s.rng != nil {
		noise := s.rng.Intn(2*s.evalNoise+1) - s.evalNoise
		score += noise
	}
	return score
}

func (s *Search) TemperatureSelect(bestScore int) {
	if s.temperature <= 0 || s.rng == nil || len(s.rootMoves) <= 1 {
		return
	}
	if s.fullMove > s.tempMoveLimit {
		return
	}
	if bestScore > MATE_THRESHOLD || bestScore < -MATE_THRESHOLD {
		return
	}

	const scoreWindow = 75
	var candidates []ScoredMove
	for _, sm := range s.rootMoves {
		if bestScore-sm.score <= scoreWindow {
			candidates = append(candidates, sm)
		}
	}
	if len(candidates) <= 1 {
		return
	}

	maxScore := candidates[0].score
	for _, c := range candidates {
		if c.score > maxScore {
			maxScore = c.score
		}
	}
	weights := make([]float64, len(candidates))
	var total float64
	for i, c := range candidates {
		weights[i] = math.Exp(float64(c.score-maxScore) / (s.temperature * 100.0))
		total += weights[i]
	}

	r := s.rng.Float64() * total
	var cum float64
	for i, w := range weights {
		cum += w
		if r <= cum {
			s.bestMove = candidates[i].move
			return
		}
	}
	s.bestMove = candidates[len(candidates)-1].move
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
	var lastBestMove Move
	const ASPIRATION_WINDOW = 50

	for i := 1; i <= s.maxDepth; i++ {
		s.rootMoves = s.rootMoves[:0]
		alpha := MIN_SCORE
		beta := MAX_SCORE

		// Use aspiration windows from depth 4 onward, centered on the previous score.
		if i >= 4 {
			alpha = bestScore - ASPIRATION_WINDOW
			beta = bestScore + ASPIRATION_WINDOW
		}

		savedNodes := s.nodes
		bestScore = s.Negamax(i, alpha, beta)

		// Re-search with full window on fail-high or fail-low.
		// Reset nodes so the info output only reflects the final search.
		if !s.stop && (bestScore <= alpha || bestScore >= beta) {
			s.nodes = savedNodes
			s.rootMoves = s.rootMoves[:0]
			bestScore = s.Negamax(i, MIN_SCORE, MAX_SCORE)
		}

		// If stopped during search (e.g. timeout), the result may be incomplete.
		// Fall back to the last fully-completed iteration's best move.
		if s.stop {
			if lastBestMove != 0 {
				s.bestMove = lastBestMove
			}
			break
		}
		lastBestMove = s.bestMove
		s.completedDepth = i
		s.SendInfo(i, bestScore, startTime)
	}
	if !s.stop {
		s.TemperatureSelect(bestScore)
	}
	return bestScore
}

func (s *Search) SendInfo(depth int, score int, startTime time.Time) {
	if s.quiet {
		return
	}
	var standard strings.Builder
	elapsedTimeMs := time.Since(startTime).Milliseconds() + 1
	nps := int64(s.nodes*1000) / elapsedTimeMs
	fmt.Fprintf(&standard, "info depth %v", depth)
	if score > MATE_THRESHOLD {
		movesToMate := (MAX_SCORE - score + 1) / 2
		fmt.Fprintf(&standard, " score mate %v", movesToMate)
	} else if score < -MATE_THRESHOLD {
		movesToMate := (MIN_SCORE - score) / 2
		fmt.Fprintf(&standard, " score mate %v", movesToMate)
	} else {
		fmt.Fprintf(&standard, " score cp %v", score)
	}
	fmt.Fprintf(&standard, " nodes %v", s.nodes)
	fmt.Fprintf(&standard, " nps %v", nps)
	fmt.Fprintf(&standard, " time %v ", elapsedTimeMs)
	fmt.Fprint(&standard, s.GetPvLineString())
	fmt.Println(standard.String())

	var debug strings.Builder
	fmt.Fprintf(&debug, "info string fhf/fh %.2f%%", (s.fhf/s.fh)*100)
	fmt.Println(debug.String())
}

func (s *Search) QSearch(alpha int, beta int) int {
	if s.nodes&2047 == 0 {
		s.TimeCheck()
	}

	s.nodes++

	if s.ply >= int(MAX_GAME_MOVES)-1 {
		return s.NoisyEvaluate()
	}

	if s.Repetition() || s.halfMove >= 100 {
		return DRAW
	}

	inCheck := s.CoordAttacked(s.kingCoords[s.sideToMove], s.sideToMove)

	var standPat int
	if !inCheck {
		standPat = s.NoisyEvaluate()

		if standPat >= beta {
			return beta
		}

		if standPat > alpha {
			alpha = standPat
		}
	}

	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	var legalMoves int
	if inCheck {
		s.GenerateMoves(&moves, s.sideToMove, false)
	} else {
		s.GenerateMoves(&moves, s.sideToMove, true)
	}

	for i := 0; i < len(moves); i++ {
		move := s.PickNextMove(i, &moves, 0)

		// Delta pruning: skip captures that cannot raise alpha.
		if !inCheck {
			var mu MoveUnpacked
			move.Unpack(&mu)
			// Never delta-prune promotions — they gain significant material.
			if mu.moveKind < KNIGHT_PROMOTION {
				capturedValue := SQUARE_SCORES[mu.dstSquare]
				if capturedValue < 0 {
					capturedValue = -capturedValue
				}
				if standPat+capturedValue+DELTA_MARGIN <= alpha {
					continue
				}
			}
		}

		err := s.MakeMove(move)
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
			return alpha
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
			if s.currentDepth == 0 {
				s.bestMove = move
			}
		}
	}

	if inCheck && legalMoves == 0 {
		// Mate distance from search root. Note: currentDepth is also incremented
		// during null-move search, but null-move cutoffs return beta (not the score),
		// so the off-by-one doesn't propagate to the root.
		return MIN_SCORE + s.currentDepth
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

	if s.ply >= int(MAX_GAME_MOVES)-1 {
		return s.NoisyEvaluate()
	}

	if s.Repetition() || s.halfMove >= 100 {
		return DRAW
	}

	s.nodes++

	// Probe transposition table.
	ttMove := Move(0)
	isPvNode := beta-alpha > 1
	if entry, ok := s.ProbeTT(); ok {
		ttMove = entry.move
		if int(entry.depth) >= depth && !isPvNode {
			ttScore := s.TTScoreToSearch(entry.score)
			switch entry.flag {
			case TT_EXACT:
				return ttScore
			case TT_ALPHA:
				if ttScore <= alpha {
					return alpha
				}
			case TT_BETA:
				if ttScore >= beta {
					return beta
				}
			}
		}
	}

	inCheck := s.CoordAttacked(s.kingCoords[s.sideToMove], s.sideToMove)
	if inCheck {
		depth++
	}

	// Null move pruning.
	// nullMoveAllowed prevents consecutive null moves (works because search is single-threaded).
	if !inCheck && s.nullMoveAllowed && depth >= 3 {
		hasPieces := false
		if s.sideToMove == WHITE {
			hasPieces = s.bbWN|s.bbWB|s.bbWR|s.bbWQ > 0
		} else {
			hasPieces = s.bbBN|s.bbBB|s.bbBR|s.bbBQ > 0
		}
		if hasPieces {
			oldEp := s.epCoord
			oldHalfMove := s.halfMove
			s.MakeNullMove()
			s.currentDepth++
			s.nullMoveAllowed = false
			R := 2
			if depth >= 6 {
				R = 3
			}
			score := -s.Negamax(depth-1-R, -beta, -beta+1)
			s.nullMoveAllowed = true
			s.currentDepth--
			s.UndoNullMove(oldEp, oldHalfMove)
			if s.stop {
				return alpha
			}
			if score >= beta {
				return beta
			}
		}
	}

	// Compute static eval for futility pruning and reverse futility pruning.
	var staticEval int
	futilityPruning := false
	if !inCheck && !isPvNode && depth <= 3 {
		staticEval = s.NoisyEvaluate()
		// Reverse futility pruning: if static eval is well above beta, prune.
		if staticEval-FUTILITY_MARGINS[depth] >= beta {
			return beta
		}
		futilityPruning = staticEval+FUTILITY_MARGINS[depth] <= alpha
	}

	moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
	var legalMoves int
	s.GenerateMoves(&moves, s.sideToMove, false)

	origAlpha := alpha
	var bestMove Move

	for i := 0; i < len(moves); i++ {
		move := s.PickNextMove(i, &moves, ttMove)
		err := s.MakeMove(move)
		s.currentDepth++
		if err != nil {
			s.currentDepth--
			s.UndoMove()
			continue
		}
		legalMoves++
		givesCheck := s.CoordAttacked(s.kingCoords[s.sideToMove], s.sideToMove)

		var score int
		moveKind := MoveKind(move & MOVE_KIND_MASK)

		// Futility pruning: skip quiet moves that cannot improve alpha.
		if futilityPruning && !givesCheck && moveKind == QUIET && legalMoves > 1 {
			s.currentDepth--
			s.UndoMove()
			continue
		}

		// Late move reductions (skip for moves that give check)
		if legalMoves > 4 && depth >= 3 && !inCheck && !givesCheck && moveKind == QUIET {
			score = -s.Negamax(depth-2, -beta, -alpha)
			if score > alpha {
				score = -s.Negamax(depth-1, -beta, -alpha)
			}
		} else {
			score = -s.Negamax(depth-1, -beta, -alpha)
		}

		s.currentDepth--
		s.UndoMove()

		if s.stop {
			return alpha
		}

		if s.currentDepth == 0 {
			s.rootMoves = append(s.rootMoves, ScoredMove{move, score})
		}

		if score >= beta {
			bestMove = move
			if legalMoves == 1 {
				s.fhf++
			}
			if moveKind == QUIET && s.currentDepth < KILLERS_DEPTH {
				s.killers[1][s.currentDepth] = s.killers[0][s.currentDepth]
				s.killers[0][s.currentDepth] = move
			}
			s.fh++
			s.StoreTT(move, beta, depth, TT_BETA)
			return beta
		}

		if score > alpha {
			alpha = score
			bestMove = move
			if s.currentDepth == 0 {
				s.bestMove = move
			}
			if moveKind == QUIET {
				originCoord := Coord((move & MOVE_ORIGIN_COORD_MASK) >> MOVE_ORIGIN_COORD_SHIFT)
				dstCoord := Coord((move & MOVE_DST_COORD_MASK) >> MOVE_DST_COORD_SHIFT)
				s.alphaHistory[originCoord][dstCoord] += depth
				if s.alphaHistory[originCoord][dstCoord] > MAX_HISTORY_VALUE {
					s.alphaHistory[originCoord][dstCoord] = MAX_HISTORY_VALUE
				}
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

	// Store result in transposition table.
	if alpha > origAlpha {
		s.StoreTT(bestMove, alpha, depth, TT_EXACT)
	} else {
		s.StoreTT(bestMove, alpha, depth, TT_ALPHA)
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

// StoreTT writes a search result into the transposition table.
// Mate scores are adjusted from distance-from-root to distance-from-position.
func (s *Search) StoreTT(move Move, score int, depth int, flag uint8) {
	index := s.hash & TT_MASK
	storeScore := score
	if storeScore > MATE_THRESHOLD {
		storeScore += s.currentDepth
	} else if storeScore < -MATE_THRESHOLD {
		storeScore -= s.currentDepth
	}
	s.tt[index] = TTEntry{
		hash:  s.hash,
		move:  move,
		score: int32(storeScore),
		depth: int8(depth),
		flag:  flag,
	}
}

// ProbeTT retrieves a TT entry for the current position.
// Returns the entry and whether the hash matched.
func (s *Search) ProbeTT() (TTEntry, bool) {
	index := s.hash & TT_MASK
	entry := s.tt[index]
	if entry.hash == s.hash && entry.flag != TT_NONE {
		return entry, true
	}
	return TTEntry{}, false
}

// TTScoreToSearch converts a stored mate score back to distance-from-root.
func (s *Search) TTScoreToSearch(score int32) int {
	sc := int(score)
	if sc > MATE_THRESHOLD {
		sc -= s.currentDepth
	} else if sc < -MATE_THRESHOLD {
		sc += s.currentDepth
	}
	return sc
}

func (s *Search) SetPvLine() {
	s.pvLine = s.pvLine[:0]
	var count int
	for count < s.maxDepth {
		entry, ok := s.ProbeTT()
		if !ok {
			break
		}
		if !s.MoveExists(entry.move) {
			break
		}
		if err := s.MakeMove(entry.move); err != nil {
			s.UndoMove()
			break
		}
		s.pvLine = append(s.pvLine, entry.move)
		count++
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

func (s *Search) PickNextMove(index int, movesPtr *[]Move, pvMove Move) Move {
	moves := *movesPtr
	var bestOrder MoveOrder
	bestNum := index
	for i := index; i < len(moves); i++ {
		var mu MoveUnpacked
		moves[i].Unpack(&mu)
		order := MVV_LVA_SCORES[mu.dstSquare][mu.originSquare]

		// Assign special ordering for promotions and EP captures
		switch mu.moveKind {
		case EP_CAPTURE:
			order = MVV_LVA_EN_PASSANT
		case KNIGHT_PROMOTION:
			order = MVV_LVA_KNIGHT_PROMOTION
		case BISHOP_PROMOTION:
			order = MVV_LVA_BISHOP_PROMOTION
		case ROOK_PROMOTION:
			order = MVV_LVA_ROOK_PROMOTION
		case QUEEN_PROMOTION:
			order = MVV_LVA_QUEEN_PROMOTION
		case KNIGHT_PROMOTION_CAPTURE:
			order = MVV_LVA_KNIGHT_PROMOTION_CAPTURE
		case BISHOP_PROMOTION_CAPTURE:
			order = MVV_LVA_BISHOP_PROMOTION_CAPTURE
		case ROOK_PROMOTION_CAPTURE:
			order = MVV_LVA_ROOK_PROMOTION_CAPTURE
		case QUEEN_PROMOTION_CAPTURE:
			order = MVV_LVA_QUEEN_PROMOTION_CAPTURE
		}

		if pvMove == moves[i] {
			order = 255
		} else if order > 0 {
			// Already has capture/promotion ordering
		} else if s.currentDepth < KILLERS_DEPTH && s.killers[0][s.currentDepth] == moves[i] {
			order = 9
		} else if s.currentDepth < KILLERS_DEPTH && s.killers[1][s.currentDepth] == moves[i] {
			order = 8
		} else if order == 0 {
			histVal := s.alphaHistory[mu.originCoord][mu.dstCoord]
			if histVal > 0 {
				// Scale history into 1-7 range proportionally
				order = MoveOrder(1 + histVal*6/MAX_HISTORY_VALUE)
				if order > 7 {
					order = 7
				}
			}
		}
		if order > bestOrder {
			bestOrder = order
			bestNum = i
		}
	}
	if index != bestNum {
		moves[index], moves[bestNum] = moves[bestNum], moves[index]
	}
	return moves[index]
}
