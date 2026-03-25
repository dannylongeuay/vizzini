package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

const lichessBaseURL = "https://lichess.org"

// LichessBot holds the state for a Lichess bot session.
type LichessBot struct {
	token   string
	client  *http.Client
	botID   string
	playing atomic.Bool
}

// JSON types for Lichess API responses.

type lichessEvent struct {
	Type      string          `json:"type"`
	Challenge json.RawMessage `json:"challenge,omitempty"`
	Game      json.RawMessage `json:"game,omitempty"`
}

type lichessChallenge struct {
	ID      string `json:"id"`
	Rated   bool   `json:"rated"`
	Variant struct {
		Key string `json:"key"`
	} `json:"variant"`
	Speed string `json:"speed"`
}

type lichessGameStart struct {
	GameID string `json:"gameId"`
	FullID string `json:"fullId"`
	Color  string `json:"color"`
}

type lichessGameFull struct {
	ID         string           `json:"id"`
	White      lichessPlayer    `json:"white"`
	Black      lichessPlayer    `json:"black"`
	InitialFen string           `json:"initialFen"`
	State      lichessGameState `json:"state"`
	Type       string           `json:"type"`
}

type lichessGameState struct {
	Type   string `json:"type"`
	Moves  string `json:"moves"`
	Wtime  int64  `json:"wtime"`
	Btime  int64  `json:"btime"`
	Winc   int64  `json:"winc"`
	Binc   int64  `json:"binc"`
	Status string `json:"status"`
}

type lichessPlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type lichessOpponentGone struct {
	Type              string `json:"type"`
	Gone              bool   `json:"gone"`
	ClaimWinInSeconds int    `json:"claimWinInSeconds"`
}

type lichessAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// ModeLichess is the entry point for the lichess subcommand.
func ModeLichess() {
	token := os.Getenv("LICHESS_TOKEN")
	if token == "" {
		log.Fatal("LICHESS_TOKEN environment variable is required")
	}

	bot := &LichessBot{
		token:  token,
		client: &http.Client{},
	}

	if err := bot.fetchAccount(); err != nil {
		log.Fatalf("Failed to fetch account: %v", err)
	}
	log.Printf("Logged in as %s", bot.botID)

	bot.streamEvents()
}

func (b *LichessBot) fetchAccount() error {
	resp, err := b.apiGet("/api/account")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET /api/account returned %d: %s", resp.StatusCode, body)
	}

	var account lichessAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return err
	}
	b.botID = strings.ToLower(account.ID)
	return nil
}

// HTTP helpers

func (b *LichessBot) apiRequest(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, lichessBaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+b.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return b.client.Do(req)
}

func (b *LichessBot) apiGet(path string) (*http.Response, error) {
	return b.apiRequest("GET", path, nil)
}

func (b *LichessBot) apiPost(path string, body io.Reader) (*http.Response, error) {
	return b.apiRequest("POST", path, body)
}

func (b *LichessBot) apiPostEmpty(path string) {
	resp, err := b.apiPost(path, nil)
	if err != nil {
		log.Printf("POST %s error: %v", path, err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("POST %s returned %d", path, resp.StatusCode)
	}
}

// Event stream

func (b *LichessBot) streamEvents() {
	backoff := 5 * time.Second
	const maxBackoff = 60 * time.Second

	for {
		err := b.streamEventsOnce()
		if err != nil {
			log.Printf("Event stream error: %v, reconnecting in %v...", err, backoff)
			time.Sleep(backoff)
			backoff = min(backoff*2, maxBackoff)
		} else {
			// Stream ended cleanly (e.g. server closed), reset backoff
			backoff = 5 * time.Second
			time.Sleep(backoff)
		}
	}
}

func (b *LichessBot) streamEventsOnce() error {
	resp, err := b.apiGet("/api/stream/event")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET /api/stream/event returned %d: %s", resp.StatusCode, body)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var event lichessEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			log.Printf("Failed to parse event: %v", err)
			continue
		}

		switch event.Type {
		case "challenge":
			var c lichessChallenge
			if err := json.Unmarshal(event.Challenge, &c); err != nil {
				log.Printf("Failed to parse challenge: %v", err)
				continue
			}
			b.handleChallenge(c)
		case "gameStart":
			var g lichessGameStart
			if err := json.Unmarshal(event.Game, &g); err != nil {
				log.Printf("Failed to parse gameStart: %v", err)
				continue
			}
			b.handleGameStart(g)
		case "gameFinish":
			b.playing.Store(false)
			log.Println("Game finished")
		case "challengeCanceled":
			log.Println("Challenge canceled")
		case "challengeDeclined":
			log.Println("Challenge declined")
		}
	}
	return scanner.Err()
}

// Challenge handling

func (b *LichessBot) handleChallenge(c lichessChallenge) {
	if b.playing.Load() {
		log.Printf("Declining challenge %s: already playing", c.ID)
		b.declineChallenge(c.ID)
		return
	}
	if c.Rated {
		log.Printf("Declining challenge %s: rated", c.ID)
		b.declineChallenge(c.ID)
		return
	}
	if c.Variant.Key != "standard" {
		log.Printf("Declining challenge %s: variant %s", c.ID, c.Variant.Key)
		b.declineChallenge(c.ID)
		return
	}
	if c.Speed != "bullet" && c.Speed != "blitz" && c.Speed != "rapid" {
		log.Printf("Declining challenge %s: speed %s", c.ID, c.Speed)
		b.declineChallenge(c.ID)
		return
	}

	log.Printf("Accepting challenge %s", c.ID)
	b.acceptChallenge(c.ID)
}

func (b *LichessBot) acceptChallenge(id string) {
	b.apiPostEmpty("/api/challenge/" + id + "/accept")
}

