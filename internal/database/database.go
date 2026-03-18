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

// New создает новое соединение с базой данных SQLite
// Инициализирует базу данных, создает необходимые таблицы и возвращает объект базы данных
// Принимает путь к файлу базы данных и возвращает указатель на Db и возможную ошибку
func New(dbPath string) (*Db, error) {

	// Создаем директорию для базы данных, если она не существует
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Открываем соединение с базой данных SQLite
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Проверяем соединение с базой данных
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Создаем необходимые таблицы в базе данных
	if err = createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")

	return &Db{
		db: db,
	}, nil
}

// createTables создает необходимые таблицы в базе данных
// В текущей реализации создает таблицу typing_sessions для хранения сессий пользователей
// Возвращает ошибку в случае неудачи при создании таблицы
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

// SaveTypingLesson сохраняет результат сессии ввода текста в базу данных
// Принимает структуру LessonResults с результатами сессии
// Возвращает ошибку в случае неудачи при сохранении
func (d *Db) SaveTypingLeрsson(r models.LessonResults) error {
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

// CloseDB закрывает соединение с базой данных
// Безопасно закрывает соединение и выводит сообщение о состоянии операции
func (d *Db) CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v\n", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

// GetTypingLessons возвращает историю сессий для указанного пользователя
// Пока не реализована, возвращает пустой результат
// Параметр user - имя пользователя, для которого запрашивается история
// Возвращает срез LessonResults и возможную ошибку
func (d *Db) GetTypingLessons(user string) ([]models.LessonResults, error) {
	return nil, nil
}
