package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	switch os.Args[1] {
	case "play":
		ModePlayerVsEngine(scanner)
	case "uci":
		ModeUCI(scanner)
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: vizzini <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  play   Start Player vs Engine mode")
	fmt.Println("  uci    Start UCI mode")
}

func ModePlayerVsEngine(scanner *bufio.Scanner) {
	search, err := NewSearch(STARTING_FEN, DEFAULT_MAX_DEPTH, DEFAULT_MAX_NODES)
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Printf("%v\n\n", search.ToString())
		fmt.Print("Submit move: ")

		moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		search.GenerateMoves(&moves, search.sideToMove, false)

		for scanner.Scan() {
			input := scanner.Text()
			submittedMove, err := search.ParseUCIMove(input)
			if err != nil {
				fmt.Println(err)
			}

			for _, move := range moves {
				if submittedMove == move {
					err := search.MakeMove(move)
					if err != nil {
						search.UndoMove()
						break
					}
					goto SEARCH
				}
			}
			fmt.Println("*** Not a valid move ***")
			fmt.Print("Submit move: ")
		}

	SEARCH:
		fmt.Printf("%v\n\n", search.ToString())
		search.Reset()
		search.stopTime = time.Now().Add(2500 * time.Millisecond)
		score := search.IterativeDeepening()
		search.MakeMove(search.bestMove)
		fmt.Printf("%v with score of %v\n", search.bestMove.ToString(), score)
	}

}
