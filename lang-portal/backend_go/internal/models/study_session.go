package models

import (
	"database/sql"
	"time"
)

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudySessionDetail struct {
	ID              int64  `json:"id"`
	ActivityName    string `json:"activity_name"`
	GroupName       string `json:"group_name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	ReviewItemCount int    `json:"review_items_count"`
}

func GetStudySession(db *sql.DB, id int64) (*StudySessionDetail, error) {
	var session StudySessionDetail
	err := db.QueryRow(`
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at as start_time,
			MAX(wri.created_at) as end_time,
			COUNT(wri.id) as review_items_count
		FROM study_sessions ss
		JOIN study_activities sa ON sa.id = ss.study_activity_id
		JOIN groups g ON g.id = ss.group_id
		LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
		WHERE ss.id = ?
		GROUP BY ss.id
	`, id).Scan(
		&session.ID,
		&session.ActivityName,
		&session.GroupName,
		&session.StartTime,
		&session.EndTime,
		&session.ReviewItemCount,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &session, nil
} 