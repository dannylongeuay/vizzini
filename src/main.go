package main

import (
	"fmt"
	"time"
)

func main() {
	seedKeys(time.Now().UTC().UnixNano())
	board, err := newBoard(STARTING_FEN)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(board.toString())
}
