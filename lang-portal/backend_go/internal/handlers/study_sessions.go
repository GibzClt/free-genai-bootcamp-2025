package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
)

func GetStudySessions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 100
		offset := (page - 1) * perPage

		// Get total count
		var total int
		err := db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
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
			GROUP BY ss.id
			ORDER BY ss.created_at DESC
			LIMIT ? OFFSET ?
		`, perPage, offset)
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

func GetStudySession(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		session, err := models.GetStudySession(db, id)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func GetStudySessionWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		rows, err := db.Query(`
			SELECT 
				w.id,
				w.japanese,
				w.romaji,
				w.english,
				wri.correct,
				wri.created_at as reviewed_at
			FROM words w
			JOIN word_review_items wri ON wri.word_id = w.id
			WHERE wri.study_session_id = ?
			ORDER BY wri.created_at
		`, sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var words []struct {
			models.Word
			Correct    bool      `json:"correct"`
			ReviewedAt time.Time `json:"reviewed_at"`
		}

		for rows.Next() {
			var word struct {
				models.Word
				Correct    bool      `json:"correct"`
				ReviewedAt time.Time `json:"reviewed_at"`
			}
			err := rows.Scan(
				&word.ID,
				&word.Japanese,
				&word.Romaji,
				&word.English,
				&word.Correct,
				&word.ReviewedAt,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			words = append(words, word)
		}

		c.JSON(http.StatusOK, gin.H{"items": words})
	}
}

func CreateWordReview(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		// Verify session exists
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM study_sessions WHERE id = ?)", sessionID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Study session not found"})
			return
		}

		wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
			return
		}

		var request struct {
			Correct bool `json:"correct" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		review := &models.WordReview{
			WordID:         wordID,
			StudySessionID: sessionID,
			Correct:        request.Correct,
			CreatedAt:      time.Now(),
		}

		if err := models.CreateWordReview(db, review); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":          true,
			"word_id":         wordID,
			"study_session_id": sessionID,
			"correct":         request.Correct,
			"created_at":      review.CreatedAt,
		})
	}
} 