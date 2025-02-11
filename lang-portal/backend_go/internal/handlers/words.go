package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
)

func GetWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage := 100
		offset := (page - 1) * perPage

		// Get total count
		var total int
		err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
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
			LEFT JOIN word_review_items wri ON wri.word_id = w.id
			GROUP BY w.id
			LIMIT ? OFFSET ?
		`, perPage, offset)
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

func GetWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
			return
		}

		word, err := models.GetWord(db, id)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		stats, err := models.GetWordStats(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get groups this word belongs to
		rows, err := db.Query(`
			SELECT g.id, g.name
			FROM groups g
			JOIN word_groups wg ON wg.group_id = g.id
			WHERE wg.word_id = ?
		`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var groups []models.Group
		for rows.Next() {
			var group models.Group
			if err := rows.Scan(&group.ID, &group.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			groups = append(groups, group)
		}

		c.JSON(http.StatusOK, gin.H{
			"word":   word,
			"stats":  stats,
			"groups": groups,
		})
	}
} 