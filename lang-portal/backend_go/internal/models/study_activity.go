package models

import (
	"database/sql"
	"time"
)

type StudyActivity struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ThumbnailURL string  `json:"thumbnail_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetStudyActivity(db *sql.DB, id int64) (*StudyActivity, error) {
	var activity StudyActivity
	err := db.QueryRow(`
		SELECT id, name, description, thumbnail_url, created_at
		FROM study_activities
		WHERE id = ?
	`, id).Scan(&activity.ID, &activity.Name, &activity.Description, &activity.ThumbnailURL, &activity.CreatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &activity, nil
} 