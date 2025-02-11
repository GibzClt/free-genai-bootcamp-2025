//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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