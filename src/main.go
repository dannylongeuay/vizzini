package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	SeedKeys(time.Now().UTC().UnixNano())
	board, err := NewBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(os.Stdin)

MAIN_LOOP:
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
					goto MAIN_LOOP
				}
			}
			fmt.Println("*** Not a valid move ***")
		}
	}
}
