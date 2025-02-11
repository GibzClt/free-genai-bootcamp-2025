package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetQuickStats(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/dashboard/quick-stats", GetQuickStats(db))

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/dashboard/quick-stats", nil)
	r.ServeHTTP(w, req)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		SuccessRate        float64 `json:"success_rate"`
		TotalStudySessions int     `json:"total_study_sessions"`
		TotalActiveGroups  int     `json:"total_active_groups"`
		StudyStreakDays    int     `json:"study_streak_days"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
} 