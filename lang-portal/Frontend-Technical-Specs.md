# Frontend Technical Specs

## Pages

### Dashboard "/dashboard"
#### Purpose
    The purpose of this page is to provide a quick overview of the user's study progress and to act as the default page when the user visits the portal.
#### Components
- Last Study Session
    - shows last activity used
    - shows when last activity used
    - summarizes wrong vs correct from last activity
    - has a link to the group
- Study Progress
    - total words study eg 5/110
        - total words studied across all study sessions of all possible words in our database
    - display a mastery progress eg 0%
- Quick Stats
    - success rate eg 70%
    - total study sessions eg 10
    - total active groups eg 3
    - study streak eg 3 days
- Start Studying Button
    - goes to study activities page
#### Required API Endpoints
    GET /api/dashboard/last-study-session
    GET /api/dashboard/study-progress
    GET /api/dashboard/quick-stats

### Study Activities Index Page "/study-activities"
#### Purpose
    The purpose of this page is to provide a list of study activities with a thumbnail and its name to either launch or view the study activity.
#### Components
    - Study Activity Card
        - thumbnail of the study activity
        - name of the study activity
        - launch button to launch the study activity
        - view button to view more information about past study

#### Required API Endpoints
    - GET /api/study-activities
    
### Study Activity "/study_activity/:id"
#### Purpose
    The purpose of this page is to provide a detailed view of a specific study activity.
#### Components
    - thumbnail of the study activity
    - name of the study activity
    - launch button to launch the study activity
    - Study Activities Paginated List
        - id
        - activity name
        - group name
        - start time
        - end time (inferred from last word_review_item submitted)
        - number of review items
#### Required API Endpoints
    - GET /api/study-activity/:id
    - GET /api/study-activity/:id/study-sessions

### Study Activity Launcher "/study-activity/:id/launch"
#### Purpose
    The purpose of this page is to provide a launcher for a specific study activity.
#### Components
    - Name of study activity
    - Launch form
        - select field for group
        - launch now button
#### Required API Endpoints
    - POST /api/study-activities

## Behavior
After the form is submitted a new tab opens with the study activity based on the url provided in the database.
Also after the form is submitted the page will redirect to the study session show page

### Words Index Page "/words"
#### Purpose
    The purpose of this page is to provide a list of all words in our database.
#### Components
    - Paginated Word List
        - Columns
            - Japanese
            - Romaji
            - English
            - Correct count
            - Wrong count
        - Pagination with 100 items per page
        - Clicking the Japanese word will take us to the word show page
#### Required API Endpoints
    - GET /api/words

### Word Show Page "/words/:id"
#### Purpose
    The purpose of this page is to provide a detailed information of a specific word.
#### Components
    - Word Card
        - Japanese
        - Romaji
        - English
        - Study statistics
            - Correct count
            - Wrong count
        - Word Groups
            - show on a series of pills eg.tags
            - when group name is clicked, it will take us to the group show page
#### Required API Endpoints
    - GET /api/words/:id

### Word Groups Index Page "/groups"
#### Purpose
    The purpose of this page is to provide a list of groups in our database.
#### Components
    - Paginated Group List
        - Columns
            - Group Name
            - Word Count
            - Group Words Count
        - Clicking the Group Name will take us to the group show page
#### Required API Endpoints
    - GET /api/groups


### Group Show Page "/groups/:id"
#### Purpose
    The purpose of this page is to provide a detailed information of a specific group.
#### Components
    - Group Name
    - Group Statistics
        - Total word count
    - Words in Group (Paginated list of words)
        - Should use the same component as the words index page
    - Study Sessions (Paginated list of study sessions)
        - Should use the same component as the study sessions index page
#### Required API Endpoints
    - GET /api/groups/:id (the name and group statistics)
    - GET /api/groups/:id/words
    - GET /api/groups/:id/study-sessions

### Study Sessions Index Page "/study-sessions"
#### Purpose
    The purpose of this page is to provide a list of all study sessions in our database.
#### Components
    - Paginated Study Session List
        - Columns
            - id
            - activity name
            - group name
            - start time
            - end time
            - Number of review items
        - Pagination with 100 items per page
        - Clicking the study session id will take us to the study session show page
#### Required API Endpoints
    - GET /api/study-sessions

### Study Session Show Page "/study-sessions/:id"
#### Purpose
    The purpose of this page is to provide a detailed information of a specific study session.
#### Components
    - Study Session Details
        - Activity Name
        - Group Name
        - Start Time
        - End Time
        - Number of review items
    - Words Reviewed (Paginated list of words reviewed)
        - Should use the same component as the words index page
#### Required API Endpoints
    - GET /api/study-sessions/:id
    - GET /api/study-sessions/:id/words

### Settings Page "/settings"
#### Purpose
    The purpose of this page is to provide a page to update the configurations of the portal.
#### Components
    - Theme selection
        - light
        - dark
        - system default
    - Reset history button
        - will reset all study sessions and words reviewed
    - Full Reset Data
        - will drop all tables and re-create with seed data
#### Required API Endpoints
    - POST /api/settings/reset-history
    - POST /api/settings/full-reset
