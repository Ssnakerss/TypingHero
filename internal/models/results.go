package models

import "time"

// User представляет пользователя приложения
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login"`
}

// LessonResults представляет результаты одной сессии ввода текста
// Содержит все метрики производительности пользователя
// Используется для хранения данных в базе данных и передачи между компонентами приложения
type LessonResults struct {
	Id             int           // Уникальный идентификатор сессии
	UserID         int           // Идентификатор пользователя в базе данных
	When           time.Time     // Время начала сессии
	Difficulty     int           // Уровень сложности текста (1-10)
	Wpm, ErrorRate float64       // Скорость ввода (слов в минуту) и процент ошибок
	TimeTaken      time.Duration // Общее время, затраченное на сессию
}

// Storage интерфейс для взаимодействия с хранилищем данных
// Определяет методы для сохранения и получения результатов сессий
// Позволяет использовать различные реализации хранилища (база данных, файлы и т.д.)
type Storage interface {
	// SaveTypingLesson сохраняет результаты сессии ввода текста
	// Принимает структуру LessonResults и возвращает ошибку в случае неудачи
	SaveTypingLesson(result LessonResults) error

	// GetTypingLessons возвращает историю сессий для указанного пользователя
	// Пока не реализована, будет использоваться для анализа прогресса пользователя
	// Принимает имя пользователя и возвращает срез результатов и возможную ошибку
	GetTypingLessons(user string) ([]LessonResults, error)

	// CreateUser создает нового пользователя или возвращает существующего
	// Принимает имя и никнейм пользователя, возвращает ID пользователя
	CreateUser(name, nickname string) (int, error)

	// GetUsers возвращает список всех пользователей
	GetUsers() ([]User, error)

	// GetUserStats возвращает статистику пользователя
	// Принимает ID пользователя и возвращает структуру с статистикой
	GetUserStats(userID int) (*UserStats, error)
}

// UserStats структура для хранения статистики пользователя
type UserStats struct {
	TotalSessions int     `json:"totalSessions"`
	AvgWpm        float64 `json:"avgWpm"`
	BestAccuracy  float64 `json:"bestAccuracy"`
}
