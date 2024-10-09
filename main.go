package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/notnil/chess"
)

type Lichess struct {
	Game struct {
		ID string `json:"id"`
		PGN string `json:"pgn"`
	} `json:"game"`
	
	Puzzle struct {
		ID string `json:"id"`
		Solution []string `json:"solution"`
		Themes []string `json:"themes"`
	} `json:"puzzle"`
}


func main() {
	res, err := http.Get("https://lichess.org/api/puzzle/daily") // daily puzzle api 
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		 log.Fatalf("Lichess API returned status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var lichess Lichess
	err = json.Unmarshal(body, &lichess)
	if err != nil {
		log.Fatal(err)
	}

	game, puzzle := lichess.Game, lichess.Puzzle

	chessGame := chess.NewGame()
	err = chessGame.UnmarshalText([]byte(game.PGN))
	if err != nil {
		log.Fatalf("Failed to parse PGN: %v", err)
	}

	chessGameBoard := chessGame.Position().Board()


	isBlackTurn := chessGame.Position().Turn() == chess.Black

	// The notnil/chess library doesn't have a function for reversing the board's row and column headers in the drawing, so I implemented my own draw function
	drawBoard(chessGameBoard, isBlackTurn)

	handleUserInput(chessGame, puzzle.Solution)
}

// Draws the board from FEN notation with the correct orientation
func drawBoard(board *chess.Board, isBlackTurn bool) {
	if isBlackTurn {
		fmt.Printf("\n")
		fmt.Println("  h g f e d c b a")
		fmt.Println(" +----------------+")
		for row := 0; row < 8; row++ {
			fmt.Printf("%d|", row+1)
			for col := 7; col >= 0; col-- {
				square := chess.Square((row*8 + col))
				piece := board.Piece(square)
				fmt.Printf("%s ", pieceASCII(piece))
			}
			fmt.Printf("|%d\n", row+1)
		}
		fmt.Println(" +----------------+")
		fmt.Println("  h g f e d c b a")
		fmt.Printf("\n")
	} else {
		fmt.Printf("\n")
		fmt.Println("  a b c d e f g h")
		fmt.Println(" +----------------+")
		for row := 7; row >= 0; row-- {
			fmt.Printf("%d|", row+1)
			for col := 0; col < 8; col++ {
				square := chess.Square((row*8 + col))
				piece := board.Piece(square)
				fmt.Printf("%s ", pieceASCII(piece))
			}
			fmt.Printf("|%d\n", row+1)
		}
		fmt.Println(" +----------------+")
		fmt.Println("  a b c d e f g h")
		fmt.Printf("\n")
	}

}

// Converts chess pieces to ASCII and color them
func pieceASCII(piece chess.Piece) string {
	if piece == chess.NoPiece {
		return "."
	}

	pieceChar := piece.String() // The piece's string representation like 'P', 'N', etc.

	switch piece.Color() {
	case chess.White:
		return color.New(color.FgBlue).SprintFunc()(pieceChar)
	case chess.Black:
		return color.New(color.FgRed).SprintFunc()(pieceChar)
	default:
		return pieceChar
	}
}

func handleUserInput(game *chess.Game, solution []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s to move.\n\n", getColor(game.Position().Turn()))

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
		correctMoveAlgebraic := chess.AlgebraicNotation{}.Encode(game.Position(), move)
		correctMoveAlgebraic = strings.ToLower(correctMoveAlgebraic)
		
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
			printHelp()
		case "":
			fmt.Printf("The correct move was %s.\n", correctMoveUCI)
			// Applies the correct move in UCI format
			game.Move(move)

			if !lastMove {
				// Applies the opponent's move in UCI format
				opponentMove, err := chess.UCINotation{}.Decode(game.Position(), opponentMoveUCI)
				if err != nil {
					fmt.Printf("Failed to parse opponent UCI move: %v\n", err)
					return
				}
				fmt.Printf("Opponent played %s.\n", opponentMove)
				game.Move(opponentMove)
			}
			moveProcessed = true
		default:
			// Allows user to enter either UCI or SAN notation
			if userInput == correctMoveUCI || userInput == correctMoveAlgebraic {
				fmt.Println("Correct!")
				fmt.Print("\n")
				game.Move(move)

				if !lastMove {
					opponentMove, err := chess.UCINotation{}.Decode(game.Position(), opponentMoveUCI)
					if err != nil {
						fmt.Printf("Failed to parse opponent UCI move: %v\n", err)
						return
					}
					fmt.Printf("Opponent played %s.\n", opponentMove)
					game.Move(opponentMove)
				}
				moveProcessed = true
			} else {
				fmt.Printf("Incorrect! Try again.\n\n")
			}
		}

		if moveProcessed {
			i += 2 // This works because the solution will always have an odd number of moves since the user always finishes the puzzle with a move
			if ! lastMove {
				drawBoard(game.Position().Board(), game.Position().Turn() == chess.Black)
			} else {
				drawBoard(game.Position().Board(), game.Position().Turn() != chess.Black)
			}
		}
	}

	fmt.Println("Puzzle completed.")
}

func printHelp() {
	fmt.Println("\nHelp:")
	fmt.Println(" - Enter your move (e.g., 'e4', 'Nf3', 'Qxd7').")
	fmt.Println(" - Enter '?' to display this help message.")
	fmt.Println(" - Enter nothing to skip the move and reveal the correct answer.")
	fmt.Println()
}

func getColor(turn chess.Color) string {
	if turn == chess.White {
		return "White"
	}
	return "Black"
}

