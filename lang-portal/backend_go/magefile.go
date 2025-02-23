//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "words.db"

// InitDB initializes the SQLite database
func InitDB() error {
	fmt.Println("Initializing database...")
	
	if _, err := os.Stat(dbName); err == nil {
		fmt.Printf("Database %s already exists\n", dbName)
		return nil
	}

	file, err := os.Create(dbName)
	if err != nil {
		return fmt.Errorf("error creating database file: %v", err)
	}
	file.Close()

	fmt.Printf("Created database %s\n", dbName)
	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	fmt.Println("Running migrations...")

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}

	// Get list of migration files
	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("error finding migration files: %v", err)
	}
	sort.Strings(files)

	// Run each migration
	for _, file := range files {
		name := filepath.Base(file)
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = ?)", name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error checking migration status: %v", err)
		}

		if exists {
			fmt.Printf("Skipping %s (already applied)\n", name)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %v", err)
		}

		// Split the content into individual statements
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = tx.Exec(stmt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error executing migration %s: %v", name, err)
			}
		}

		_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %v", name, err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing migration %s: %v", name, err)
		}

		fmt.Printf("Applied %s\n", name)
	}

	return nil
}

// SeedWord represents a word in our seed files
type SeedWord struct {
	Japanese string          `json:"japanese"`
	Romaji   string         `json:"romaji"`
	English  string         `json:"english"`
	Parts    map[string]any `json:"parts"`
}

// SeedFile represents the structure of our seed files
type StudyActivity struct {
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

type ConfigFile struct {
	Groups          []struct {
		Name       string `json:"name"`
		SourceFile string `json:"source_file"`
	} `json:"groups"`
	StudyActivities []StudyActivity `json:"study_activities"`
}

type SeedFile struct {
	GroupName string     `json:"group_name"`
	Words     []SeedWord `json:"words"`
}

// Seed loads initial data into the database
func Seed() error {
	fmt.Println("Seeding database...")

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// First, process config.json for study activities
	configContent, err := os.ReadFile("db/seeds/config.json")
	if err != nil {
		return fmt.Errorf("error reading config.json: %v", err)
	}

	var config ConfigFile
	if err := json.Unmarshal(configContent, &config); err != nil {
		return fmt.Errorf("error parsing config.json: %v", err)
	}

	// Begin transaction for study activities
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Clear existing study activities
	_, err = tx.Exec("DELETE FROM study_activities")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error clearing study activities: %v", err)
	}

	// Insert study activities
	for _, activity := range config.StudyActivities {
		_, err = tx.Exec(`
			INSERT INTO study_activities (name, description, thumbnail_url)
			VALUES (?, ?, ?)
		`, activity.Name, activity.Description, activity.ThumbnailURL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting study activity: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing study activities: %v", err)
	}

	// Add random study sessions and word review items
	tx, err = db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction for study sessions: %v", err)
	}

	// Get all group IDs
	rows, err := tx.Query("SELECT id FROM groups")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting group IDs: %v", err)
	}
	var groupIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			tx.Rollback()
			return fmt.Errorf("error scanning group ID: %v", err)
		}
		groupIDs = append(groupIDs, id)
	}
	rows.Close()

	// Get all activity IDs
	rows, err = tx.Query("SELECT id FROM study_activities")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting activity IDs: %v", err)
	}
	var activityIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			tx.Rollback()
			return fmt.Errorf("error scanning activity ID: %v", err)
		}
		activityIDs = append(activityIDs, id)
	}
	rows.Close()

	// Get all word IDs
	rows, err = tx.Query("SELECT id FROM words")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting word IDs: %v", err)
	}
	var wordIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			tx.Rollback()
			return fmt.Errorf("error scanning word ID: %v", err)
		}
		wordIDs = append(wordIDs, id)
	}
	rows.Close()

	// Clear existing study sessions and word review items
	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error clearing word review items: %v", err)
	}
	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error clearing study sessions: %v", err)
	}

	// Create random study sessions over the last 30 days
	for i := 0; i < 50; i++ { // Create 50 study sessions
		// Random group and activity
		groupID := groupIDs[rand.Intn(len(groupIDs))]
		activityID := activityIDs[rand.Intn(len(activityIDs))]

		// Random date within last 30 days
		daysAgo := rand.Intn(30)
		createdAt := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02 15:04:05")

		// Insert study session
		var sessionID int64
		err = tx.QueryRow(`
			INSERT INTO study_sessions (group_id, study_activity_id, created_at)
			VALUES (?, ?, ?)
			RETURNING id
		`, groupID, activityID, createdAt).Scan(&sessionID)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting study session: %v", err)
		}

		// Create 5-15 word review items for each session
		numReviews := rand.Intn(11) + 5
		for j := 0; j < numReviews; j++ {
			wordID := wordIDs[rand.Intn(len(wordIDs))]
			correct := rand.Float32() < 0.7 // 70% chance of correct answer

			_, err = tx.Exec(`
				INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
				VALUES (?, ?, ?, ?)
			`, wordID, sessionID, correct, createdAt)

			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting word review item: %v", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing study sessions and reviews: %v", err)
	}

	// Get list of seed files
	files, err := filepath.Glob("db/seeds/*.json")
	if err != nil {
		return fmt.Errorf("error finding seed files: %v", err)
	}

	for _, file := range files {
		fmt.Printf("Processing %s...\n", filepath.Base(file))

		// Read and parse seed file
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading seed file %s: %v", filepath.Base(file), err)
		}

		var seedFile SeedFile
		if err := json.Unmarshal(content, &seedFile); err != nil {
			return fmt.Errorf("error parsing seed file %s: %v", filepath.Base(file), err)
		}

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %v", err)
		}

		// Create or get group
		var groupID int64
		err = tx.QueryRow(
			"INSERT OR IGNORE INTO groups (name) VALUES (?) RETURNING id",
			seedFile.GroupName,
		).Scan(&groupID)
		
		if err == sql.ErrNoRows {
			err = tx.QueryRow(
				"SELECT id FROM groups WHERE name = ?",
				seedFile.GroupName,
			).Scan(&groupID)
		}
		
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating/getting group: %v", err)
		}

		// Insert words and create word-group associations
		for _, word := range seedFile.Words {
			partsJSON, err := json.Marshal(word.Parts)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error marshaling parts: %v", err)
			}

			var wordID int64
			err = tx.QueryRow(`
				INSERT INTO words (japanese, romaji, english, parts)
				VALUES (?, ?, ?, ?)
				RETURNING id
			`, word.Japanese, word.Romaji, word.English, partsJSON).Scan(&wordID)
			
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting word: %v", err)
			}

			_, err = tx.Exec(`
				INSERT OR IGNORE INTO word_groups (word_id, group_id)
				VALUES (?, ?)
			`, wordID, groupID)
			
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error creating word-group association: %v", err)
			}
		}

		// Create default study activity if it doesn't exist
		_, err = tx.Exec(`
			INSERT OR IGNORE INTO study_activities (name, description)
			VALUES ('Vocabulary Quiz', 'Practice your vocabulary with flashcards')
		`)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating default study activity: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing transaction: %v", err)
		}

		fmt.Printf("Successfully processed %s\n", filepath.Base(file))
	}

	return nil
} 