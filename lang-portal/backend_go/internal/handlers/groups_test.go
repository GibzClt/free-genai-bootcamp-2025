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
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Items, tt.wantItems)
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
		wantStatus int
		wantItems  int
	}{
		{
			name:       "Valid group",
			groupID:    "1",
			wantStatus: http.StatusOK,
			wantItems:  1, // One word in Basic Greetings
		},
		{
			name:       "Invalid group",
			groupID:    "999",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/groups/%s/words", tt.groupID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
