package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kjalba/chess-trainer-cli/chessboard"
	"github.com/kjalba/chess-trainer-cli/game"
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
	dailyFlag := flag.Bool("dailyPuzzle", false, "Fetch the daily puzzle from Lichess")
	randomFlag := flag.Bool("randomPuzzle", false, "Fetch a random puzzle based on theme and/or rating")

	flag.Parse()

	if *randomFlag {
		fmt.Println("Fetching a random puzzle...")
		handleRandomPuzzle()
	} else if *dailyFlag {
		fmt.Println("Fetching the daily puzzle...")
		handleDailyPuzzle()
	} else {
		fmt.Println("Fetching the daily puzzle...")
		handleDailyPuzzle()
	}
}

func handleDailyPuzzle() {
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

	lichessGame, puzzle := lichess.Game, lichess.Puzzle

	chessGame := chess.NewGame()
	err = chessGame.UnmarshalText([]byte(lichessGame.PGN))
	if err != nil {
		log.Fatalf("Failed to parse PGN: %v", err)
	}

	chessGameBoard := chessGame.Position().Board()

	// The notnil/chess library doesn't have a function for reversing the board's row and column headers in the drawing, so I implemented my own draw function
	chessboard.DrawBoard(chessGameBoard, chessGame.Position().Turn() == chess.Black)
	game.HandleUserInput(chessGame, puzzle.Solution)
}

func handleRandomPuzzle() {
	// Logic to handle fetching random puzzles can be implemented here
	fmt.Println("Random puzzle feature is under development.")
}
