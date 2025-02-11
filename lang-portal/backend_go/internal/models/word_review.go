package models

import (
	"database/sql"
	"time"
)

type WordReview struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

func CreateWordReview(db *sql.DB, review *WordReview) error {
	result, err := db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)
	`, review.WordID, review.StudySessionID, review.Correct, review.CreatedAt)
	
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	review.ID = id
	return nil
} 