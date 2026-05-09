package engine

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type ScoreRecord struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type HighScoreBoard struct {
	Scores []ScoreRecord `json:"scores"`
}

const HighScoreFile = "highscores.json"

func LoadHighScores() HighScoreBoard {
	var board HighScoreBoard
	data, err := os.ReadFile(HighScoreFile)
	if err != nil {
		return HighScoreBoard{Scores: []ScoreRecord{}}
	}
	if err := json.Unmarshal(data, &board); err != nil {
		return HighScoreBoard{Scores: []ScoreRecord{}}
	}
	return board
}

func (b *HighScoreBoard) AddScore(name string, score int) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Player"
	}
	if len(name) > 10 {
		name = name[:10]
	}

	b.Scores = append(b.Scores, ScoreRecord{Name: name, Score: score})
	sort.Slice(b.Scores, func(i, j int) bool {
		return b.Scores[i].Score > b.Scores[j].Score
	})
	if len(b.Scores) > 10 {
		b.Scores = b.Scores[:10]
	}
}

func (b *HighScoreBoard) Save() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(HighScoreFile, data, 0644)
}

func (b *HighScoreBoard) IsHighScore(score int) bool {
	if score <= 0 {
		return false
	}
	if len(b.Scores) < 10 {
		return true
	}
	return score > b.Scores[len(b.Scores)-1].Score
}
