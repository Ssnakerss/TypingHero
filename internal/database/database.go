package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ssnakerss/TypingHero/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Db struct {
	db *sql.DB
}

func New(dbPath string) (*Db, error) {

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Open database connection
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create tables
	if err = createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")

	return &Db{
		db: db,
	}, nil
}

// createTables creates the necessary tables in the database
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS typing_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user TEXT NOT NULL,
		start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		difficulty INTEGER NOT NULL,
		wpm REAL NOT NULL,
		error_rate REAL NOT NULL,
		time_taken INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create typing_sessions table: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// SaveTypingLesson(result LessonResults) error
// 	GetTypingLessons(user string) ([]LessonResults, error)

func (d *Db) GetTypingLessons(user string) ([]models.LessonResults, error) {
	return nil, nil
}

// saveTypingSession saves a typing session to the database
func (d *Db) SaveTypingLesson(r models.LessonResults) error {
	query := `
	INSERT INTO typing_sessions (user, start, difficulty, wpm, errorRate, timeTaken)
	VALUES ( ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, r.User, r.When, r.Difficulty, r.Wpm, r.ErrorRate, r.TimeTaken)
	if err != nil {
		return fmt.Errorf("failed to save typing session: %v", err)
	}

	return nil
}

// closeDB closes the database connection
func (d *Db) CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v\n", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}
