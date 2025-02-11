package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLastStudySession(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var session struct {
			ID              int64  `json:"id"`
			GroupID         int64  `json:"group_id"`
			GroupName       string `json:"group_name"`
			StudyActivityID int64  `json:"study_activity_id"`
			CreatedAt       string `json:"created_at"`
		}

		err := db.QueryRow(`
			SELECT 
				ss.id,
				ss.group_id,
				g.name,
				ss.study_activity_id,
				ss.created_at
			FROM study_sessions ss
			JOIN groups g ON g.id = ss.group_id
			ORDER BY ss.created_at DESC
			LIMIT 1
		`).Scan(
			&session.ID,
			&session.GroupID,
			&session.GroupName,
			&session.StudyActivityID,
			&session.CreatedAt,
		)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No study sessions found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func GetStudyProgress(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var progress struct {
			TotalWordsStudied    int `json:"total_words_studied"`
			TotalAvailableWords int `json:"total_available_words"`
		}

		err := db.QueryRow(`
			SELECT 
				(SELECT COUNT(DISTINCT word_id) FROM word_review_items) as studied,
				(SELECT COUNT(*) FROM words) as total
		`).Scan(&progress.TotalWordsStudied, &progress.TotalAvailableWords)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, progress)
	}
}

func GetQuickStats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stats struct {
			SuccessRate       float64 `json:"success_rate"`
			TotalStudySessions int    `json:"total_study_sessions"`
			TotalActiveGroups int    `json:"total_active_groups"`
			StudyStreakDays   int    `json:"study_streak_days"`
		}

		// Get success rate and total study sessions
		err := db.QueryRow(`
			WITH review_stats AS (
				SELECT 
					CAST(SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) AS FLOAT) as correct_count,
					COUNT(*) as total_count
				FROM word_review_items
			)
			SELECT 
				CASE 
					WHEN total_count > 0 THEN (correct_count / total_count) * 100 
					ELSE 0 
				END,
				(SELECT COUNT(DISTINCT id) FROM study_sessions)
			FROM review_stats
		`).Scan(&stats.SuccessRate, &stats.TotalStudySessions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get total active groups
		err = db.QueryRow(`
			SELECT COUNT(DISTINCT group_id) 
			FROM study_sessions
		`).Scan(&stats.TotalActiveGroups)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Calculate study streak
		err = db.QueryRow(`
			WITH RECURSIVE dates AS (
				SELECT date(MAX(created_at)) as date
				FROM study_sessions
				UNION ALL
				SELECT date(date, '-1 day')
				FROM dates
				WHERE EXISTS (
					SELECT 1 
					FROM study_sessions 
					WHERE date(created_at) = date(date, '-1 day')
				)
			)
			SELECT COUNT(*) FROM dates
		`).Scan(&stats.StudyStreakDays)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
} 