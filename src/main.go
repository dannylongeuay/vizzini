package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	seedKeys(time.Now().UTC().UnixNano())
	board, err := newBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(board.toString())
		fmt.Print("\n\nMoves:\n\n")
		moves := board.generateMoves(board.sideToMove)
		for i, m := range moves {
			fmt.Println(i, coordBySquareIndex(m.origin), coordBySquareIndex(m.target))
		}
		fmt.Print("\n\nSubmit move: ")
		scanner.Scan()
		input := scanner.Text()

		numMove, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
		} else if numMove >= 0 && numMove < len(moves) {
			merr := board.makeMove(moves[numMove])
			if merr != nil {
				fmt.Println(merr)
			}
		} else {
			fmt.Println("Not a valid move.")
		}
	}
}