func (b *LichessBot) declineChallenge(id string) {
	b.apiPostEmpty("/api/challenge/" + id + "/decline")
}

// Game start

func (b *LichessBot) handleGameStart(g lichessGameStart) {
	if !b.playing.CompareAndSwap(false, true) {
		log.Printf("Ignoring gameStart %s: already playing", g.GameID)
		return
	}
	log.Printf("Starting game %s", g.GameID)
	go b.playGame(g.GameID)
}

// Game loop

func (b *LichessBot) playGame(gameID string) {
	defer b.playing.Store(false)

	resp, err := b.apiGet("/api/bot/game/stream/" + gameID)
	if err != nil {
		log.Printf("Failed to stream game %s: %v", gameID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Game stream %s returned %d: %s", gameID, resp.StatusCode, body)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	var botColor Color
	var initialFen string
	firstEvent := true

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if firstEvent {
			firstEvent = false

			var full lichessGameFull
			if err := json.Unmarshal([]byte(line), &full); err != nil {
				log.Printf("Failed to parse gameFull: %v", err)
				return
			}

			if strings.ToLower(full.White.ID) == b.botID {
				botColor = WHITE
			} else {
				botColor = BLACK
			}

			initialFen = full.InitialFen
			if initialFen == "startpos" {
				initialFen = STARTING_FEN
			}

			log.Printf("Game %s: playing as %s", gameID, colorName(botColor))
			b.sendChat(gameID, "player", "glhf")

			b.handleGameState(gameID, botColor, initialFen, full.State)
			continue
		}

		// Peek at the type field
		var peek struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(line), &peek); err != nil {
			log.Printf("Failed to parse event type: %v", err)
			continue
		}

		switch peek.Type {
		case "gameState":
			var state lichessGameState
			if err := json.Unmarshal([]byte(line), &state); err != nil {
				log.Printf("Failed to parse gameState: %v", err)
				continue
			}
			if state.Status != "started" {
				log.Printf("Game %s ended: %s", gameID, state.Status)
				b.sendChat(gameID, "player", "gg")
				return
			}
			b.handleGameState(gameID, botColor, initialFen, state)
		case "opponentGone":
			var gone lichessOpponentGone
			if err := json.Unmarshal([]byte(line), &gone); err != nil {
				log.Printf("Failed to parse opponentGone: %v", err)
				continue
			}
			if gone.Gone {
				wait := time.Duration(gone.ClaimWinInSeconds) * time.Second
				log.Printf("Game %s: opponent gone, claiming victory in %v", gameID, wait)
				time.AfterFunc(wait, func() {
					b.claimVictory(gameID)
				})
			}
		case "chatLine":
			// ignore
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Game stream %s scanner error: %v", gameID, err)
	}
}

func colorName(color Color) string {
	if color == WHITE {
		return "white"
	}
	return "black"
}

// Move calculation

func (b *LichessBot) handleGameState(gameID string, botColor Color, initialFen string, state lichessGameState) {
	// Determine side to move from move count
	moves := strings.Fields(state.Moves)
	var sideToMove Color = WHITE
	if len(moves)%2 == 1 {
		sideToMove = BLACK
	}

	if sideToMove != botColor {
		return
	}

	// Build position args for NewBoardFromUCIPosition
	var args []string
	if initialFen == STARTING_FEN {
		args = []string{"startpos"}
	} else {
		fenParts := strings.Fields(initialFen)
		args = append([]string{"fen"}, fenParts...)
	}
	if len(moves) > 0 {
		args = append(args, "moves")
		args = append(args, moves...)
	}

	board, err := NewBoardFromUCIPosition(args)
	if err != nil {
		log.Printf("Failed to parse position: %v", err)
		return
	}

	search, err := NewSearch(board.ToFEN(), DEFAULT_MAX_DEPTH, DEFAULT_MAX_NODES)
	if err != nil {
		log.Printf("Failed to create search: %v", err)
		return
	}
	search.quiet = true

	// Time management: remainingTime/30 + increment - SEARCH_BUFFER
	var remainingTime, increment int64
	if botColor == WHITE {
		remainingTime = state.Wtime
		increment = state.Winc
	} else {
		remainingTime = state.Btime
		increment = state.Binc
	}

	mtime := remainingTime/30 + increment - SEARCH_BUFFER
	if mtime < 1 {
		mtime = 1
	}
	search.stopTime = time.Now().Add(time.Millisecond * time.Duration(mtime))

	search.IterativeDeepening()

	bestMove := search.bestMove
	if bestMove == 0 {
		// Fallback: pick first legal move
		legalMoves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		search.GenerateMoves(&legalMoves, search.sideToMove, false)
		for _, m := range legalMoves {
			if err := search.MakeMove(m); err == nil {
				search.UndoMove()
				bestMove = m
				break
			}
		}
	}

	if bestMove == 0 {
		log.Printf("Game %s: no legal moves available", gameID)
		return
	}

	uciMove := bestMove.ToUCIString()
	log.Printf("Game %s: playing %s", gameID, uciMove)

	// Submit move
	b.apiPostEmpty("/api/bot/game/" + gameID + "/move/" + uciMove)
}

// Chat

func (b *LichessBot) sendChat(gameID, room, text string) {
	data := url.Values{}
	data.Set("room", room)
	data.Set("text", text)
	resp, err := b.apiPost("/api/bot/game/"+gameID+"/chat", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to send chat: %v", err)
		return
	}
	resp.Body.Close()
}

// Game action helpers

func (b *LichessBot) claimVictory(gameID string) {
	b.apiPostEmpty("/api/bot/game/" + gameID + "/claim-victory")
}

func (b *LichessBot) claimDraw(gameID string) {
	b.apiPostEmpty("/api/bot/game/" + gameID + "/claim-draw")
}
