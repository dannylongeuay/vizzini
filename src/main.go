package main

import "fmt"

func main() {
	board, err := newBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(board.toString())
}
