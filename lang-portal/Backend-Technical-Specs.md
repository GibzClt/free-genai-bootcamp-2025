# Backend Server Technical Specs

## Business Goal: 
A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps
 
 ## Technical Requirements
 - The backend server will be written in Go
 - The database will be SQLite3
 - The API will be built using Gin framework
 - The API will always return JSON
 - There will be no authentication or authorization
 - Everything will be treated as a single user

 ## Database Schema
 We have the following tables:
 - words - stored vocabulary words
    - id integer
    - japanese string
    - romaji string
    - english string
    - parts json
 - words_groups - join table for words and groups many to many relationship
    - id integer
    - word_id integer
    - group_id integer
 - groups - thematic groups of words
    - id integer
    - name string
 - study_sessions - records of study sessions grouping word_review_items
    - id integer
    - word_id integer
    - group_id integer
    - created_at datetime
    - study_activity_id integer
 - study_activities - a specific study activity linking a study session to group
    - id integer
    - study_session_id integer
    - group_id integer
    - created_at datetime
 - word_review_items - a record of word practice, determining if the word is correct or not
    - id integer
    - word_id integer
    - study_session_id integer
    - correct boolean
    - created_at datetime

## API Endpoints
- GET /api/dashboard/last-study-session
- GET /api/dashboard/study-progress
- GET /api/dashboard/quick-stats
- GET /api/study-activity/:id
- GET /api/study-activity/:id/study-sessions
- POST /api/study-activities
   - required params
      - group_id
      - study_activity_id
- GET /api/words
   - paginated with 100 items per page
- GET /api/words/:id
- GET /api/groups
   - paginated with 100 items per page
 - GET /api/groups/:id
- GET /api/groups/:id/words
   - paginated with 100 items per page
- GET /api/groups/:id/study-sessions
   - paginated with 100 items per page
- GET /api/study-sessions
   - paginated with 100 items per page
- GET /api/study-sessions/:id
- GET /api/study-sessions/:id/words
- GET /api/study-sessions/:id/words/:word_id/review
   - required params
      - correct boolean
- POST /api/settings/reset-history
- POST /api/settings/full-reset
