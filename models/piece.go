package models

// Tetromino represents a falling piece in the game.
type Tetromino struct {
	Shape    [][]int
	X, Y     int
	Color    int // Using int to represent termbox.Attribute
}

// Vector2 represents a position on the board.
type Vector2 struct {
	X, Y int
}

// GetShapes returns the matrix representations of all Tetromino types.
func GetShapes() [][][]int {
	return [][][]int{
		// I
		{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		// J
		{
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		// L
		{
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 0},
		},
		// O
		{
			{1, 1},
			{1, 1},
		},
		// S
		{
			{0, 1, 1},
			{1, 1, 0},
			{0, 0, 0},
		},
		// T
		{
			{0, 1, 0},
			{1, 1, 1},
			{0, 0, 0},
		},
		// Z
		{
			{1, 1, 0},
			{0, 1, 1},
			{0, 0, 0},
		},
	}
}

// Rotate returns a new rotated matrix for the Tetromino.
// It rotates the shape 90 degrees clockwise.
func Rotate(shape [][]int) [][]int {
	n := len(shape)
	m := len(shape[0])
	newShape := make([][]int, m)
	for i := range newShape {
		newShape[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			newShape[j][n-1-i] = shape[i][j]
		}
	}
	return newShape
}
