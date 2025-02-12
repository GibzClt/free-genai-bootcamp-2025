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

func TestGetStudyActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/study-activity/:id", GetStudyActivity(db))

	tests := []struct {
		name       string
		activityID string
		wantStatus int
	}{
		{
			name:       "Valid activity",
			activityID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid activity",
			activityID: "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Invalid ID format",
			activityID: "abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/study-activity/%s", tt.activityID), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response struct {
					ID           int64  `json:"id"`
					Name         string `json:"name"`
					Description  string `json:"description"`
					ThumbnailURL string `json:"thumbnail_url"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.NotEmpty(t, response.Name)
				assert.NotEmpty(t, response.Description)
			}
		})
	}
}

func TestGetStudyActivitySessions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.GET("/api/study-activity/:id/study-sessions", GetStudyActivitySessions(db))

	tests := []struct {
		name       string
		activityID string
		page       string
		wantStatus int
		wantItems  int
	}{
		{
			name:       "Valid activity first page",
			activityID: "1",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  1, // One session in test data
		},
		{
			name:       "Valid activity empty page",
			activityID: "1",
			page:       "999",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
		{
			name:       "Invalid activity",
			activityID: "999",
			page:       "1",
			wantStatus: http.StatusOK,
			wantItems:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/study-activity/%s/study-sessions?page=%s", tt.activityID, tt.page)
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

func TestCreateStudyActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	db := setupTestDB(t)
	defer db.Close()

	r.POST("/api/study-activities", CreateStudyActivity(db))

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid activity",
			payload: map[string]interface{}{
				"group_id":          float64(1),
				"study_activity_id": float64(1),
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing group_id",
			payload: map[string]interface{}{
				"study_activity_id": float64(1),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing study_activity_id",
			payload: map[string]interface{}{
				"group_id": float64(1),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid group",
			payload: map[string]interface{}{
				"group_id":          float64(999),
				"study_activity_id": float64(1),
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid activity",
			payload: map[string]interface{}{
				"group_id":          float64(1),
				"study_activity_id": float64(999),
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/study-activities", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var response struct {
					ID      int64  `json:"id"`
					GroupID int64  `json:"group_id"`
					Success bool   `json:"success"`
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.NotZero(t, response.ID)
				assert.Equal(t, int64(tt.payload["group_id"].(float64)), response.GroupID)
			}
		})
	}
}
