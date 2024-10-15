package main

import (
	"flag"
	"fmt"

	"github.com/kjalba/chess-trainer-cli/puzzle"
)

var (
	dailyPuzzle bool
	randomPuzzle bool
)

func init() {
	const (
		defaultDaily = false
		dailyUsage = "Fetch the daily puzzle from Lichess"
		defaultRandom = false
		randomUsage = "Fetch a random puzzle based on theme and/or rating"
	)
	flag.BoolVar(&dailyPuzzle, "dailyPuzzle", defaultDaily, dailyUsage)
	flag.BoolVar(&dailyPuzzle, "dp", defaultDaily, dailyUsage+" (shorthand)")

	flag.BoolVar(&randomPuzzle, "randomPuzzle", defaultRandom, randomUsage)
	flag.BoolVar(&randomPuzzle, "rp", defaultRandom, randomUsage+" (shorthand)")
}


func main() {
	flag.Parse()
	if randomPuzzle {
		fmt.Println("Fetching a random puzzle...")
		puzzle.HandleRandomPuzzle()
	} else if dailyPuzzle {
		fmt.Println("Fetching the daily puzzle...")
		puzzle.HandleDailyPuzzle()
	} else {
		fmt.Println("Fetching the daily puzzle...")
		puzzle.HandleDailyPuzzle()
	}
}
