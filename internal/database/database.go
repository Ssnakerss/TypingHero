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
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Проверяем соединение с базой данных
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db := &Db{
		db: conn,
	}

	// Создаем необходимые таблицы в базе данных
	if err = db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")

	return db, nil
}

// createTables создает необходимые таблицы в базе данных
// В текущей реализации создает таблицы users и typing_sessions
// Возвращает ошибку в случае неудачи при создании таблиц
func (d *Db) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		nickname TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS typing_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		difficulty INTEGER NOT NULL,
		wpm REAL NOT NULL,
		error_rate REAL NOT NULL,
		time_taken INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	if _, err := d.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create typing_sessions table: %v", err)
	}

	log.Println("Database tables created successfully")
	return nil
}

// SaveTypingLesson сохраняет результат сессии ввода текста в базу данных
// Принимает структуру LessonResults с результатами сессии
// Возвращает ошибку в случае неудачи при сохранении
func (d *Db) SaveTypingLesson(r models.LessonResults) error {
	query := `
	INSERT INTO typing_sessions (user_id, start, difficulty, wpm, error_rate, time_taken)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query, r.UserID, r.When, r.Difficulty, r.Wpm, r.ErrorRate, r.TimeTaken)
	if err != nil {
		return fmt.Errorf("failed to save typing session: %v", err)
	}

	return nil
}

// CloseDB закрывает соединение с базой данных
// Безопасно закрывает соединение и выводит сообщение о состоянии операции
func (d *Db) CloseDB() {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
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

// CreateUser создает нового пользователя в базе данных
// Проверяет, существует ли уже пользователь с таким именем или ником
// Если пользователь существует, обновляет дату последнего входа
// Если пользователь новый, создает запись с текущей датой регистрации и последнего входа
// Возвращает идентификатор пользователя и возможную ошибку
func (d *Db) CreateUser(name, nickname string) (int, error) {
	// Сначала проверяем, существует ли пользователь
	var userID int
	query := `SELECT id FROM users WHERE name = ? OR nickname = ?`
	err := d.db.QueryRow(query, name, nickname).Scan(&userID)

	if err == sql.ErrNoRows {
		// Пользователь не существует, создаем нового
		insertQuery := `INSERT INTO users (name, nickname, created_at, last_login) VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
		result, err := d.db.Exec(insertQuery, name, nickname)
		if err != nil {
			return 0, fmt.Errorf("failed to create user: %v", err)
		}
		id, _ := result.LastInsertId()
		return int(id), nil
	} else if err != nil {
		// Ошибка при выполнении запроса
		return 0, fmt.Errorf("failed to check user existence: %v", err)
	}

	// Пользователь существует, обновляем дату последнего входа
	updateQuery := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = d.db.Exec(updateQuery, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to update last login: %v", err)
	}

	return userID, nil
}

// GetUsers возвращает список всех зарегистрированных пользователей
// Сортирует пользователей по дате последнего входа (сначала последние)
// Возвращает срез пользователей и возможную ошибку
func (d *Db) GetUsers() ([]models.User, error) {
	query := `SELECT id, name, nickname, created_at, last_login FROM users ORDER BY last_login DESC`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Nickname, &user.CreatedAt, &user.LastLogin); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

// GetUserStats возвращает статистику пользователя
// Рассчитывает общее количество сессий, среднюю скорость ввода и лучшую точность
func (d *Db) GetUserStats(userID int) (*models.UserStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_sessions,
			AVG(wpm) as avg_wpm,
			MAX(error_rate) as best_accuracy
		FROM typing_sessions 
		WHERE user_id = ?
	`

	var totalSessions int
	var avgWpm float64
	var bestAccuracy float64

	err := d.db.QueryRow(query, userID).Scan(&totalSessions, &avgWpm, &bestAccuracy)
	if err != nil {
		return nil, fmt.Errorf("failed to query user stats: %v", err)
	}

	stats := &models.UserStats{
		TotalSessions: totalSessions,
		AvgWpm:        avgWpm,
		BestAccuracy:  100 - bestAccuracy, // Convert error rate to accuracy
	}

	return stats, nil
}
