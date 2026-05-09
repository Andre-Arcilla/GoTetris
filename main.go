package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/user/golang-tetris/engine"
)

const (
	UIOffsetX = 10
	UIOffsetY = 2
)

type GameMode int

const (
	ModePlaying GameMode = iota
	ModeGameOver
	ModeScoreEntry
	ModeHighScores
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event, 10)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	highScores := engine.LoadHighScores()
	
	for {
		if !runGameLoop(&highScores, eventQueue) {
			break
		}
	}
}

func runGameLoop(highScores *engine.HighScoreBoard, eventQueue chan termbox.Event) bool {
	game := engine.NewGame()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	mode := ModePlaying
	playerName := ""

	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		switch mode {
		case ModePlaying:
			render(game)
			if game.GameOver {
				if highScores.IsHighScore(game.Score) {
					mode = ModeScoreEntry
				} else {
					mode = ModeHighScores
				}
			}
		case ModeScoreEntry:
			drawScoreEntry(game.Score, playerName)
		case ModeHighScores:
			drawHighScores(highScores)
		}

		termbox.Flush()

		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
					return false // Quit application
				}

				switch mode {
				case ModePlaying:
					handlePlayingInput(game, ev)
				case ModeScoreEntry:
					if ev.Key == termbox.KeyEnter && len(playerName) > 0 {
						highScores.AddScore(playerName, game.Score)
						highScores.Save()
						mode = ModeHighScores
					} else if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
						if len(playerName) > 0 {
							playerName = playerName[:len(playerName)-1]
						}
					} else if ev.Ch != 0 && len(playerName) < 3 {
						playerName += string(ev.Ch)
					}
				case ModeHighScores:
					return true // Any key to restart
				}
			}
		case <-ticker.C:
			if mode == ModePlaying {
				game.Tick()
			}
		}
	}
}

func handlePlayingInput(game *engine.GameState, ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyArrowLeft:
		game.MovePiece(-1, 0)
	case termbox.KeyArrowRight:
		game.MovePiece(1, 0)
	case termbox.KeyArrowDown:
		game.SoftDrop()
	case termbox.KeyArrowUp:
		game.RotatePiece()
	case termbox.KeySpace:
		game.HardDrop()
	}
}

func render(g *engine.GameState) {
	drawRect(UIOffsetX-1, UIOffsetY-1, engine.BoardWidth*2+2, engine.BoardHeight+2)
	for y := 0; y < engine.BoardHeight; y++ {
		for x := 0; x < engine.BoardWidth; x++ {
			cell := g.Board.Grid[y][x]
			if cell != 0 {
				drawCell(x*2+UIOffsetX, y+UIOffsetY, termbox.Attribute(cell))
			} else {
				termbox.SetCell(x*2+UIOffsetX, y+UIOffsetY, '.', termbox.ColorDarkGray, termbox.ColorDefault)
				termbox.SetCell(x*2+UIOffsetX+1, y+UIOffsetY, ' ', termbox.ColorDarkGray, termbox.ColorDefault)
			}
		}
	}

	ghostY := g.GetGhostY()
	if !g.GameOver {
		for i, row := range g.CurrentPiece.Shape {
			for j, cell := range row {
				if cell != 0 {
					drawGhostCell(g.CurrentPiece.X*2+j*2+UIOffsetX, ghostY+i+UIOffsetY, termbox.Attribute(g.CurrentPiece.Color))
				}
			}
		}
		for i, row := range g.CurrentPiece.Shape {
			for j, cell := range row {
				if cell != 0 {
					drawCell(g.CurrentPiece.X*2+j*2+UIOffsetX, g.CurrentPiece.Y+i+UIOffsetY, termbox.Attribute(g.CurrentPiece.Color))
				}
			}
		}
	}

	sideX := UIOffsetX + engine.BoardWidth*2 + 4
	drawText(sideX, UIOffsetY, fmt.Sprintf("Score: %d", g.Score))
	drawText(sideX, UIOffsetY+2, "Next:")
	for i, row := range g.NextPiece.Shape {
		for j, cell := range row {
			if cell != 0 {
				drawCell(sideX+j*2, UIOffsetY+3+i, termbox.Attribute(g.NextPiece.Color))
			}
		}
	}
	drawText(sideX, UIOffsetY+8, "Controls:")
	drawText(sideX, UIOffsetY+9, "Arrows: Move/Rotate")
	drawText(sideX, UIOffsetY+10, "Space: Hard Drop")
	drawText(sideX, UIOffsetY+11, "Esc: Quit")
}

func drawCell(x, y int, color termbox.Attribute) {
	termbox.SetCell(x, y, '[', termbox.ColorWhite, color)
	termbox.SetCell(x+1, y, ']', termbox.ColorWhite, color)
}

func drawGhostCell(x, y int, color termbox.Attribute) {
	termbox.SetCell(x, y, '{', color, termbox.ColorDefault)
	termbox.SetCell(x+1, y, '}', color, termbox.ColorDefault)
}

func drawRect(x, y, w, h int) {
	for i := 0; i < w; i++ {
		termbox.SetCell(x+i, y, '-', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x+i, y+h-1, '-', termbox.ColorWhite, termbox.ColorDefault)
	}
	for i := 0; i < h; i++ {
		termbox.SetCell(x, y+i, '|', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x+w-1, y+i, '|', termbox.ColorWhite, termbox.ColorDefault)
	}
}

func drawText(x, y int, text string) {
	for i, r := range text {
		termbox.SetCell(x+i, y, r, termbox.ColorWhite, termbox.ColorDefault)
	}
}

func drawGameOverOverlay() {
	text := " GAME OVER "
	sub := " Press 'R' to Restart "
	x := UIOffsetX + engine.BoardWidth - len(text)/2
	y := UIOffsetY + engine.BoardHeight/2
	drawText(x, y, text)
	drawText(UIOffsetX+engine.BoardWidth-len(sub)/2, y+1, sub)
}

func drawScoreEntry(score int, name string) {
	header := " NEW HIGH SCORE! "
	prompt := fmt.Sprintf("Score: %d", score)
	input := "Enter Name: " + name + "_"
	instruct := "Press ENTER to save"
	
	x := UIOffsetX + engine.BoardWidth - 10
	y := UIOffsetY + engine.BoardHeight/2 - 2
	
	drawText(x, y, header)
	drawText(x, y+1, prompt)
	drawText(x, y+3, input)
	drawText(x, y+5, instruct)
}

func drawHighScores(board *engine.HighScoreBoard) {
	header := " HIGH SCORES "
	headerX := UIOffsetX + engine.BoardWidth - len(header)/2
	drawText(headerX, UIOffsetY+2, header)
	
	// Alignment:
	// 1. XXX 000000
	// ^  ^   ^
	// Rank Name Score
	
	startX := UIOffsetX + engine.BoardWidth - 8
	for i, s := range board.Scores {
		rank := fmt.Sprintf("%2d.", i+1)
		name := fmt.Sprintf("%-3s", s.Name)
		score := fmt.Sprintf("%d", s.Score)
		
		drawText(startX, UIOffsetY+4+i, rank)
		drawText(startX+4, UIOffsetY+4+i, name)
		drawText(startX+9, UIOffsetY+4+i, score)
	}
	
	footer := " Press any key to Play Again "
	drawText(UIOffsetX+engine.BoardWidth-len(footer)/2, UIOffsetY+16, footer)
}
