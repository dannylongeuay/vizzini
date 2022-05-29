package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UCI struct {
	*Search
	debug      bool
	input      string
	ponderMode bool
}

func NewUCI() *UCI {
	var board Board
	var search Search
	search.Board = &board
	search.maxDepth = UCI_DEFAULT_DEPTH
	search.pvTable = make([]PvMove, PV_TABLE_SIZE)
	var uci UCI
	uci.Search = &search
	return &uci
}

func ModeUCI(scanner *bufio.Scanner) {
	uci := NewUCI()
	uci.SendOk()

	for scanner.Scan() {
		uci.input = scanner.Text()
		words := strings.Split(uci.input, " ")

		if len(words) == 0 {
			continue
		}
		command := words[0]

		var args []string
		if len(words) > 0 {
			args = words[1:]
		}

		switch command {
		case "uci":
			uci.SendOk()
		case "debug":
			uci.SetDebug(args)
		case "isready":
			uci.SendReady()
		case "setoption":
			uci.SetOption(args)
		case "register":
			uci.SendRegistration()
		case "ucinewgame":
			uci.SetNewGame()
		case "position":
			uci.SetPosition(args)
		case "go":
			uci.SendCalculations(args)
		case "stop":
			uci.SetStop()
		case "ponderhit":
			uci.SetPonder()
		case "quit":
			goto QUIT

		}
	}
QUIT:
}

func (u *UCI) SendOk() {
	fmt.Println("id name Vizzini")
	fmt.Println("id author Daniel")
	// fmt.Println("option")
	fmt.Println("uciok")
}

func (u *UCI) SetDebug(args []string) {
	if args[0] == "on" {
		u.debug = true

	} else if args[0] == "off" {
		u.debug = false
	}

}

func (u *UCI) SendReady() {
	fmt.Println("readyok")
}

func (u *UCI) SetOption(args []string) {
}

func (u *UCI) SendRegistration() {
	fmt.Println("register later")
}

func (u *UCI) SetNewGame() {
	args := []string{"startpos"}
	u.pvTable = make([]PvMove, PV_TABLE_SIZE)
	u.SetPosition(args)
}

func (u *UCI) SetPosition(args []string) {
	var err error
	u.Board, err = NewBoardFromUCIPosition(args)
	if err != nil {
		panic(err)
	}
}

func (u *UCI) SendCalculations(args []string) {
	u.SetGoParams(args)
	u.IterativeDeepening()
	fmt.Println("bestmove", u.bestMove.ToUCIString())
	u.bestMove = 0
}

func (u *UCI) SetStop() {
	u.stop = true
}

func (u *UCI) SetPonder() {
	panic("ponderhit not implemented")
}

func NewBoardFromUCIPosition(args []string) (*Board, error) {
	fen := STARTING_FEN
	moveTokenIndex := 1
	if args[0] == "fen" {
		moveTokenIndex = 7
		fen = strings.Join(args[1:moveTokenIndex], " ")
	}
	board, err := NewBoard(fen)
	if err != nil {
		return board, err
	}
	if len(args) > moveTokenIndex && args[moveTokenIndex] == "moves" {
		for _, uciMove := range args[moveTokenIndex+1:] {
			move, err := board.ParseUCIMove(uciMove)
			if err != nil {
				return board, err
			}
			err = board.MakeMove(move)
			if err != nil {
				return board, err
			}
		}
	}

	return board, nil
}

func (u *UCI) SetGoParams(args []string) {
	u.maxDepth = UCI_DEFAULT_DEPTH
	u.maxNodes = 0
	u.stop = false

	var maxMoveTime int64
	var remainingTime int64
	var increment int64
	var remainingMoves int64
	var searchInfinite bool

	for i, arg := range args {
		switch arg {
		case "searchmoves":
			panic("searchmoves not implemented")
		case "ponder":
			panic("ponder not implemented")
		case "wtime":
			if u.sideToMove != WHITE {
				continue
			}
			wtime, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			remainingTime = int64(wtime)
		case "btime":
			if u.sideToMove != BLACK {
				continue
			}
			btime, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			remainingTime = int64(btime)
		case "winc":
			if u.sideToMove != WHITE {
				continue
			}
			winc, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			increment = int64(winc)
		case "binc":
			if u.sideToMove != BLACK {
				continue
			}
			binc, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			increment = int64(binc)
		case "movestogo":
			movestogo, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			remainingMoves = int64(movestogo)
		case "depth":
			maxDepth, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			u.maxDepth = maxDepth
		case "nodes":
			maxNodes, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			u.maxNodes = maxNodes
		case "mate":
			panic("mate not implemented")
		case "movetime":
			moveTime, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
			maxMoveTime = int64(moveTime)
		case "infinite":
			searchInfinite = true
		}
	}

	if searchInfinite {
		u.stopTime = time.Now().Add(time.Minute)
	} else if remainingTime > 0 {
		mtime := remainingTime/remainingMoves + increment - SEARCH_BUFFER
		if maxMoveTime > 0 {
			mtime = maxMoveTime - SEARCH_BUFFER
		}
		u.stopTime = time.Now().Add(time.Millisecond * time.Duration(mtime))
	} else {
		u.stopTime = time.Now().Add(time.Second * 3)
	}
}

