package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	Setup()
	HandleInput()
}

func HandleInput() {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if input == "uci" {
		ModeUCI(scanner)
	} else {
		ModePlayerVsEngine(scanner)
	}

}

func Setup() {
	SeedKeys(time.Now().UTC().UnixNano())
}

func ModePlayerVsEngine(scanner *bufio.Scanner) {
	board, err := NewBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Printf("%v\n\n", board.ToString())
		fmt.Print("Submit move: ")

		moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		board.GenerateMoves(&moves, board.sideToMove)

		for scanner.Scan() {
			input := scanner.Text()
			submittedMove, err := board.ParseUCIMove(input)
			if err != nil {
				fmt.Println(err)
			}

			for _, move := range moves {
				if submittedMove == move {
					err := board.MakeMove(move)
					if err != nil {
						board.UndoMove()
						break
					}
					goto SEARCH
				}
			}
			fmt.Println("*** Not a valid move ***")
			fmt.Print("Submit move: ")
		}

	SEARCH:
		fmt.Printf("%v\n\n", board.ToString())
		search := Search{Board: board}
		score := search.Negamax(5, math.MinInt+1, math.MaxInt)
		board.MakeMove(search.bestMove)
		fmt.Printf("%v with score of %v\n", search.bestMove.ToString(), score)
	}

}
