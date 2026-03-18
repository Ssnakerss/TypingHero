package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"

	"github.com/Ssnakerss/TypingHero/internal/models"
)

// HandlerMaster управляет всеми обработчиками HTTP-запросов
// Содержит ссылки на все обработчики и хранилище данных
// Является центральной точкой маршрутизации запросов
type HandlerMaster struct {
	GetTextHandler         http.HandlerFunc // Обработчик для получения текста
	CalculateResultHandler http.HandlerFunc // Обработчик для расчета результата
	HomeHandler            http.Handler     // Обработчик для главной страницы

	Storage models.Storage // Интерфейс для доступа к данным
}

// NewHandlerMaster создает новый экземпляр HandlerMaster
// Инициализирует все обработчики и устанавливает ссылку на хранилище
// Возвращает указатель на новый HandlerMaster
func NewHandlerMaster(db models.Storage) *HandlerMaster {
	return &HandlerMaster{
		GetTextHandler:         getTextHandler,
		CalculateResultHandler: calculateResultHandler,
		HomeHandler:            http.FileServer(http.Dir("./html")),
		Storage:                db,
	}
}

// TextRequest представляет запрос на получение текста
// Содержит уровень сложности, для которого нужно получить текст
type TextRequest struct {
	Difficulty int `json:"difficulty"`
}

// ResultRequest представляет запрос на расчет результата
// Содержит оригинальный текст, введенный текст и время ввода
type ResultRequest struct {
	OriginalText string `json:"originalText"`
	UserInput    string `json:"userInput"`
	TimeTakenSec int    `json:"timeTakenSec"`
}

// TextResponse представляет ответ с текстом
// Используется для отправки текста клиенту
type TextResponse struct {
	Text string `json:"text"`
}

// ResultResponse представляет ответ с результатами ввода
// Содержит все метрики производительности пользователя
type ResultResponse struct {
	Errors    int     `json:"errors"`
	WPM       float64 `json:"wpm"`
	Accuracy  float64 `json:"accuracy"`
	CharCount int     `json:"charCount"`
	TimeTaken int     `json:"timeTaken"`
}

// generateText генерирует текст для заданного уровня сложности
// Для сложности 5 и выше комбинирует несколько предложений для увеличения длины
// Возвращает случайно выбранный текст из пула текстов
func generateText(difficulty int) string {
	// Ограничиваем уровень сложности в пределах 1-10
	if difficulty < 1 {
		difficulty = 1
	}
	if difficulty > 10 {
		difficulty = 10
	}

	// Получаем тексты для выбранного уровня сложности
	samples := models.TextPools[difficulty]
	// Выбираем случайный текст
	selected := samples[rand.Intn(len(samples))]

	// Для более сложных уровней комбинируем несколько предложений
	if difficulty >= 5 {
		count := difficulty / 2
		var result strings.Builder
		result.WriteString(selected)
		for i := 1; i < count; i++ {
			result.WriteString(" ")
			result.WriteString(samples[rand.Intn(len(samples))])
		}
		return result.String()
	}

	return selected
}

// abs возвращает абсолютное значение целого числа
// Используется для расчета количества ошибок при сравнении текстов
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// getTextHandler обрабатывает запросы на получение текста
// Принимает уровень сложности и возвращает случайный текст
// Поддерживает CORS для работы с фронтендом
func getTextHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка preflight-запросов CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование JSON-запроса
	var req TextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Генерация текста
	text := generateText(req.Difficulty)

	// Установка заголовков ответа
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	// Кодирование и отправка JSON-ответа
	json.NewEncoder(w).Encode(TextResponse{Text: text})
}

// calculateResultHandler обрабатывает запросы на расчет результатов
// Принимает оригинальный текст, введенный текст и время, и возвращает метрики производительности
// Выполняет сравнение текстов, расчет ошибок, точности и скорости ввода
func calculateResultHandler(w http.ResponseWriter, r *http.Request) {
	// Обработка preflight-запросов CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование JSON-запроса
	var req ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Расчет количества ошибок (по символам)
	errors := 0
	original := strings.TrimSpace(req.OriginalText)
	userInput := strings.TrimSpace(req.UserInput)

	// Определяем минимальную длину для сравнения
	minLen := len(original)
	if len(userInput) < minLen {
		minLen = len(userInput)
	}

	// Сравниваем символы по позициям
	for i := 0; i < minLen; i++ {
		if original[i] != userInput[i] {
			errors++
		}
	}

	// Добавляем разницу в длине как ошибки
	errors += abs(len(original) - len(userInput))

	// Расчет точности
	charCount := len(original)
	var accuracy float64 = 100
	if charCount > 0 {
		accuracy = float64(charCount-errors) / float64(charCount) * 100
		// Ограничиваем точность не ниже 0%
		if accuracy < 0 {
			accuracy = 0
		}
	}

	// Расчет скорости ввода (WPM - слова в минуту)
	// Стандарт: 5 символов = 1 слово
	words := float64(charCount) / 5.0
	var wpm float64 = 0
	// Рассчитываем WPM только если время больше 0
	if req.TimeTakenSec > 0 {
		wpm = words / float64(req.TimeTakenSec) * 60
	}

	// Установка заголовков ответа
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	// Кодирование и отправка JSON-ответа с результатами
	json.NewEncoder(w).Encode(ResultResponse{
		Errors:    errors,
		WPM:       wpm,
		Accuracy:  accuracy,
		CharCount: charCount,
		TimeTaken: req.TimeTakenSec,
	})
}