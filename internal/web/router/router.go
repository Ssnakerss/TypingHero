package router

import (
	"github.com/Ssnakerss/TypingHero/internal/web/handlers"
	"github.com/go-chi/chi"
)

// New создает новый маршрутизатор HTTP-запросов
// Настраивает все маршруты для веб-приложения Typing Hero
// Принимает HandlerMaster для доступа к обработчикам и возвращает настроенный маршрутизатор
func New(hm *handlers.HandlerMaster) *chi.Mux {
	r := chi.NewRouter()

	// Маршрут для главной страницы
	// Обслуживает статические файлы из директории ./html
	r.Handle("/", hm.HomeHandler)

	// API маршруты для получения текста и расчета результатов
	r.Post("/api/text", hm.GetTextHandler)        // Получение текста для ввода
	r.Post("/api/result", hm.CalculateResultHandler) // Расчет результатов ввода

	return r
}