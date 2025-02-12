package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetStudySessions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/study-sessions", GetStudySessions(db))

	tests := []struct {
		name       string
		page       string
		wantStatus int
		wantItems  int
	}{
		{
			name:       "First page",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  1, // One session from test data
		},
		{
			name:       "Empty page",
			page:       "999",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/study-sessions?page=%s", tt.page)
			req, _ := http.NewRequest("GET", url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				Items []struct {
					ID              int64  `json:"id"`
					ActivityName    string `json:"activity_name"`
					GroupName       string `json:"group_name"`
					StartTime       string `json:"start_time"`
					EndTime         string `json:"end_time"`
					ReviewItemCount int    `json:"review_items_count"`
				} `json:"items"`
				Pagination struct {
					CurrentPage  int `json:"current_page"`
					TotalPages   int `json:"total_pages"`
					TotalItems   int `json:"total_items"`
					ItemsPerPage int `json:"items_per_page"`
				} `json:"pagination"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Items, tt.wantItems)

			if tt.wantItems > 0 {
				assert.NotEmpty(t, response.Items[0].ActivityName)
				assert.NotEmpty(t, response.Items[0].GroupName)
				assert.NotEmpty(t, response.Items[0].StartTime)
			}
		})
	}
}

func TestGetStudySession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/study-sessions/:id", GetStudySession(db))

	tests := []struct {
		name       string
		sessionID  string
		wantStatus int
	}{
		{
			name:       "Valid session",
			sessionID:  "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid session",
			sessionID:  "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Invalid ID format",
			sessionID:  "abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/study-sessions/%s", tt.sessionID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response struct {
					ID              int64  `json:"id"`
					ActivityName    string `json:"activity_name"`
					GroupName       string `json:"group_name"`
					StartTime       string `json:"start_time"`
					EndTime         string `json:"end_time"`
					ReviewItemCount int    `json:"review_items_count"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.NotEmpty(t, response.ActivityName)
				assert.NotEmpty(t, response.GroupName)
			}
		})
	}
}

func TestGetStudySessionWords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/study-sessions/:id/words", GetStudySessionWords(db))

	tests := []struct {
		name       string
		sessionID  string
		wantStatus int
		wantItems  int
	}{
		{
			name:       "Valid session",
			sessionID:  "1",
			wantStatus: http.StatusOK,
			wantItems:  1, // One word review in test data
		},
		{
			name:       "Invalid session",
			sessionID:  "999",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/study-sessions/%s/words", tt.sessionID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				Items []struct {
					ID         int64  `json:"id"`
					Japanese   string `json:"japanese"`
					Romaji     string `json:"romaji"`
					English    string `json:"english"`
					Correct    bool   `json:"correct"`
					ReviewedAt string `json:"reviewed_at"`
				} `json:"items"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Items, tt.wantItems)

			if tt.wantItems > 0 {
				assert.NotEmpty(t, response.Items[0].Japanese)
				assert.NotEmpty(t, response.Items[0].English)
				assert.NotEmpty(t, response.Items[0].ReviewedAt)
			}
		})
	}
}

func TestCreateWordReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.POST("/api/study-sessions/:id/words/:word_id/review", CreateWordReview(db))

	tests := []struct {
		name       string
		sessionID  string
		wordID     string
		correct    bool
		wantStatus int
	}{
		{
			name:       "Valid review",
			sessionID:  "1",
			wordID:     "1",
			correct:    true,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Invalid session",
			sessionID:  "999",
			wordID:     "1",
			correct:    true,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Invalid word",
			sessionID:  "1",
			wordID:     "999",
			correct:    true,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Invalid session ID format",
			sessionID:  "abc",
			wordID:     "1",
			correct:    true,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]bool{"correct": tt.correct}
			payloadBytes, _ := json.Marshal(payload)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/study-sessions/%s/words/%s/review", tt.sessionID, tt.wordID)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response struct {
					Success        bool   `json:"success"`
					WordID         int64  `json:"word_id"`
					StudySessionID int64  `json:"study_session_id"`
					Correct        bool   `json:"correct"`
					CreatedAt      string `json:"created_at"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, tt.correct, response.Correct)
				assert.NotEmpty(t, response.CreatedAt)
			}
		})
	}
}
