package chessboard

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/notnil/chess"
)

// Draws the board from FEN notation with the correct orientation
func DrawBoard(board *chess.Board, isBlackTurn bool) {
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