func (b *Board) ParseUCIMove(s string) (Move, error) {
	sParts := []rune(s)

	if len(sParts) < 4 || len(sParts) > 5 {
		return 0, fmt.Errorf("Invalid UCI move notation: %v", s)
	}

	originCoord, err := StringToCoord(string(sParts[:2]))
	if err != nil {
		return 0, err
	}

	dstCoord, err := StringToCoord(string(sParts[2:4]))
	if err != nil {
		return 0, err
	}

	originSquare := b.squares[originCoord]
	dstSquare := b.squares[dstCoord]

	moveKind := QUIET
	if dstSquare != EMPTY {
		moveKind = CAPTURE
	}

	if originSquare == WHITE_PAWN && dstCoord-originCoord == 16 {
		moveKind = DOUBLE_PAWN_PUSH
	}

	if originSquare == BLACK_PAWN && originCoord-dstCoord == 16 {
		moveKind = DOUBLE_PAWN_PUSH
	}

	if dstCoord == b.epCoord && (originSquare == WHITE_PAWN || originSquare == BLACK_PAWN) {
		moveKind = EP_CAPTURE
	}

	if originSquare == WHITE_KING {
		switch dstCoord {
		case G1:
			if b.castleRights&CASTLING_RIGHTS_WHITE_KING_MASK > 0 {
				moveKind = KING_CASTLE
			}
		case C1:
			if b.castleRights&CASTLING_RIGHTS_WHITE_QUEEN_MASK > 0 {
				moveKind = QUEEN_CASTLE
			}
		}
	}

	if originSquare == BLACK_KING {
		switch dstCoord {
		case G8:
			if b.castleRights&CASTLING_RIGHTS_BLACK_KING_MASK > 0 {
				moveKind = KING_CASTLE
			}
		case C8:
			if b.castleRights&CASTLING_RIGHTS_BLACK_QUEEN_MASK > 0 {
				moveKind = QUEEN_CASTLE
			}
		}
	}

	if len(sParts) == 5 {
		switch sParts[4] {
		case 'n':
			if moveKind == CAPTURE {
				moveKind = KNIGHT_PROMOTION_CAPTURE
			} else {
				moveKind = KNIGHT_PROMOTION
			}
		case 'b':
			if moveKind == CAPTURE {
				moveKind = BISHOP_PROMOTION_CAPTURE
			} else {
				moveKind = BISHOP_PROMOTION
			}
		case 'r':
			if moveKind == CAPTURE {
				moveKind = ROOK_PROMOTION_CAPTURE
			} else {
				moveKind = ROOK_PROMOTION
			}
		case 'q':
			if moveKind == CAPTURE {
				moveKind = QUEEN_PROMOTION_CAPTURE
			} else {
				moveKind = QUEEN_PROMOTION
			}
		default:
			return 0, fmt.Errorf("Invalid promotion: %v", s)
		}
	}
	return NewMove(originCoord, dstCoord, originSquare, dstSquare, moveKind), nil
}

func (m *Move) ToUCIString() string {
	var mu MoveUnpacked
	m.Unpack(&mu)
	var p string
	switch mu.moveKind {
	case KNIGHT_PROMOTION_CAPTURE:
		fallthrough
	case KNIGHT_PROMOTION:
		p = "n"
	case BISHOP_PROMOTION_CAPTURE:
		fallthrough
	case BISHOP_PROMOTION:
		p = "b"
	case ROOK_PROMOTION_CAPTURE:
		fallthrough
	case ROOK_PROMOTION:
		p = "r"
	case QUEEN_PROMOTION_CAPTURE:
		fallthrough
	case QUEEN_PROMOTION:
		p = "q"
	}
	oc := COORD_STRINGS[mu.originCoord]
	dc := COORD_STRINGS[mu.dstCoord]
	s := fmt.Sprint(strings.ToLower(oc), strings.ToLower(dc), p)
	return s
}
