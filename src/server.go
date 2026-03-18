package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SERVER_PORT            = ":8080"
	DEFAULT_ALLOWED_ORIGIN = "https://chess.cyberdan.dev"
)

// Request types

type ValidMovesRequest struct {
	FEN string `json:"fen"`
}

type SubmitMoveRequest struct {
	FEN  string `json:"fen"`
	Move string `json:"move"`
}

type BestMoveRequest struct {
	FEN       string `json:"fen"`
	Depth     *int   `json:"depth"`
	TimeoutMs *int   `json:"timeout_ms"`
}

// Shared types

type MoveInfo struct {
	UCI       string  `json:"uci"`
	SAN       string  `json:"san"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Capture   bool    `json:"capture"`
	Castling  bool    `json:"castling"`
	Check     bool    `json:"check"`
	Promotion *string `json:"promotion"`
}

// Response types

type ValidMovesResponse struct {
	FEN        string     `json:"fen"`
	SideToMove string     `json:"side_to_move"`
	Status     string     `json:"status"`
	MoveCount  int        `json:"move_count"`
	Moves      []MoveInfo `json:"moves"`
}

type SubmitMoveResponse struct {
	UCI        string     `json:"uci"`
	SAN        string     `json:"san"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	FEN        string     `json:"fen"`
	Status     string     `json:"status"`
	SideToMove string     `json:"side_to_move"`
	MoveCount  int        `json:"move_count"`
	Moves      []MoveInfo `json:"moves"`
}

type BestMoveResponse struct {
	FEN    string  `json:"fen"`
	Depth  int     `json:"depth"`
	UCI    *string `json:"uci"`
	SAN    *string `json:"san"`
	From   *string `json:"from"`
	To     *string `json:"to"`
	Score  int     `json:"score"`
	Nodes  int     `json:"nodes"`
	Source string  `json:"source"`
}

type SubmitBestMoveResponse struct {
	UCI        string     `json:"uci"`
	SAN        string     `json:"san"`
	From       string     `json:"from"`
	To         string     `json:"to"`
	FEN        string     `json:"fen"`
	Status     string     `json:"status"`
	SideToMove string     `json:"side_to_move"`
	MoveCount  int        `json:"move_count"`
	Moves      []MoveInfo `json:"moves"`
	Depth      int        `json:"depth"`
	Score      int        `json:"score"`
	Nodes      int        `json:"nodes"`
	Source     string     `json:"source"`
}

func ModeServe() {
	allowedOrigin := DEFAULT_ALLOWED_ORIGIN
	if os.Getenv("CORS_PERMISSIVE") != "" {
		allowedOrigin = "*"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /validmoves", handleValidMoves)
	mux.HandleFunc("POST /submitmove", handleSubmitMove)
	mux.HandleFunc("POST /bestmove", handleBestMove)
	mux.HandleFunc("POST /submitbestmove", handleSubmitBestMove)
	mux.HandleFunc("/", handleNotFound)

	handler := maxBodyMiddleware(corsMiddleware(mux, allowedOrigin), 1<<20)

	fmt.Printf("Listening on %s (CORS origin: %s)\n", SERVER_PORT, allowedOrigin)
	if err := http.ListenAndServe(SERVER_PORT, handler); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func maxBodyMiddleware(next http.Handler, maxBytes int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code string, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": code, "message": message})
}

func buildMoveInfo(board *Board, move Move) MoveInfo {
	var mu MoveUnpacked
	move.Unpack(&mu)

	san := board.ToSAN(move)
	isCheck := strings.HasSuffix(san, "+") || strings.HasSuffix(san, "#")

	return MoveInfo{
		UCI:       move.ToUCIString(),
		SAN:       san,
		From:      stringCoordLower(mu.originCoord),
		To:        stringCoordLower(mu.dstCoord),
		Capture:   mu.IsCapture(),
		Castling:  mu.IsCastling(),
		Check:     isCheck,
		Promotion: mu.PromotionString(),
	}
}

func buildPositionInfo(board *Board) (string, string, int, []MoveInfo) {
	legal := board.LegalMoves()
	status := board.GameStatus(legal)
	fen := board.ToFEN()
	moves := make([]MoveInfo, len(legal))
	for i, m := range legal {
		moves[i] = buildMoveInfo(board, m)
	}
	return fen, status, len(legal), moves
}

func sideToMoveStr(board *Board) string {
	if board.sideToMove == WHITE {
		return "white"
	}
	return "black"
}

func runSearch(fen string, depth *int, timeoutMs *int) (*Search, int, error) {
	maxDepth := DEFAULT_MAX_DEPTH
	search, err := NewSearch(fen, maxDepth, 0)
	if err != nil {
		return nil, 0, err
	}
	search.quiet = true

	var bestScore int
	if timeoutMs != nil {
		dur := time.Duration(*timeoutMs)*time.Millisecond - time.Duration(SEARCH_BUFFER)*time.Millisecond
		if dur < 0 {
			dur = 0
		}
		search.stopTime = time.Now().Add(dur)
		if depth != nil && *depth < DEFAULT_MAX_DEPTH {
			search.maxDepth = *depth
		}
	} else if depth != nil {
		d := *depth
		if d > 20 {
			d = 20
		}
		search.maxDepth = d
	} else {
		search.stopTime = time.Now().Add(5 * time.Second)
	}

	bestScore = search.IterativeDeepening()
	return search, bestScore, nil
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, "not_found", fmt.Sprintf("no route for %s %s", r.Method, r.URL.Path), http.StatusNotFound)
}

func handleValidMoves(w http.ResponseWriter, r *http.Request) {
	var req ValidMovesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "could not decode request body", http.StatusBadRequest)
		return
	}
	if req.FEN == "" {
		writeError(w, "missing_fen", "fen field is required", http.StatusBadRequest)
		return
	}
	board, err := NewBoard(req.FEN)
	if err != nil {
		writeError(w, "invalid_fen", err.Error(), http.StatusBadRequest)
		return
	}

	fen, status, count, moves := buildPositionInfo(board)
	writeJSON(w, ValidMovesResponse{
		FEN:        fen,
		SideToMove: sideToMoveStr(board),
		Status:     status,
		MoveCount:  count,
		Moves:      moves,
	})
}

