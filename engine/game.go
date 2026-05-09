package engine

import (
	"math/rand"

	"github.com/user/golang-tetris/models"
)

const (
	// LockDelay is the number of ticks a piece can stay on the ground before locking.
	LockDelay = 2
	// MoveResetLimit prevents infinite rotations by capping resets at 15.
	MoveResetLimit = 15
)

// --- SCORING VALUES ---
// Edit these constants to change how points are awarded.
const (
	ScoreSingle   = 100
	ScoreDouble   = 300
	ScoreTriple   = 500
	ScoreTetris   = 800
	ScoreSoftDrop = 1 // Points per cell for soft drop
	ScoreHardDrop = 2 // Points per cell for hard drop
)

// GameState holds the main game state logic.
type GameState struct {
	Board            *Board
	CurrentPiece     models.Tetromino
	NextPiece        models.Tetromino
	Score            int
	GameOver         bool
	Shapes           [][][]int
	Colors           []int
	pieceBag         []int
	bagIndex         int
	LockTimer        int // Tracks how long a piece has been touching a surface
	ResetsCount      int // Tracks how many times lock timer was reset
	PendingDropScore int // Points accumulated from drops, awarded only if lines clear
}

// NewGame initializes a new game state.
func NewGame() *GameState {
	gs := &GameState{
		Board:  NewBoard(),
		Shapes: models.GetShapes(),
		// Termbox colors: 1-7 are standard colors
		Colors: []int{2, 3, 4, 5, 6, 7, 8},
	}
	gs.fillBag()
	gs.NextPiece = gs.generateRandomPiece()
	gs.SpawnPiece()
	return gs
}

func (g *GameState) fillBag() {
	g.pieceBag = rand.Perm(len(g.Shapes))
	g.bagIndex = 0
}

func (g *GameState) nextBagIndex() int {
	if g.bagIndex >= len(g.pieceBag) {
		g.fillBag()
	}
	idx := g.pieceBag[g.bagIndex]
	g.bagIndex++
	return idx
}

func (g *GameState) generateRandomPiece() models.Tetromino {
	typeIdx := g.nextBagIndex()
	return models.Tetromino{
		Shape: g.Shapes[typeIdx],
		Color: g.Colors[typeIdx],
	}
}

// SpawnPiece creates a new random Tetromino at the top of the board.
func (g *GameState) SpawnPiece() {
	g.CurrentPiece = g.NextPiece
	g.CurrentPiece.X = BoardWidth/2 - len(g.CurrentPiece.Shape[0])/2
	g.CurrentPiece.Y = 0
	g.LockTimer = 0
	g.ResetsCount = 0
	g.PendingDropScore = 0

	g.NextPiece = g.generateRandomPiece()

	// Check for immediate collision (Game Over)
	if g.Board.IsCollision(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y, g.CurrentPiece.Shape) {
		g.GameOver = true
	}
}

// MovePiece attempts to move the current piece. Returns true if successful.
func (g *GameState) MovePiece(dx, dy int) bool {
	if g.GameOver {
		return false
	}

	newX := g.CurrentPiece.X + dx
	newY := g.CurrentPiece.Y + dy

	if !g.Board.IsCollision(g.CurrentPiece, newX, newY, g.CurrentPiece.Shape) {
		g.CurrentPiece.X = newX
		g.CurrentPiece.Y = newY
		// Reset lock timer if touching surface and under limit
		if g.IsTouchingBottom() && g.ResetsCount < MoveResetLimit {
			g.LockTimer = 0
			g.ResetsCount++
		}
		return true
	}

	return false
}

// SoftDrop moves the current piece down and tracks potential bonus.
func (g *GameState) SoftDrop() bool {
	if g.MovePiece(0, 1) {
		g.PendingDropScore += ScoreSoftDrop
		return true
	}
	return false
}

// HardDrop instantly drops the current piece to the bottom and tracks potential bonus.
func (g *GameState) HardDrop() {
	if g.GameOver {
		return
	}

	for g.MovePiece(0, 1) {
		g.PendingDropScore += ScoreHardDrop
	}

	g.LockPiece()
}

// RotatePiece attempts to rotate the current piece with simple wall kicks.
func (g *GameState) RotatePiece() {
	if g.GameOver {
		return
	}

	newShape := models.Rotate(g.CurrentPiece.Shape)

	// Wall kick offsets
	offsets := []models.Vector2{
		{X: 0, Y: 0},
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: -2, Y: 0},
		{X: 2, Y: 0},
	}

	for _, offset := range offsets {
		newX := g.CurrentPiece.X + offset.X
		newY := g.CurrentPiece.Y + offset.Y
		if !g.Board.IsCollision(g.CurrentPiece, newX, newY, newShape) {
			g.CurrentPiece.X = newX
			g.CurrentPiece.Y = newY
			g.CurrentPiece.Shape = newShape
			if g.IsTouchingBottom() && g.ResetsCount < MoveResetLimit {
				g.LockTimer = 0
				g.ResetsCount++
			}
			return
		}
	}
}

// IsTouchingBottom checks if the current piece is resting on something.
func (g *GameState) IsTouchingBottom() bool {
	return g.Board.IsCollision(g.CurrentPiece, g.CurrentPiece.X, g.CurrentPiece.Y+1, g.CurrentPiece.Shape)
}

// GetGhostY calculates the Y position where the current piece would land.
func (g *GameState) GetGhostY() int {
	ghostY := g.CurrentPiece.Y
	for !g.Board.IsCollision(g.CurrentPiece, g.CurrentPiece.X, ghostY+1, g.CurrentPiece.Shape) {
		ghostY++
	}
	return ghostY
}

func scoreForLines(lines int) int {
	switch lines {
	case 1:
		return ScoreSingle
	case 2:
		return ScoreDouble
	case 3:
		return ScoreTriple
	case 4:
		return ScoreTetris
	default:
		return lines * 200
	}
}

// LockPiece fixes the current piece to the board and prepares for the next.
func (g *GameState) LockPiece() int {
	g.Board.PlacePiece(g.CurrentPiece)
	cleared := g.Board.ClearLines()
	g.Score += scoreForLines(cleared)
	
	// Bonus drop score is ONLY awarded if lines were cleared
	if cleared > 0 {
		g.Score += g.PendingDropScore
	}
	
	g.SpawnPiece()
	return cleared
}

// Tick handles the gravity update.
func (g *GameState) Tick() {
	if g.GameOver {
		return
	}

	if !g.MovePiece(0, 1) {
		g.LockTimer++
		if g.LockTimer >= LockDelay {
			g.LockPiece()
		}
	} else {
		g.LockTimer = 0
	}
}
