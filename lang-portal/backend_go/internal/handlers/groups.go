package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"lang-portal/backend_go/internal/models"

	"github.com/gin-gonic/gin"
)

func GetGroups(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 100
		offset := (page - 1) * perPage

		// Get total count
		var total int
		err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get paginated groups with word counts
		rows, err := db.Query(`
			SELECT 
				g.id,
				g.name,
				COUNT(DISTINCT wg.word_id) as word_count
			FROM groups g
			LEFT JOIN word_groups wg ON wg.group_id = g.id
			GROUP BY g.id
			LIMIT ? OFFSET ?
		`, perPage, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var groups []struct {
			models.Group
			WordCount int `json:"word_count"`
		}

		for rows.Next() {
			var group struct {
				models.Group
				WordCount int `json:"word_count"`
			}
			err := rows.Scan(&group.ID, &group.Name, &group.WordCount)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			groups = append(groups, group)
		}

		c.JSON(http.StatusOK, gin.H{
			"items": groups,
			"pagination": gin.H{
				"current_page":   page,
				"total_pages":    (total + perPage - 1) / perPage,
				"total_items":    total,
				"items_per_page": perPage,
			},
		})
	}
}

func GetGroup(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
			return
		}

		group, err := models.GetGroup(db, id)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		stats, err := models.GetGroupStats(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    group.ID,
			"name":  group.Name,
			"stats": stats,
		})
	}
}

func GetGroupWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
			return
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 100
		offset := (page - 1) * perPage

		// Get total count
		var total int
		err = db.QueryRow(`
			SELECT COUNT(*) 
			FROM word_groups 
			WHERE group_id = ?
		`, groupID).Scan(&total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get paginated words with stats
		rows, err := db.Query(`
			SELECT 
				w.id,
				w.japanese,
				w.romaji,
				w.english,
				COUNT(CASE WHEN wri.correct = 1 THEN 1 END) as correct_count,
				COUNT(CASE WHEN wri.correct = 0 THEN 1 END) as wrong_count
			FROM words w
			JOIN word_groups wg ON wg.word_id = w.id
			LEFT JOIN word_review_items wri ON wri.word_id = w.id
			WHERE wg.group_id = ?
			GROUP BY w.id
			LIMIT ? OFFSET ?
		`, groupID, perPage, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var words []struct {
			models.Word
			CorrectCount int `json:"correct_count"`
			WrongCount   int `json:"wrong_count"`
		}

		for rows.Next() {
			var word struct {
				models.Word
				CorrectCount int `json:"correct_count"`
				WrongCount   int `json:"wrong_count"`
			}
			err := rows.Scan(
				&word.ID,
				&word.Japanese,
				&word.Romaji,
				&word.English,
				&word.CorrectCount,
				&word.WrongCount,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			words = append(words, word)
		}

		c.JSON(http.StatusOK, gin.H{
			"items": words,
			"pagination": gin.H{
				"current_page":   page,
				"total_pages":    (total + perPage - 1) / perPage,
				"total_items":    total,
				"items_per_page": perPage,
			},
		})
	}
}

func GetGroupStudySessions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
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
			WHERE group_id = ?
		`, groupID).Scan(&total)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Count error: " + err.Error()})
			return
		}

		// Get paginated study sessions with simpler query first
		rows, err := db.Query(`
			SELECT 
				ss.id,
				sa.name,
				g.name,
				ss.created_at,
				ss.created_at,
				0
			FROM study_sessions ss
			JOIN study_activities sa ON sa.id = ss.study_activity_id
			JOIN groups g ON g.id = ss.group_id
			WHERE ss.group_id = ?
			ORDER BY ss.created_at DESC
			LIMIT ? OFFSET ?
		`, groupID, perPage, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error: " + err.Error()})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
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

// ... continuing with more handlers ...
