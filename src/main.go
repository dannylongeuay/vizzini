package main

import "fmt"

func main() {
	board, err := newBoard(StartingFEN)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(board.toString())
}
