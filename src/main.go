package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	SeedKeys(time.Now().UTC().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if input == "uci" {
		UCIMode()
	} else {
		PlaySinglePlayer()
	}
}

func PlaySinglePlayer() {
	board, err := NewBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%v\n\n", board.ToString())

		moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		board.GenerateMoves(&moves, board.sideToMove)

		for {
			fmt.Print("Submit move: ")
			scanner.Scan()
			input := scanner.Text()
			submittedMove, err := board.UCIParseMove(input)
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
		}

	SEARCH:
		fmt.Printf("%v\n\n", board.ToString())
		search := Search{Board: board}
		score := search.Negamax(5, math.MinInt+1, math.MaxInt)
		board.MakeMove(search.bestMove)
		fmt.Printf("%v with score of %v\n", search.bestMove.ToString(), score)
	}

}
