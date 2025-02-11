package handlers

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run migrations
	err = runTestMigrations(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Insert test data
	err = seedTestData(db)
	if err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	return db
}

func runTestMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			japanese TEXT NOT NULL,
			romaji TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT NOT NULL
		)`,
		`CREATE TABLE groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE word_groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
			UNIQUE(word_id, group_id)
		)`,
		`CREATE TABLE study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			thumbnail_url TEXT,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			study_activity_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE word_review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			study_session_id INTEGER NOT NULL,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
		)`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedTestData(db *sql.DB) error {
	testData := []string{
		`INSERT INTO words (japanese, romaji, english, parts) VALUES 
		('こんにちは', 'konnichiwa', 'hello', '{"type":"greeting"}')`,
		`INSERT INTO words (japanese, romaji, english, parts) VALUES 
		('さようなら', 'sayounara', 'goodbye', '{"type":"greeting"}')`,
		`INSERT INTO words (japanese, romaji, english, parts) VALUES 
		('ありがとう', 'arigatou', 'thank you', '{"type":"greeting"}')`,
		`INSERT INTO groups (name) VALUES ('Basic Greetings')`,
		`INSERT INTO word_groups (word_id, group_id) VALUES (1, 1)`,
		`INSERT INTO study_activities (name, description) 
		VALUES ('Vocabulary Quiz', 'Practice your vocabulary with flashcards')`,
		`INSERT INTO study_sessions (group_id, study_activity_id) 
		VALUES (1, 1)`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct) 
		VALUES (1, 1, true)`,
	}

	for _, data := range testData {
		_, err := db.Exec(data)
		if err != nil {
			return err
		}
	}

	return nil
}
