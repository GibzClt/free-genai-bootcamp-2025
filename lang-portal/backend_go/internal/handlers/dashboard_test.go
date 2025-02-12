package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLastStudySession(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/dashboard/last-study-session", GetLastStudySession(db))

	// Test cases
	tests := []struct {
		name       string
		wantStatus int
		wantEmpty  bool
	}{
		{
			name:       "Has study session",
			wantStatus: http.StatusOK,
			wantEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/dashboard/last-study-session", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				ID              int64  `json:"id"`
				GroupID         int64  `json:"group_id"`
				GroupName       string `json:"group_name"`
				StudyActivityID int64  `json:"study_activity_id"`
				CreatedAt       string `json:"created_at"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if !tt.wantEmpty {
				assert.NotZero(t, response.ID)
				assert.NotZero(t, response.GroupID)
				assert.NotEmpty(t, response.GroupName)
				assert.NotZero(t, response.StudyActivityID)
				assert.NotEmpty(t, response.CreatedAt)
			}
		})
	}
}

func TestGetStudyProgress(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/dashboard/study-progress", GetStudyProgress(db))

	// Test case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/dashboard/study-progress", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		TotalWordsStudied   int `json:"total_words_studied"`
		TotalAvailableWords int `json:"total_available_words"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// We should have 3 total words from test data
	assert.Equal(t, 3, response.TotalAvailableWords)
	// And 1 studied word from the test data
	assert.Equal(t, 1, response.TotalWordsStudied)
}

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

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		SuccessRate        float64 `json:"success_rate"`
		TotalStudySessions int     `json:"total_study_sessions"`
		TotalActiveGroups  int     `json:"total_active_groups"`
		StudyStreakDays    int     `json:"study_streak_days"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// From our test data:
	// - We have 1 correct review out of 1 total = 100% success rate
	assert.Equal(t, float64(100), response.SuccessRate)
	// - We have 1 study session
	assert.Equal(t, 1, response.TotalStudySessions)
	// - We have 1 active group
	assert.Equal(t, 1, response.TotalActiveGroups)
	// - We have 1 day streak (today)
	assert.Equal(t, 1, response.StudyStreakDays)
}
