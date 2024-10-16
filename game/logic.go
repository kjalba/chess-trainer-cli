package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kjalba/chess-trainer-cli/chessboard"
	"github.com/kjalba/chess-trainer-cli/utils"
	"github.com/notnil/chess"
)

func HandleUserInput(game *chess.Game, solution []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s to move.\n\n", utils.GetColor(game.Position().Turn()))

	for i := 0; i < len(solution); {
		moveProcessed := false // Flag to track if we should move to the next solution move
		lastMove := false      // Flag to track if we are at the last move in solution
		if i == len(solution)-1 {
			lastMove = true
		}
		correctMoveUCI := solution[i]

		// Parses the UCI move directly and apply it
		move, err := chess.UCINotation{}.Decode(game.Position(), correctMoveUCI)
		if err != nil {
			fmt.Printf("Failed to parse UCI move: %v\n", err)
			return
		}

		// Converts the UCI move to algebraic (SAN) format to allow user to enter this format as well
		correctMoveAlgebraic := utils.ConvertChessMoveToAlgebraic(game, move)
		var opponentMoveUCI string
		if !lastMove {
			opponentMoveUCI = solution[i+1]
		}

		fmt.Print("Enter the best move, or '?' for help: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		userInput = strings.ToLower(userInput)
		fmt.Print("\n")
		switch userInput {
		case "?", "help":
			utils.PrintHelp()
		case "":
			fmt.Printf("The correct move was %s.\n", correctMoveAlgebraic)
			// Applies the correct move in UCI format
			game.Move(move)

			if !lastMove {
				// Applies the opponent's move in UCI format
				opponentMove, err := chess.UCINotation{}.Decode(game.Position(), opponentMoveUCI)
				if err != nil {
					fmt.Printf("Failed to parse opponent UCI move: %v\n", err)
					return
				}
				opponentMoveAlgebraic := utils.ConvertChessMoveToAlgebraic(game, opponentMove)
				fmt.Printf("Opponent played %s.\n", opponentMoveAlgebraic)
				game.Move(opponentMove)
			}
			moveProcessed = true
		case "s", "show":
				chessboard.DrawBoard(game.Position().Board(), game.Position().Turn() == chess.Black)
		default:
			// Allows user to enter either UCI or SAN notation
			if userInput == correctMoveUCI || userInput == strings.ToLower(correctMoveAlgebraic) {
				fmt.Println("Correct!")
				fmt.Print("\n")
				game.Move(move)

				if !lastMove {
					opponentMove, err := chess.UCINotation{}.Decode(game.Position(), opponentMoveUCI)
					if err != nil {
						fmt.Printf("Failed to parse opponent UCI move: %v\n", err)
						return
					}
					opponentMoveAlgebraic := utils.ConvertChessMoveToAlgebraic(game, opponentMove)
					fmt.Printf("Opponent played %s.\n", opponentMoveAlgebraic)
					game.Move(opponentMove)
				}
				moveProcessed = true
			} else {
				fmt.Printf("Incorrect! Try again.\n\n")
			}
		}

		if moveProcessed {
			i += 2 // This works because the solution will always have an odd number of moves since the user always finishes the puzzle with a move
			if lastMove {
				chessboard.DrawBoard(game.Position().Board(), game.Position().Turn() != chess.Black)
			}
		}
	}

	fmt.Println("Puzzle completed.")
}