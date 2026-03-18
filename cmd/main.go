package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/TypingHero/internal/console"
	"github.com/Ssnakerss/TypingHero/internal/database"
	"github.com/Ssnakerss/TypingHero/internal/web"
)

// main - основная функция приложения Typing Hero
// Инициализирует базу данных, настраивает обработку сигналов и запускает консольную и веб-версии приложения
func main() {
	// Инициализация базы данных
	dbPath := "data/typing_sessions.sqlite"

	// Создаем новое соединение с базой данных
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Warning: Could not initialize database: %v", err)
	}
	// Отложенное закрытие соединения с базой данных
	defer db.CloseDB()

	// Создаем контекст с возможностью отмены для управления жизненным циклом приложения
	ctx, cancel := context.WithCancel(context.Background())
	// Отложенная отмена контекста при завершении функции
	defer cancel()

	// Настройка обработки системных сигналов для корректного завершения приложения
	go func() {
		// Канал для получения системных сигналов
		exit := make(chan os.Signal, 1)
		// Регистрируем интерес к сигналам завершения
		signal.Notify(exit,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		// Ожидаем получения сигнала
		sig := <-exit
		// Логируем полученный сигнал
		slog.Warn("signal received", "termination", sig)
		slog.Info("stopping server")
		// Отменяем контекст, что приведет к завершению всех горутин
		cancel()
	}()

	// Запускаем консольную версию игры в отдельной горутине
	go console.RunGame(ctx, cancel, db)
	// Запускаем веб-сервер в отдельной горутине
	go web.StartWeb(ctx, db)
	// Ожидаем отмены контекста (завершения приложения)
	<-ctx.Done()
}