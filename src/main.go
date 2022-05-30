package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
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

func ModePlayerVsEngine(scanner *bufio.Scanner) {
	search, err := NewSearch(STARTING_FEN, UCI_DEFAULT_DEPTH, 0)
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Printf("%v\n\n", search.ToString())
		fmt.Print("Submit move: ")

		moves := make([]Move, 0, INITIAL_MOVES_CAPACITY)
		search.GenerateMoves(&moves, search.sideToMove)

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
		score := search.Negamax(5, math.MinInt+1, math.MaxInt)
		search.MakeMove(search.bestMove)
		fmt.Printf("%v with score of %v\n", search.bestMove.ToString(), score)
	}

}
