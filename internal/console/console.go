package console

import (
	"bufio"
	"context"
	"fmt"

	"os"
	"strings"
	"time"

	"github.com/Ssnakerss/TypingHero/internal/models"
	"github.com/Ssnakerss/TypingHero/internal/utils"
)

// Константы для цветового оформления вывода в терминале
// Используются escape-коды ANSI для изменения цвета и стиля текста
const (
	ColorCyan   = "\033[36m" // Голубой цвет для заголовков и инструкций
	ColorGreen  = "\033[32m" // Зеленый цвет для правильных символов и положительных сообщений
	ColorRed    = "\033[31m" // Красный цвет для ошибок и неправильных символов
	ColorYellow = "\033[33m" // Желтый цвет для предупреждений и средней производительности
	ColorReset  = "\033[0m"  // Сброс цвета и стиля к значениям по умолчанию
	ColorBold   = "\033[1m"  // Жирный шрифт для выделения важной информации
)

// Структура для хранения статистики пользователя за все сессии
// Позволяет отслеживать прогресс и лучшие результ��ты
var stats = struct {
	maxWPM     float64 // Максимальная скорость ввода, достигнутая пользователем
	minWPM     float64 // Минимальная скорость ввода (исключая нулевые значения)
	totalWPM   float64 // Сумма всех показателей WPM для расчета среднего значения
	attempts   int     // Общее количество попыток ввода текста
	bestText   string  // Текст, с которым была достигнута лучшая производительность
	bestErrors int     // Количество ошибок в лучшей попытке
}{
	maxWPM:   -1, // Инициализация значением -1, чтобы первое положительное значение стало максимумом
	minWPM:   -1, // Инициализация значением -1 для отслеживания первого положительного значения
	totalWPM: 0,
	attempts: 0,
}

// Функция вывода приветственного сообщения и инструкций
// Отображает анимационный логотип и основные правила игры
func printWelcome() {
	fmt.Println()
	fmt.Println(ColorCyan + ColorBold + `
    ____  ___      __  __  ___      ____  ____  ___  ____    ____  ___ __  __  ___  
   (  _ \( ___)   (  )(  )( __)    (  _ \(  _)/ __)(_  _)  (  __)/ __)(  )(  )( __) 
    ) _ < ) _)     )(__)(  ) _)      ) _ < )_( \__ \  )      ) _)( (__  ) __ ( (_ \ 
   (____/(____)   (______)(____)    (____/____)(___/ (__)   (____)\___)(__)(__)(___/ 
` + ColorReset)
	fmt.Println(ColorYellow + "         Welcome to Typing Hero - Improve your typing skills!" + ColorReset)
	fmt.Println()
	fmt.Println(ColorCyan + "Instructions:" + ColorReset)
	fmt.Println("  1. Select language (eng/rus)")
	fmt.Println("  2. Select difficulty level (1-10)")
	fmt.Println("  3. Type the displayed text as fast and accurate as you can")
	fmt.Println("  4. View your typing speed (WPM) and error rate")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
}

// Переменная для хранения предыдущего выбранного уровня сложности
// Позволяет пользователю пропустить ввод, используя предыдущее значение
var prevDiff = 1

// Переменная для хранения предыдущего выбранного языка
// Позволяет пользователю пропустить ввод, используя предыдущее значение
var prevLang = "eng"

// Функция получения языка от пользователя
// Обеспечивает ввод и валидацию значения (eng или rus)
func getLanguage() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(ColorCyan + "Select language (eng/rus): " + ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(ColorRed + "Error reading input. Please try again." + ColorReset)
			continue
		}

		input = strings.TrimSpace(input)
		// Если ввод пустой, возвращаем предыдущее значение языка
		if input == "" {
			return prevLang
		}
		// Проверяем корректность ввода: должно быть "eng" или "rus"
		if input != "eng" && input != "rus" {
			fmt.Println(ColorRed + "Invalid input. Please enter 'eng' or 'rus'." + ColorReset)
			continue
		}
		// Сохраняем выбранное значение как предыдущее
		prevLang = input
		return input
	}
}

// Функция получения уровня сложности от пользователя
// Обеспечивает ввод и валидацию значения в диапазоне от 1 до 10
func getDifficulty() int {
	reader := bufio.NewReader(os.Stdin)
	var difficulty int

	for {
		fmt.Print(ColorCyan + "Select difficulty (1-10): " + ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(ColorRed + "Error reading input. Please try again." + ColorReset)
			continue
		}

		input = strings.TrimSpace(input)
		// Если ввод пустой, возвращаем предыдущее значение сложности
		if input == "" {
			return prevDiff
		}
		// Пытаемся распарсить ввод как целое число
		_, err = fmt.Sscanf(input, "%d", &difficulty)
		// Проверяем корректность ввода: должно быть числом в диапазоне 1-10
		if err != nil || difficulty < 1 || difficulty > 10 {
			fmt.Println(ColorRed + "Invalid input. Please enter a number between 1 and 10." + ColorReset)
			continue
		}
		// Сохраняем выбранное значение как предыдущее
		prevDiff = difficulty
		return difficulty
	}
}

