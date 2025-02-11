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

func TestGetWords(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/words", GetWords(db))

	// Test cases
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
			wantItems:  3, // Assuming we have seeded 3 words
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
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/words?page=%s", tt.page), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response struct {
				Items []struct {
					ID    int64  `json:"id"`
					Words string `json:"words"`
				} `json:"items"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Len(t, response.Items, tt.wantItems)
		})
	}
} 