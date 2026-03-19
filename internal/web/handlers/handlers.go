package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ssnakerss/TypingHero/internal/models"
	"github.com/Ssnakerss/TypingHero/internal/utils"
)

// HandlerMaster структура для управления всеми обработчиками
// Содержит ссылку на хранилище данных для доступа к базе данных
type HandlerMaster struct {
	storage models.Storage
}

// NewHandlerMaster создает новый экземпляр HandlerMaster
// Принимает интерфейс хранилища для работы с данными
func NewHandlerMaster(storage models.Storage) *HandlerMaster {
	return &HandlerMaster{storage: storage}
}

// HomeHandler serves the main HTML page.
func (hm *HandlerMaster) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/login.html")
}

// TypingHandler serves the typing game page.
func (hm *HandlerMaster) TypingHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/typing.html")
}

// LoginHandler обрабатывает запросы на вход пользователя
// Проверяет или создает пользователя в базе данных и возвращает его данные
func (hm *HandlerMaster) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
	}

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if user.Name == "" || user.Nickname == "" {
		json.NewEncoder(w).Encode(map[string]string{"message": "Name and nickname are required"})
		return
	}

	// Создаем или получаем существующего пользователя
	userID, err := hm.storage.CreateUser(user.Name, user.Nickname)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Возвращаем данные пользователя
	responseData := map[string]interface{}{
		"id":       userID,
		"name":     user.Name,
		"nickname": user.Nickname,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

// UsersHandler возвращает список всех зарегистрированных пользователей
// Используется для отображения списка пользователей на странице входа
func (hm *HandlerMaster) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем список пользователей из базы данных
	dbUsers, err := hm.storage.GetUsers()
	if err != nil {
		log.Printf("Error querying users: %v", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	var users []map[string]interface{}
	for _, u := range dbUsers {
		users = append(users, map[string]interface{}{
			"id":       u.ID,
			"name":     u.Name,
			"nickname": u.Nickname,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetTextHandler возвращает случайный текст для указанного уровня сложности
// Принимает сложность через POST-запрос в формате JSON и возвращает текст
func (hm *HandlerMaster) GetTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Структура для получения данных о сложности
	type difficultyRequest struct {
		Difficulty int `json:"difficulty"`
	}

	var req difficultyRequest
	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем диапазон сложности
	if req.Difficulty < 1 || req.Difficulty > 10 {
		http.Error(w, "Difficulty must be between 1 and 10", http.StatusBadRequest)
		return
	}

	// Получаем случайный текст для указанного уровня сложности
	text := models.GetText(req.Difficulty)

	// Формируем ответ
	response := map[string]string{"text": text}

	// Устанавливаем заголовок Content-Type и отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getCurrentUser извлекает ID пользователя из cookies или localStorage
// В реальном приложении нужно реализовать полноценную аутентификацию
// Здесь используется упрощенная версия для демонстрации
func getCurrentUser(r *http.Request) (int, error) {
	// Попробуем получить данные из cookie
	cookie, err := r.Cookie("currentUser")
	if err == nil {
		var userData map[string]interface{}
		if json.Unmarshal([]byte(cookie.Value), &userData) == nil {
			if id, ok := userData["id"].(float64); ok {
				return int(id), nil
			}
		}
	}

	// Альтернативно, можно проверить заголовки или тело запроса
	// Это зависит от реализации фронтенда

	return 0, fmt.Errorf("user not authenticated")
}

// SaveResultHandler обрабатывает сохранение результатов сессии
// Принимает данные о сессии и сохраняет их в базу данных
func (hm *HandlerMaster) SaveResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Структура для получения данных о результате
	type resultRequest struct {
		OriginalText string  `json:"originalText"`
		UserInput    string  `json:"userInput"`
		TimeTakenSec float64 `json:"timeTakenSec"`
		Difficulty   int     `json:"difficulty"`
	}

	var req resultRequest
	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if req.OriginalText == "" || req.UserInput == "" {
		http.Error(w, "Original text and user input are required", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из localStorage (передается через фронтенд)
	// В реальной реализации нужно добавить проверку аутентификации
	userID, err := getCurrentUser(r)
	if err != nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Рассчитываем метрики
	duration := time.Duration(req.TimeTakenSec * float64(time.Second))
	charsTyped := len(req.UserInput)
	wpm := utils.CalculateWPM(charsTyped, duration)
	errorRate := utils.CalculateErrorRate(req.UserInput, req.OriginalText)

	// Создаем структуру результата
	result := models.LessonResults{
		UserID:     userID,
		When:       time.Now(),
		Difficulty: req.Difficulty,
		Wpm:        wpm,
		ErrorRate:  errorRate,
		TimeTaken:  duration,
	}

	// Сохраняем результат в базу данных
	if err := hm.storage.SaveTypingLesson(result); err != nil {
		log.Printf("Error saving typing lesson: %v", err)
		http.Error(w, "Failed to save result", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с расчетными метриками
	response := map[string]interface{}{
		"wpm":       wpm,
		"errors":    int(errorRate * float64(len(req.OriginalText)) / 100),
		"accuracy":  100 - errorRate,
		"timeTaken": req.TimeTakenSec,
	}

	// Устанавливаем заголовок Content-Type и отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
