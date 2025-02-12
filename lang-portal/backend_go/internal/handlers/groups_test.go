package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/groups", GetGroups(db))

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
			wantItems:  1, // Basic Greetings group
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
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/groups?page=%s", tt.page), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				Items []struct {
					ID        int64  `json:"id"`
					Name      string `json:"name"`
					WordCount int    `json:"word_count"`
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
		})
	}
}

func TestGetGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/groups/:id", GetGroup(db))

	tests := []struct {
		name       string
		groupID    string
		wantStatus int
	}{
		{
			name:       "Valid group",
			groupID:    "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid group",
			groupID:    "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Invalid ID format",
			groupID:    "abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/groups/%s", tt.groupID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response struct {
					ID    int64  `json:"id"`
					Name  string `json:"name"`
					Stats struct {
						TotalWords    int     `json:"total_words"`
						StudiedWords  int     `json:"studied_words"`
						SuccessRate   float64 `json:"success_rate"`
						LastStudiedAt string  `json:"last_studied_at"`
					} `json:"stats"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.NotEmpty(t, response.Name)
			}
		})
	}
}

func TestGetGroupWords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/groups/:id/words", GetGroupWords(db))

	tests := []struct {
		name       string
		groupID    string
		page       string
		wantStatus int
		wantItems  int
	}{
		{
			name:       "Valid group first page",
			groupID:    "1",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  1, // One word in Basic Greetings
		},
		{
			name:       "Valid group empty page",
			groupID:    "1",
			page:       "999",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
		{
			name:       "Invalid group",
			groupID:    "999",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/groups/%s/words?page=%s", tt.groupID, tt.page)
			req, _ := http.NewRequest("GET", url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				Items []struct {
					ID           int64  `json:"id"`
					Japanese     string `json:"japanese"`
					Romaji       string `json:"romaji"`
					English      string `json:"english"`
					CorrectCount int    `json:"correct_count"`
					WrongCount   int    `json:"wrong_count"`
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
		})
	}
}

func TestGetGroupStudySessions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/groups/:id/study-sessions", GetGroupStudySessions(db))

	tests := []struct {
		name       string
		groupID    string
		page       string
		wantStatus int
		wantItems  int
	}{
		{
			name:       "Valid group with sessions",
			groupID:    "1",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  1, // Changed from 0 to 1 since we have one study session in test data
		},
		{
			name:       "Valid group empty page",
			groupID:    "1",
			page:       "999",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
		{
			name:       "Invalid group",
			groupID:    "999",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/groups/%s/study-sessions?page=%s", tt.groupID, tt.page)
			req, _ := http.NewRequest("GET", url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
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
			}
		})
	}
}
