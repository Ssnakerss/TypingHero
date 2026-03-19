package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Ssnakerss/TypingHero/internal/console"
	"github.com/Ssnakerss/TypingHero/internal/constants"
	"github.com/Ssnakerss/TypingHero/internal/database"
	"github.com/Ssnakerss/TypingHero/internal/models"
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

	// Создаем канал для передачи выбранного пользователя
	userCh := make(chan models.User, 1)

	// Запускаем веб-сервер в отдельной горутине
	go web.StartWeb(ctx, db)

	// Запускаем выбор пользователя в отдельной горутине
	go func() {
		user, err := selectUser(ctx, db)
		if err != nil {
			slog.Error("Failed to select user", "error", err)
			cancel()
			return
		}
		userCh <- user
	}()

	// Получаем выбранного пользователя
	var user models.User
	select {
	case <-ctx.Done():
		return
	case user = <-userCh:
	}

	// Запускаем консольную версию игры в отдельной горутине
	go console.RunGame(ctx, cancel, db, user.ID)

	// Ожидаем отмены контекста (завершения приложения)
	<-ctx.Done()
}

// selectUser предоставляет интерфейс для выбора пользователя
// Позволяет выбрать существующего пользователя или зарегистрировать нового
// Возвращает структуру пользователя с установленным ID
func selectUser(ctx context.Context, db *database.Db) (models.User, error) {
	var user models.User

	fmt.Println(constants.ColorCyan + "\nWelcome to Typing Hero!" + constants.ColorReset)
	fmt.Println(constants.ColorCyan + "Please select a user or register a new one." + constants.ColorReset)
	fmt.Println(constants.ColorCyan + "Web interface available at http://localhost:8080" + constants.ColorReset)
	fmt.Println(strings.Repeat("-", 60))

	for {
		// Получаем список всех пользователей
		users, err := db.GetUsers()
		if err != nil {
			return user, fmt.Errorf("failed to query users: %v", err)
		}

		// Показываем список пользователей
		if len(users) > 0 {
			fmt.Println(constants.ColorYellow + "Registered users:" + constants.ColorReset)
			for i, u := range users {
				fmt.Printf("  %d. %s (@%s)\n", i+1, u.Name, u.Nickname)
			}
			fmt.Println()
		}

		// Показываем опции
		fmt.Println(constants.ColorCyan + "Options:" + constants.ColorReset)
		if len(users) > 0 {
			fmt.Println("  1-" + fmt.Sprintf("%d", len(users)) + ". Select existing user")
		}
		if len(users) > 0 {
			fmt.Println("  " + fmt.Sprintf("%d", len(users)+1) + ". Register new user")
		} else {
			fmt.Println("  1. Register new user")
		}
		fmt.Println("  q. Quit")

		// Получаем выбор пользователя
		fmt.Print(constants.ColorCyan + "\nSelect option: " + constants.ColorReset)
		var input string
		fmt.Fscanln(os.Stdin, &input)
		input = strings.TrimSpace(input)

		// Обрабатываем выход
		if input == "q" || input == "quit" || input == "exit" {
			os.Exit(0)
		}

		// Обрабатываем регистрацию нового пользователя
		if (len(users) == 0 && input == "1") ||
			(len(users) > 0 && input == fmt.Sprintf("%d", len(users)+1)) {
			fmt.Print(constants.ColorCyan + "Enter your full name: " + constants.ColorReset)
			fmt.Fscanln(os.Stdin, &user.Name)
			user.Name = strings.TrimSpace(user.Name)

			fmt.Print(constants.ColorCyan + "Enter your nickname: " + constants.ColorReset)
			fmt.Fscanln(os.Stdin, &user.Nickname)
			user.Nickname = strings.TrimSpace(user.Nickname)

			if user.Name == "" || user.Nickname == "" {
				fmt.Println(constants.ColorRed + "Name and nickname cannot be empty." + constants.ColorReset)
				continue
			}

			// Создаем или получаем существующего пользователя
			userID, err := db.CreateUser(user.Name, user.Nickname)
			if err != nil {
				fmt.Println(constants.ColorRed + "Failed to create user: " + err.Error() + constants.ColorReset)
				continue
			}
			user.ID = userID
			fmt.Printf(constants.ColorGreen+"Welcome back, %s!\n"+constants.ColorReset, user.Name)
			return user, nil
		}

		// Обрабатываем выбор существующего пользователя
		if len(users) > 0 {
			var index int
			_, err := fmt.Sscanf(input, "%d", &index)
			if err != nil || index < 1 || index > len(users) {
				fmt.Println(constants.ColorRed + "Invalid option. Please try again." + constants.ColorReset)
				continue
			}
			user = users[index-1]
			fmt.Printf(constants.ColorGreen+"Welcome back, %s!\n"+constants.ColorReset, user.Name)
			return user, nil
		}

		fmt.Println(constants.ColorRed + "Invalid option. Please try again." + constants.ColorReset)
	}
}
