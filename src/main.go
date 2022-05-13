package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	seedKeys(time.Now().UTC().UnixNano())
	board, err := NewBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(board.ToString())
		fmt.Println()

		moves := board.generateMoves(board.sideToMove)

		var matchingMoves []Move
		for {
			fmt.Print("Submit target coord: ")
			scanner.Scan()
			input := scanner.Text()

			for _, m := range moves {
				if input == COORD_MAP[0] {
					matchingMoves = append(matchingMoves, m)
				}
			}
			if len(matchingMoves) == 0 {
				fmt.Println("*** Not a valid move ***")
			} else {
				break
			}
		}

		if len(matchingMoves) == 1 {
			err := board.makeMove(matchingMoves[0])
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		for {
			fmt.Print("Submit origin coord: ")
			scanner.Scan()
			input := scanner.Text()

			performedMoved := false

			for _, m := range matchingMoves {
				if input == COORD_MAP[0] {
					err := board.makeMove(m)
					if err != nil {
						fmt.Println(err)
					}
					performedMoved = true
				}
			}
			if performedMoved {
				break
			}
			fmt.Println("*** Not a valid move ***")
		}
	}
}
