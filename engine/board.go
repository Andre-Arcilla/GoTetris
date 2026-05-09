package engine

import (
	"github.com/user/golang-tetris/models"
)

const (
	BoardWidth  = 10
	BoardHeight = 20
)

// Board manages the game grid and line clearing logic.
type Board struct {
	Grid [BoardHeight][BoardWidth]int
}

// NewBoard creates a new empty board.
func NewBoard() *Board {
	return &Board{}
}

// IsCollision checks if the piece at the given position conflicts with board boundaries or existing blocks.
func (b *Board) IsCollision(piece models.Tetromino, x, y int, shape [][]int) bool {
	for i, row := range shape {
		for j, cell := range row {
			if cell != 0 {
				newX := x + j
				newY := y + i

				// Check horizontal bounds and bottom boundary.
				if newX < 0 || newX >= BoardWidth || newY >= BoardHeight {
					return true
				}

				// Any part above the top of the visible board is allowed for spawn/rotation.
				if newY < 0 {
					continue
				}

				// Check existing blocks
				if b.Grid[newY][newX] != 0 {
					return true
				}
			}
		}
	}
	return false
}

// PlacePiece adds the piece's blocks to the board's grid.
func (b *Board) PlacePiece(piece models.Tetromino) {
	for i, row := range piece.Shape {
		for j, cell := range row {
			if cell != 0 {
				newY := piece.Y + i
				if newY >= 0 && newY < BoardHeight {
					b.Grid[newY][piece.X+j] = piece.Color
				}
			}
		}
	}
}

// ClearLines removes full horizontal lines and returns the number of lines cleared.
func (b *Board) ClearLines() int {
	linesCleared := 0
	for i := BoardHeight - 1; i >= 0; i-- {
		isFull := true
		for j := 0; j < BoardWidth; j++ {
			if b.Grid[i][j] == 0 {
				isFull = false
				break
			}
		}

		if isFull {
			linesCleared++
			// Shift lines down
			for k := i; k > 0; k-- {
				b.Grid[k] = b.Grid[k-1]
			}
			// Clear top line
			b.Grid[0] = [BoardWidth]int{}
			// Re-check the same line index as it now contains the shifted line
			i++
		}
	}
	return linesCleared
}
