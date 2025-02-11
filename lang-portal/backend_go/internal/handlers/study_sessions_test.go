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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/study-sessions", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := map[string]bool{"correct": tt.correct}
			payloadBytes, _ := json.Marshal(payload)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				"POST",
				fmt.Sprintf("/api/study-sessions/%s/words/%s/review", tt.sessionID, tt.wordID),
				bytes.NewBuffer(payloadBytes),
			)
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
