package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"lang-portal/backend_go/internal/models"

	"github.com/gin-gonic/gin"
)

func GetStudyActivity(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
			return
		}

		var activity struct {
			ID           int64  `json:"id"`
			Name         string `json:"name"`
			Description  string `json:"description"`
			ThumbnailURL string `json:"thumbnail_url"`
		}

		err = db.QueryRow(`
			SELECT id, name, description, thumbnail_url
			FROM study_activities
			WHERE id = ?
		`, id).Scan(&activity.ID, &activity.Name, &activity.Description, &activity.ThumbnailURL)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, activity)
	}
}

func GetStudyActivitySessions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 100
		offset := (page - 1) * perPage

		// Get total count
		var total int
		err = db.QueryRow(`
			SELECT COUNT(*) 
			FROM study_sessions 
			WHERE study_activity_id = ?
		`, activityID).Scan(&total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get paginated study sessions
		rows, err := db.Query(`
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
			WHERE ss.study_activity_id = ?
			GROUP BY ss.id
			ORDER BY ss.created_at DESC
			LIMIT ? OFFSET ?
		`, activityID, perPage, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var sessions []models.StudySessionDetail
		for rows.Next() {
			var session models.StudySessionDetail
			err := rows.Scan(
				&session.ID,
				&session.ActivityName,
				&session.GroupName,
				&session.StartTime,
				&session.EndTime,
				&session.ReviewItemCount,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sessions = append(sessions, session)
		}

		c.JSON(http.StatusOK, gin.H{
			"items": sessions,
			"pagination": gin.H{
				"current_page":   page,
				"total_pages":    (total + perPage - 1) / perPage,
				"total_items":    total,
				"items_per_page": perPage,
			},
		})
	}
}

func CreateStudyActivity(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			GroupID         int64 `json:"group_id" binding:"required"`
			StudyActivityID int64 `json:"study_activity_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec(`
			INSERT INTO study_sessions (group_id, study_activity_id)
			VALUES (?, ?)
		`, request.GroupID, request.StudyActivityID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sessionID, _ := result.LastInsertId()

		c.JSON(http.StatusCreated, gin.H{
			"id":       sessionID,
			"group_id": request.GroupID,
			"success":  true,
			"message": "Study activity created successfully",
		})
	}
}
