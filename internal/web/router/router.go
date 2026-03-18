package router

import (
	"net/http"

	"github.com/Ssnakerss/TypingHero/internal/web/handlers"
)

// New создает новый HTTP-маршрутизатор для приложения
// Настраивает все необходимые маршруты для обработки запросов
// Принимает HandlerMaster для доступа к обработчикам и возвращает http.Handler
func New(hm *handlers.HandlerMaster) http.Handler {
	// Создаем новый маршрутизатор
	mux := http.NewServeMux()

	// Настраиваем маршруты для API
	mux.HandleFunc("/api/login", hm.LoginHandler)       // Аутентификация пользователя
	mux.HandleFunc("/api/users", hm.UsersHandler)       // Получение списка пользователей
	mux.HandleFunc("/api/text", hm.GetTextHandler)      // Получение текста для печати
	mux.HandleFunc("/api/result", hm.SaveResultHandler) // Сохранение результатов сессии

	// Настраиваем статические файлы
	// Главная страница
	mux.Handle("/", http.StripPrefix("/", http.HandlerFunc(hm.HomeHandler)))
	// Статические ресурсы (CSS, JS, изображения)
	mux.Handle("/css/", http.StripPrefix("/", http.FileServer(http.Dir("./cmd/html"))))
	mux.Handle("/js/", http.StripPrefix("/", http.FileServer(http.Dir("./cmd/html"))))
	mux.Handle("/images/", http.StripPrefix("/", http.FileServer(http.Dir("./cmd/html"))))

	// Возвращаем маршрутизатор как http.Handler
	return mux
}