// Функция отображения текста для ввода
// Выводит инструкцию и сам текст, который нужно набрать
func displayText(text string) {
	fmt.Println()
	fmt.Println(ColorYellow + "Type the following text:" + ColorReset)
	fmt.Println()
	fmt.Println(ColorBold + text + ColorReset)
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println()
	fmt.Print(ColorCyan + "Start typing: " + ColorReset)
}

// Функция получения ввода пользователя
// Считывает строку из стандартного ввода и очищает от пробелов
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Функция отображения результатов сессии
// Выводит статистику ввода с цветовым кодированием в зависимости от производительности
func displayResults(ls models.LessonResults,
	target string,
	typed string,
) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(ColorCyan + ColorBold + "                    RESULTS" + ColorReset)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Определяем цвет для отображения скорости в зависимости от результата
	wpmColor := ColorReset
	if ls.Wpm >= 60 {
		wpmColor = ColorGreen // Отличная производительность
	} else if ls.Wpm >= 40 {
		wpmColor = ColorYellow // Хорошая производительность
	} else {
		wpmColor = ColorRed // Низкая производительность
	}

	// Выводим основные результаты
	fmt.Printf("  %sTyping Speed:%s  %s%.2f WPM%s\n", ColorCyan, ColorReset, wpmColor, ls.Wpm, ColorReset)
	fmt.Printf("  %sError Rate:%s    %.2f%%\n", ColorCyan, ColorReset, ls.ErrorRate)
	fmt.Printf("  %sTime Taken:%s    %.2f seconds\n", ColorCyan, ColorReset, ls.TimeTaken.Seconds())
	fmt.Println()

	// Показываем введенный текст с визуальной обратной связью
	fmt.Println(ColorCyan + "Your typing:" + ColorReset)
	showTypedText(target, typed)
	fmt.Println()

	// Выводим сообщение о производительности на основе результатов
	if ls.ErrorRate < 5 && ls.Wpm > 50 {
		fmt.Println(ColorGreen + ColorBold + "  ★ Excellent! Outstanding typing performance! ★" + ColorReset)
	} else if ls.ErrorRate < 10 && ls.Wpm > 30 {
		fmt.Println(ColorYellow + ColorBold + "  ★ Good job! Keep practicing! ★" + ColorReset)
	} else {
		fmt.Println(ColorRed + "  Keep practicing! Try again for better results." + ColorReset)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
}

// Функция визуального отображения введенного текста
// Показывает, какие символы были введены правильно, а какие нет
func showTypedText(target, typed string) {
	// Преобразуем строки в руны для корректной работы с Unicode
	targetRunes := []rune(target)
	typedRunes := []rune(typed)

	fmt.Print("  ")
	// Перебираем все символы эталонного текста
	for i := 0; i < len(targetRunes); i++ {
		if i >= len(typedRunes) {
			// Символ еще не введен - отображаем как есть
			fmt.Printf("%c", targetRunes[i])
		} else if typedRunes[i] == targetRunes[i] {
			// Символ введен правильно - отображаем зеленым
			fmt.Printf(ColorGreen+"%c"+ColorReset, targetRunes[i])
		} else {
			// Символ введен неправильно - отображаем красным
			fmt.Printf(ColorRed+"%c"+ColorReset, targetRunes[i])
		}
	}
	fmt.Println()
}

// Функция запроса на повторную игру
// Запрашивает у пользователя, хочет ли он сыграть еще раз
func playAgain() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(ColorCyan + "\nPlay again? (y/n): " + ColorReset)
		input, _ := reader.ReadString('\n')
		// Очищаем ввод и приводим к нижнему регистру
		input = strings.TrimSpace(strings.ToLower(input))

		// Обрабатываем различные варианты ответа
		switch input {
		case "y", "yes", "":
			return true // Пользователь хочет сыграть еще раз
		case "n", "no":
			return false // Пользователь хочет завершить игру
		}
		// Если ввод некорректен, запрашиваем повторно
		fmt.Println(ColorRed + "Please enter 'y' or 'n'." + ColorReset)
	}
}

