package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ssnakerss/TypingHero/internal/models"
	"github.com/Ssnakerss/TypingHero/internal/web/handlers"
	"github.com/Ssnakerss/TypingHero/internal/web/router"
)

// StartWeb запускает веб-сервер для приложения Typing Hero
// Создает HTTP-сервер, который обрабатывает запросы на порту 8080
// Запускается как горутина и работает до отмены контекста
// Принимает контекст для управления жизненным циклом и хранилище для доступа к данным
func StartWeb(ctx context.Context, db models.Storage) {
	hm := handlers.NewHandlerMaster(db)
	r := router.New(hm)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Typing Trainer server running at http://localhost:8080")
	go server.ListenAndServe()
	<-ctx.Done()
	server.Shutdown(ctx)

}