func handleSubmitMove(w http.ResponseWriter, r *http.Request) {
	var req SubmitMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "could not decode request body", http.StatusBadRequest)
		return
	}
	if req.FEN == "" {
		writeError(w, "missing_fen", "fen field is required", http.StatusBadRequest)
		return
	}
	if req.Move == "" {
		writeError(w, "missing_move", "move field is required", http.StatusBadRequest)
		return
	}

	board, err := NewBoard(req.FEN)
	if err != nil {
		writeError(w, "invalid_fen", err.Error(), http.StatusBadRequest)
		return
	}

	move, err := board.ParseMove(req.Move)
	if err != nil {
		writeError(w, "invalid_move", err.Error(), http.StatusBadRequest)
		return
	}

	uci := move.ToUCIString()
	var mu MoveUnpacked
	move.Unpack(&mu)
	from := stringCoordLower(mu.originCoord)
	to := stringCoordLower(mu.dstCoord)
	san := board.ToSAN(move)

	if moveErr := board.MakeMove(move); moveErr != nil {
		writeError(w, "invalid_move", moveErr.Error(), http.StatusBadRequest)
		return
	}

	fen, status, count, moves := buildPositionInfo(board)
	writeJSON(w, SubmitMoveResponse{
		UCI:        uci,
		SAN:        san,
		From:       from,
		To:         to,
		FEN:        fen,
		Status:     status,
		SideToMove: sideToMoveStr(board),
		MoveCount:  count,
		Moves:      moves,
	})
}

func handleBestMove(w http.ResponseWriter, r *http.Request) {
	var req BestMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "could not decode request body", http.StatusBadRequest)
		return
	}
	if req.FEN == "" {
		writeError(w, "missing_fen", "fen field is required", http.StatusBadRequest)
		return
	}

	search, bestScore, err := runSearch(req.FEN, req.Depth, req.TimeoutMs)
	if err != nil {
		writeError(w, "invalid_fen", err.Error(), http.StatusBadRequest)
		return
	}

	resp := BestMoveResponse{
		FEN:    req.FEN,
		Depth:  search.completedDepth,
		Score:  bestScore,
		Nodes:  search.nodes,
		Source: "search",
	}

	if search.bestMove != 0 {
		uci := search.bestMove.ToUCIString()
		san := search.Board.ToSAN(search.bestMove)
		var mu MoveUnpacked
		search.bestMove.Unpack(&mu)
		from := stringCoordLower(mu.originCoord)
		to := stringCoordLower(mu.dstCoord)
		resp.UCI = &uci
		resp.SAN = &san
		resp.From = &from
		resp.To = &to
	}

	writeJSON(w, resp)
}

func handleSubmitBestMove(w http.ResponseWriter, r *http.Request) {
	var req BestMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "could not decode request body", http.StatusBadRequest)
		return
	}
	if req.FEN == "" {
		writeError(w, "missing_fen", "fen field is required", http.StatusBadRequest)
		return
	}

	search, bestScore, err := runSearch(req.FEN, req.Depth, req.TimeoutMs)
	if err != nil {
		writeError(w, "invalid_fen", err.Error(), http.StatusBadRequest)
		return
	}

	if search.bestMove == 0 {
		writeError(w, "no_move", "engine found no legal move", http.StatusUnprocessableEntity)
		return
	}

	move := search.bestMove
	board := search.Board

	uci := move.ToUCIString()
	var mu MoveUnpacked
	move.Unpack(&mu)
	from := stringCoordLower(mu.originCoord)
	to := stringCoordLower(mu.dstCoord)
	san := board.ToSAN(move)

	board.MakeMove(move)

	fen, status, count, moves := buildPositionInfo(board)
	writeJSON(w, SubmitBestMoveResponse{
		UCI:        uci,
		SAN:        san,
		From:       from,
		To:         to,
		FEN:        fen,
		Status:     status,
		SideToMove: sideToMoveStr(board),
		MoveCount:  count,
		Moves:      moves,
		Depth:      search.completedDepth,
		Score:      bestScore,
		Nodes:      search.nodes,
		Source:     "search",
	})
}
