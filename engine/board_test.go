package engine

import (
	"testing"

	"github.com/user/golang-tetris/models"
)

func TestIsCollision(t *testing.T) {
	board := NewBoard()
	piece := models.Tetromino{
		Shape: [][]int{
			{1, 1},
			{1, 1},
		},
		X: 0,
		Y: 0,
	}

	// Test boundary collision - Left
	if !board.IsCollision(piece, -1, 0, piece.Shape) {
		t.Errorf("Expected collision at X=-1, but none found")
	}

	// Test boundary collision - Right
	if !board.IsCollision(piece, BoardWidth-1, 0, piece.Shape) {
		t.Errorf("Expected collision at X=BoardWidth-1 (overflow), but none found")
	}

	// Test boundary collision - Bottom
	if !board.IsCollision(piece, 0, BoardHeight-1, piece.Shape) {
		t.Errorf("Expected collision at Y=BoardHeight-1 (overflow), but none found")
	}

	// Test block collision
	board.Grid[5][5] = 1
	if !board.IsCollision(piece, 4, 4, piece.Shape) {
		t.Errorf("Expected collision with placed block at (5,5), but none found")
	}
}

func TestIsCollisionAboveBoard(t *testing.T) {
	board := NewBoard()
	piece := models.Tetromino{
		Shape: [][]int{
			{0, 1, 0},
			{1, 1, 1},
		},
		X: 4,
		Y: -1,
	}

	if board.IsCollision(piece, piece.X, piece.Y, piece.Shape) {
		t.Errorf("Expected no collision when part of the piece is above the visible board")
	}
}

func TestClearLines(t *testing.T) {
	board := NewBoard()

	// Fill the bottom line
	for x := 0; x < BoardWidth; x++ {
		board.Grid[BoardHeight-1][x] = 1
	}

	// Fill half of the next line
	for x := 0; x < BoardWidth/2; x++ {
		board.Grid[BoardHeight-2][x] = 1
	}

	cleared := board.ClearLines()
	if cleared != 1 {
		t.Errorf("Expected 1 line cleared, got %d", cleared)
	}

	// Check if the half line shifted down
	if board.Grid[BoardHeight-1][0] != 1 || board.Grid[BoardHeight-1][BoardWidth/2] != 0 {
		t.Errorf("Line shifting logic failed")
	}

	// Check if the top line is cleared
	if board.Grid[0][0] != 0 {
		t.Errorf("Top line should be empty after clear")
	}
}
