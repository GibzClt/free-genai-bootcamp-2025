package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"lang-portal/backend_go/internal/models"

	"github.com/gin-gonic/gin"
)

func GetWords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		query := c.Query("q")
		perPage := 100
		offset := (page - 1) * perPage

		// Base query
		countQuery := "SELECT COUNT(*) FROM words"
		selectQuery := `
			SELECT id, japanese, romaji, english, parts
			FROM words
		`

		// Add search condition if query parameter exists
		var params []interface{}
		if query != "" {
			searchCond := `
				WHERE japanese LIKE ? 
				OR romaji LIKE ? 
				OR english LIKE ?
			`
			countQuery += " " + searchCond
			selectQuery += " " + searchCond
			searchPattern := "%" + query + "%"
			params = append(params, searchPattern, searchPattern, searchPattern)
		}

		// Add pagination
		selectQuery += " LIMIT ? OFFSET ?"
		params = append(params, perPage, offset)

		// Get total count
		var total int
		countErr := db.QueryRow(countQuery, params[:len(params)-2]...).Scan(&total)
		if countErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Count error: " + countErr.Error()})
			return
		}

		// Get paginated words
		rows, err := db.Query(selectQuery, params...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error: " + err.Error()})
			return
		}
		defer rows.Close()

		var words []models.Word
		for rows.Next() {
			var word models.Word
			err := rows.Scan(
				&word.ID,
				&word.Japanese,
				&word.Romaji,
				&word.English,
				&word.Parts,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
				return
			}
			words = append(words, word)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rows error: " + err.Error()})
			return
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
			"japanese": word.Japanese,
			"romaji":   word.Romaji,
			"english":  word.English,
			"stats": gin.H{
				"correct_count": stats.CorrectCount,
				"wrong_count":   stats.WrongCount,
			},
			"groups": groups,
		})
	}
}

func CreateWord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Japanese string `json:"japanese" binding:"required"`
			Romaji   string `json:"romaji" binding:"required"`
			English  string `json:"english" binding:"required"`
			Parts    string `json:"parts" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate Parts is valid JSON
		var partsJSON map[string]interface{}
		if err := json.Unmarshal([]byte(request.Parts), &partsJSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parts JSON"})
			return
		}

		result, err := db.Exec(`
			INSERT INTO words (japanese, romaji, english, parts)
			VALUES (?, ?, ?, ?)
		`, request.Japanese, request.Romaji, request.English, request.Parts)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()

		c.JSON(http.StatusCreated, gin.H{
			"id":      id,
			"success": true,
			"message": "Word created successfully",
		})
	}
}
