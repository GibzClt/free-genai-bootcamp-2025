# API Documentation

## Base URL

All API endpoints are prefixed with `/api`.

## Authentication

Currently, the API does not require authentication.

## Endpoints

### Dashboard

#### GET /api/dashboard/last-study-session
Returns information about the most recent study session.

**Response**
```json
{
  "id": 1,
  "group_id": 1,
  "group_name": "Basic Greetings",
  "study_activity_id": 1,
  "created_at": "2024-03-10T15:04:05Z"
}
```

#### GET /api/dashboard/study-progress
Returns overall study progress statistics.

**Response**
```json
{
  "total_words_studied": 50,
  "total_available_words": 100
}
```

#### GET /api/dashboard/quick-stats
Returns quick statistics about the user's study progress.

**Response**
```json
{
  "success_rate": 85.5,
  "total_study_sessions": 10,
  "total_active_groups": 3,
  "study_streak_days": 5
}
```

### Words

#### GET /api/words
Returns a paginated list of words with their study statistics.

**Query Parameters**
- `page`: Page number (default: 1)

**Response**
```json
{
  "items": [
    {
      "id": 1,
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "parts": {"type": "greeting"},
      "correct_count": 5,
      "wrong_count": 1
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 10,
    "total_items": 100,
    "items_per_page": 10
  }
}
```

### Groups

#### GET /api/groups
Returns a paginated list of word groups.

#### GET /api/groups/:id
Returns details about a specific group.

#### GET /api/groups/:id/words
Returns words belonging to a specific group.

#### GET /api/groups/:id/study-sessions
Returns study sessions for a specific group.

### Study Activities

#### GET /api/study-activity/:id
Returns details about a specific study activity.

#### GET /api/study-activity/:id/study-sessions
Returns study sessions for a specific activity.

#### POST /api/study-activities
Creates a new study activity session.

**Request Body**
```json
{
  "group_id": 1,
  "study_activity_id": 1
}
```

### Study Sessions

#### GET /api/study-sessions
Returns a paginated list of study sessions.

#### GET /api/study-sessions/:id
Returns details about a specific study session.

#### GET /api/study-sessions/:id/words
Returns words reviewed in a specific study session.

#### POST /api/study-sessions/:id/words/:word_id/review
Records a word review result.

**Request Body**
```json
{
  "correct": true
}
```

### Settings

#### POST /api/settings/reset-history
Resets all study history while preserving words and groups.

#### POST /api/settings/full-reset
Performs a complete system reset, removing all data.

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error 