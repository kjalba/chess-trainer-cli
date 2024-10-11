package utils

import (
	"fmt"

	"github.com/notnil/chess"
)

func PrintHelp() {
	fmt.Println("\nHelp:")
	fmt.Println(" - Enter your move (e.g., 'e4', 'Nf3', 'Qxd7').")
	fmt.Println(" - Enter '?' to display this help message.")
	fmt.Println(" - Enter nothing to skip the move and reveal the correct answer.")
	fmt.Println()
}

func GetColor(turn chess.Color) string {
	if turn == chess.White {
		return "White"
	}
	return "Black"
}

func ConvertChessMoveToAlgebraic(game *chess.Game, move *chess.Move) string {
	return chess.AlgebraicNotation{}.Encode(game.Position(), move)
}