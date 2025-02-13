package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"lang-portal/backend_go/internal/handlers"
)

func main() {
	// Connect to database
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Dashboard endpoints
		api.GET("/dashboard/last-study-session", handlers.GetLastStudySession(db))
		api.GET("/dashboard/study-progress", handlers.GetStudyProgress(db))
		api.GET("/dashboard/quick-stats", handlers.GetQuickStats(db))

		// Study activities endpoints
		api.GET("/study-activities", handlers.GetStudyActivities(db))
		api.GET("/study-activity/:id", handlers.GetStudyActivity(db))
		api.GET("/study-activity/:id/study-sessions", handlers.GetStudyActivitySessions(db))
		api.POST("/study-activities", handlers.CreateStudyActivity(db))

		// Words endpoints
		api.GET("/words", handlers.GetWords(db))
		api.GET("/words/:id", handlers.GetWord(db))

		// Groups endpoints
		api.GET("/groups", handlers.GetGroups(db))
		api.GET("/groups/:id", handlers.GetGroup(db))
		api.GET("/groups/:id/words", handlers.GetGroupWords(db))
		api.GET("/groups/:id/study-sessions", handlers.GetGroupStudySessions(db))

		// Study sessions endpoints
		api.GET("/study-sessions", handlers.GetStudySessions(db))
		api.GET("/study-sessions/:id", handlers.GetStudySession(db))
		api.GET("/study-sessions/:id/words", handlers.GetStudySessionWords(db))
		api.POST("/study-sessions/:id/words/:word_id/review", handlers.CreateWordReview(db))

		// Settings endpoints
		api.POST("/settings/reset-history", handlers.ResetHistory(db))
		api.POST("/settings/full-reset", handlers.FullReset(db))
	}

	// Start server
	if err := r.Run(":4000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}