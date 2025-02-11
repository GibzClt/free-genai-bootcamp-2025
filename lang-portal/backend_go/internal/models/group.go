package models

import (
	"database/sql"
)

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GroupStats struct {
	TotalWordCount int `json:"total_word_count"`
}

func GetGroup(db *sql.DB, id int64) (*Group, error) {
	var group Group
	err := db.QueryRow(`
		SELECT id, name
		FROM groups
		WHERE id = ?
	`, id).Scan(&group.ID, &group.Name)
	
	if err != nil {
		return nil, err
	}
	
	return &group, nil
}

func GetGroupStats(db *sql.DB, groupID int64) (*GroupStats, error) {
	var stats GroupStats
	err := db.QueryRow(`
		SELECT COUNT(DISTINCT word_id)
		FROM word_groups
		WHERE group_id = ?
	`, groupID).Scan(&stats.TotalWordCount)
	
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
} 