// Функция вывода статистики пользователя
// Отображает текущую статистику за все сессии, включая среднюю, лучшую и худшую производительность
func printStats() {
	// Если нет ни одной попытки, сообщаем об этом
	if stats.attempts == 0 {
		fmt.Println(ColorCyan + "No attempts recorded yet." + ColorReset)
		return
	}

	// Рассчитываем среднюю скорость ввода
	var avgWPM float64
	if stats.attempts > 0 {
		avgWPM = stats.totalWPM / float64(stats.attempts)
	}

	// Выводим статистику с форматированием
	fmt.Println()
	fmt.Println(ColorCyan + "Your typing statistics:" + ColorReset)
	fmt.Println(strings.Repeat("-", 40))

	// Определяем цвет для средней скорости на основе результата
	avgColor := ColorReset
	if avgWPM >= 60 {
		avgColor = ColorGreen
	} else if avgWPM >= 40 {
		avgColor = ColorYellow
	} else {
		avgColor = ColorRed
	}

	// Выводим различные метрики статистики
	fmt.Printf("  %sAttempts:%s       %d\n", ColorCyan, ColorReset, stats.attempts)
	fmt.Printf("  %sBest Speed:%s     %.2f WPM\n", ColorCyan, ColorReset, stats.maxWPM)
	// Показываем худшую скорость только если она была установлена
	if stats.minWPM > 0 {
		fmt.Printf("  %sWorst Speed:%s    %.2f WPM\n", ColorCyan, ColorReset, stats.minWPM)
	}
	// Выводим среднюю скорость с соответствующим цветом
	fmt.Printf("  %sAverage Speed:%s  %s%.2f WPM%s\n", ColorCyan, ColorReset, avgColor, avgWPM, ColorReset)
	fmt.Println(strings.Repeat("-", 40))
}

// Функция обновления глобальной статистики
// Обновляет статистику на основе результатов текущей сессии
func updateStats(wpm float64, text string, errors int) {
	// Инициализируем минимальную скорость на первой попытке
	if stats.attempts == 0 {
		stats.minWPM = wpm
	}

	// Обновляем максимальную скорость, если текущая лучше
	if wpm > stats.maxWPM {
		stats.maxWPM = wpm
		stats.bestText = text     // Сохраняем текст лучшей попытки
		stats.bestErrors = errors // Сохраняем количество ошибок в лучшей попытке
	}

	// Обновляем м��нимальную скорость (только для положительных значений)
	if wpm > 0 && (stats.attempts == 0 || wpm < stats.minWPM) {
		stats.minWPM = wpm
	}

	// Обновляем сумму и счетчик для расчета среднего значения
	stats.totalWPM += wpm
	stats.attempts++
}

// Основная функция запуска игры
// Управляет основным циклом игры, обработкой сессий и взаимодействием с пользователем
func RunGame(ctx context.Context, c context.CancelFunc, storage models.Storage, userID int) {
	// Создаем структуру для хранения результатов сессии
	ls := models.LessonResults{} // Создаем структуру для хранения результатов сессии
	ls.UserID = userID           // Устанавливаем ID пользователя в структуру результатов
	printWelcome()               // Выводим приветственное сообщение
	defer c()                    // Отложенная отмена контекста при завершении функции
	for {
		select {
		case <-ctx.Done(): // Проверяем, не был ли отменен контекст
			return // Если контекст отменен, завершаем игру
		default:
			// Выводим статистику в начале каждой попытки
			printStats()

			// Получаем язык от пользователя
			lang := getLanguage()
			// Получаем уровень сложности от пользователя
			ls.Difficulty = getDifficulty()
			// Получаем случайный текст для выбранного уровня сложности и языка
			text := models.GetText(lang, ls.Difficulty)
			// Отображаем текст для ввода
			displayText(text)

			// Запоминаем время начала ввода
			startTime := time.Now()
			ls.When = startTime

			// Получаем ввод пользователя
			typed := getUserInput()

			// Фиксируем время, затраченное на ввод
			ls.TimeTaken = time.Since(startTime)

			// Рассчитываем статистику по вводу
			ls.Wpm = utils.CalculateWPM(len(typed), ls.TimeTaken)
			ls.ErrorRate = utils.CalculateErrorRate(typed, text)
			errors := int(float64(len(text)) * ls.ErrorRate / 100)

			// Обновляем глобальную статистику
			updateStats(ls.Wpm, text, errors)

			// Сохраняем результат сессии в базу данных
			storage.SaveTypingLesson(ls)

			// Выводим результаты сессии
			displayResults(ls, text, typed)

			// Спрашиваем, хочет ли пользователь сыграть еще раз
			if !playAgain() {
				// Выводим финальную статистику
				fmt.Println()
				fmt.Println(ColorCyan + "Final typing statistics:" + ColorReset)
				fmt.Println(strings.Repeat("=", 50))
				printStats()

				// Выводим детали лучшей производительности
				if stats.attempts > 0 {
					fmt.Println()
					fmt.Println(ColorCyan + "Your best performance:" + ColorReset)
					fmt.Println(strings.Repeat("-", 30))
					fmt.Printf("  %sSpeed:%s    %.2f WPM\n", ColorCyan, ColorReset, stats.maxWPM)
					fmt.Printf("  %sErrors:%s   %d\n", ColorCyan, ColorReset, stats.bestErrors)
					fmt.Printf("  %sText:%s     %s\n", ColorCyan, ColorReset, stats.bestText)
				}

				fmt.Println()
				fmt.Println(ColorCyan + "Thanks for playing Typing Hero! Goodbye!" + ColorReset)
				fmt.Println()
				return // Завершаем игру
			}
		}
	}
}
