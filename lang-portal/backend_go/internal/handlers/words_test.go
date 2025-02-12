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

func TestGetWords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/words", GetWords(db))

	tests := []struct {
		name       string
		page       string
		query      string
		wantStatus int
		wantItems  int
	}{
		{
			name:       "First page",
			page:       "1",
			query:      "",
			wantStatus: http.StatusOK,
			wantItems:  3, // Three words from test data
		},
		{
			name:       "Empty page",
			page:       "999",
			query:      "",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
		{
			name:       "Search by Japanese",
			page:       "1",
			query:      "こんにちは",
			wantStatus: http.StatusOK,
			wantItems:  1,
		},
		{
			name:       "Search by Romaji",
			page:       "1",
			query:      "konnichiwa",
			wantStatus: http.StatusOK,
			wantItems:  1,
		},
		{
			name:       "Search by English",
			page:       "1",
			query:      "hello",
			wantStatus: http.StatusOK,
			wantItems:  1,
		},
		{
			name:       "No results",
			page:       "1",
			query:      "xyz",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/words?page=%s", tt.page)
			if tt.query != "" {
				url += fmt.Sprintf("&q=%s", tt.query)
			}
			req, _ := http.NewRequest("GET", url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				Items []struct {
					ID       int64  `json:"id"`
					Japanese string `json:"japanese"`
					Romaji   string `json:"romaji"`
					English  string `json:"english"`
					Parts    string `json:"parts"`
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
				assert.NotEmpty(t, response.Items[0].Japanese)
				assert.NotEmpty(t, response.Items[0].Romaji)
				assert.NotEmpty(t, response.Items[0].English)
			}
		})
	}
}

func TestGetWord(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/words/:id", GetWord(db))

	tests := []struct {
		name       string
		wordID     string
		wantStatus int
	}{
		{
			name:       "Valid word",
			wordID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid word",
			wordID:     "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Invalid ID format",
			wordID:     "abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/words/%s", tt.wordID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response struct {
					Japanese string `json:"japanese"`
					Romaji   string `json:"romaji"`
					English  string `json:"english"`
					Stats    struct {
						CorrectCount int `json:"correct_count"`
						WrongCount   int `json:"wrong_count"`
					} `json:"stats"`
					Groups []struct {
						ID   int64  `json:"id"`
						Name string `json:"name"`
					} `json:"groups"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.Japanese)
				assert.NotEmpty(t, response.Romaji)
				assert.NotEmpty(t, response.English)
				// Stats and Groups may be empty but should exist
				assert.NotNil(t, response.Stats)
				assert.NotNil(t, response.Groups)
			}
		})
	}
}

func TestCreateWord(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.POST("/api/words", CreateWord(db))

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid word",
			payload: map[string]interface{}{
				"japanese": "おはよう",
				"romaji":   "ohayou",
				"english":  "good morning",
				"parts":    `{"type":"greeting"}`,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing Japanese",
			payload: map[string]interface{}{
				"romaji":  "ohayou",
				"english": "good morning",
				"parts":   `{"type":"greeting"}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Romaji",
			payload: map[string]interface{}{
				"japanese": "おはよう",
				"english":  "good morning",
				"parts":    `{"type":"greeting"}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing English",
			payload: map[string]interface{}{
				"japanese": "おはよう",
				"romaji":   "ohayou",
				"parts":    `{"type":"greeting"}`,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Parts JSON",
			payload: map[string]interface{}{
				"japanese": "おはよう",
				"romaji":   "ohayou",
				"english":  "good morning",
				"parts":    "{invalid json}",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/words", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response struct {
					ID      int64  `json:"id"`
					Success bool   `json:"success"`
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.NotZero(t, response.ID)
			}
		})
	}
}
