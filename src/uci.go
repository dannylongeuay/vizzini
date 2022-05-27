package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

func UCIMode() {
	fmt.Println("id name Vizzini Alpha")
	fmt.Println("id author Daniel")
	// fmt.Println("option")
	fmt.Println("uciok")
	var board *Board

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		tokens := strings.Split(input, " ")

		if len(tokens) == 0 {
			continue
		}

		switch tokens[0] {
		case "isready":
			fmt.Println("readyok")
		case "setoption":
		case "ucinewgame":
			var err error
			board, err = NewBoard(STARTING_FEN)
			if err != nil {
				panic(err)
			}
		case "position":
			var err error
			board, err = UCISetPosition(tokens[1:])
			if err != nil {
				panic(err)
			}
		case "go":
			// fmt.Println(board.ToString())
			search := Search{Board: board}
			search.Negamax(3, math.MinInt+1, math.MaxInt)
			fmt.Println("bestmove", search.bestMove.ToUCIString())
		case "stop":

		}
	}

}

func UCISetPosition(tokens []string) (*Board, error) {
	fen := STARTING_FEN
	moveTokenIndex := 1
	// fmt.Println(tokens)
	if tokens[0] == "fen" {
		moveTokenIndex = 7
		fen = strings.Join(tokens[1:moveTokenIndex], " ")
	}
	board, err := NewBoard(fen)
	if err != nil {
		return board, err
	}
	if len(tokens) > moveTokenIndex && tokens[moveTokenIndex] == "moves" {
		for _, uciMove := range tokens[moveTokenIndex+1:] {
			move, err := board.UCIParseMove(uciMove)
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

func (b *Board) UCIParseMove(s string) (Move, error) {
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
