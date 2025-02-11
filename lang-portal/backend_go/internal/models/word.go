package models

import (
	"database/sql"
	"encoding/json"
)

type Word struct {
	ID       int64             `json:"id"`
	Japanese string           `json:"japanese"`
	Romaji   string           `json:"romaji"`
	English  string           `json:"english"`
	Parts    json.RawMessage  `json:"parts"`
}

type WordStats struct {
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}

func GetWord(db *sql.DB, id int64) (*Word, error) {
	var word Word
	err := db.QueryRow(`
		SELECT id, japanese, romaji, english, parts
		FROM words
		WHERE id = ?
	`, id).Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts)
	
	if err != nil {
		return nil, err
	}
	
	return &word, nil
}

func GetWordStats(db *sql.DB, wordID int64) (*WordStats, error) {
	var stats WordStats
	err := db.QueryRow(`
		SELECT 
			COUNT(CASE WHEN correct = 1 THEN 1 END) as correct_count,
			COUNT(CASE WHEN correct = 0 THEN 1 END) as wrong_count
		FROM word_review_items
		WHERE word_id = ?
	`, wordID).Scan(&stats.CorrectCount, &stats.WrongCount)
	
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
} 