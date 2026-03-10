package console

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Ssnakerss/TypingHero/internal/models"
)

// Statistics for tracking typing performance across sessions
var stats = struct {
	maxWPM     float64 // Maximum typing speed achieved
	minWPM     float64 // Minimum typing speed achieved (excluding zero)
	totalWPM   float64 // Sum of all WPM scores for calculating average
	attempts   int     // Number of typing attempts
	bestText   string  // Text associated with best performance
	bestErrors int     // Errors in best performance
}{
	maxWPM:   -1, // Initialize to -1 so first positive value will be set as max
	minWPM:   -1, // Initialize to -1 to track first positive value
	totalWPM: 0,
	attempts: 0,
}

// Color codes for terminal output
const (
	ColorCyan   = "\033[36m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

// Text pools for each difficulty level (1-10)

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
	fmt.Println("  1. Select difficulty level (1-10)")
	fmt.Println("  2. Type the displayed text as fast and accurate as you can")
	fmt.Println("  3. View your typing speed (WPM) and error rate")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
}

var prevDiff = 1

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
		if input == "" {
			return prevDiff
		}
		_, err = fmt.Sscanf(input, "%d", &difficulty)
		if err != nil || difficulty < 1 || difficulty > 10 {
			fmt.Println(ColorRed + "Invalid input. Please enter a number between 1 and 10." + ColorReset)
			continue
		}
		prevDiff = difficulty
		return difficulty
	}
}

func getText(difficulty int) string {
	texts := models.TextPools[difficulty]
	// Use time to get pseudo-random selection
	index := rand.Intn(len(texts))
	//index := time.Now().UnixNano() % int64(len(texts))
	return texts[index]
}

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

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func calculateWPM(charsTyped int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	minutes := duration.Seconds() / 60.0
	words := float64(charsTyped) / 5.0 // Standard: 5 characters = 1 word
	return words / minutes
}

func calculateErrorRate(typed, target string) float64 {
	if len(target) == 0 {
		return 0
	}

	errors := 0
	targetRunes := []rune(target)
	typedRunes := []rune(typed)

	// Count mismatched characters
	for i := 0; i < len(typedRunes) && i < len(targetRunes); i++ {
		if typedRunes[i] != targetRunes[i] {
			errors++
		}
	}

	// Add remaining characters as errors if typed is longer
	if len(typedRunes) > len(targetRunes) {
		errors += len(typedRunes) - len(targetRunes)
	}

	return float64(errors) / float64(len(target)) * 100
}

func displayResults(ls models.LessonResults,
	target string,
	typed string,
) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(ColorCyan + ColorBold + "                    RESULTS" + ColorReset)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Format WPM with color based on performance
	wpmColor := ColorReset
	if ls.Wpm >= 60 {
		wpmColor = ColorGreen
	} else if ls.Wpm >= 40 {
		wpmColor = ColorYellow
	} else {
		wpmColor = ColorRed
	}

	fmt.Printf("  %sTyping Speed:%s  %s%.2f WPM%s\n", ColorCyan, ColorReset, wpmColor, ls.Wpm, ColorReset)
	fmt.Printf("  %sError Rate:%s    %.2f%%\n", ColorCyan, ColorReset, ls.ErrorRate)
	fmt.Printf("  %sTime Taken:%s    %.2f seconds\n", ColorCyan, ColorReset, ls.TimeTaken.Seconds())
	fmt.Println()

	// Show typed text with visual feedback
	fmt.Println(ColorCyan + "Your typing:" + ColorReset)
	showTypedText(target, typed)
	fmt.Println()

	// Performance message
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

func showTypedText(target, typed string) {
	targetRunes := []rune(target)
	typedRunes := []rune(typed)

	fmt.Print("  ")
	for i := 0; i < len(targetRunes); i++ {
		if i >= len(typedRunes) {
			// Not typed yet
			fmt.Printf("%c", targetRunes[i])
		} else if typedRunes[i] == targetRunes[i] {
			// Correct
			fmt.Printf(ColorGreen+"%c"+ColorReset, targetRunes[i])
		} else {
			// Incorrect
			fmt.Printf(ColorRed+"%c"+ColorReset, targetRunes[i])
		}
	}
	fmt.Println()
}

func playAgain() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(ColorCyan + "\nPlay again? (y/n): " + ColorReset)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes", "":
			return true
		case "n", "no":
			return false
		}
		fmt.Println(ColorRed + "Please enter 'y' or 'n'." + ColorReset)
	}
}

func printStats() {
	if stats.attempts == 0 {
		fmt.Println(ColorCyan + "No attempts recorded yet." + ColorReset)
		return
	}

	// Calculate average WPM
	var avgWPM float64
	if stats.attempts > 0 {
		avgWPM = stats.totalWPM / float64(stats.attempts)
	}

	// Display statistics with formatting
	fmt.Println()
	fmt.Println(ColorCyan + "Your typing statistics:" + ColorReset)
	fmt.Println(strings.Repeat("-", 40))

	// Format average with color based on performance
	avgColor := ColorReset
	if avgWPM >= 60 {
		avgColor = ColorGreen
	} else if avgWPM >= 40 {
		avgColor = ColorYellow
	} else {
		avgColor = ColorRed
	}

	fmt.Printf("  %sAttempts:%s       %d\n", ColorCyan, ColorReset, stats.attempts)
	fmt.Printf("  %sBest Speed:%s     %.2f WPM\n", ColorCyan, ColorReset, stats.maxWPM)
	if stats.minWPM > 0 {
		fmt.Printf("  %sWorst Speed:%s    %.2f WPM\n", ColorCyan, ColorReset, stats.minWPM)
	}
	fmt.Printf("  %sAverage Speed:%s  %s%.2f WPM%s\n", ColorCyan, ColorReset, avgColor, avgWPM, ColorReset)
	fmt.Println(strings.Repeat("-", 40))
}

func updateStats(wpm float64, text string, errors int) {
	// Initialize minWPM on first attempt
	if stats.attempts == 0 {
		stats.minWPM = wpm
	}

	// Update max WPM
	if wpm > stats.maxWPM {
		stats.maxWPM = wpm
		stats.bestText = text
		stats.bestErrors = errors
	}

	// Update min WPM (only for positive values)
	if wpm > 0 && (stats.attempts == 0 || wpm < stats.minWPM) {
		stats.minWPM = wpm
	}

	// Update total and count for average
	stats.totalWPM += wpm
	stats.attempts++
}

func RunGame(storage models.Storage) {
	printWelcome()

	ls := models.LessonResults{}
	for {
		// Print statistics at the beginning of each attempt
		printStats()

		ls.Difficulty = getDifficulty()
		text := getText(ls.Difficulty)
		displayText(text)

		// Start timer
		startTime := time.Now()

		// Get user input
		typed := getUserInput()

		// Stop timer
		ls.TimeTaken = time.Since(startTime)

		// Calculate statistics
		ls.Wpm = calculateWPM(len(typed), ls.TimeTaken)
		ls.ErrorRate = calculateErrorRate(typed, text)
		errors := int(float64(len(text)) * ls.ErrorRate / 100)

		// Update global statistics
		updateStats(ls.Wpm, text, errors)

		// Save session to database
		storage.SaveTypingLesson(ls)

		// Display results
		displayResults(ls, text, typed)

		// Ask to play again
		if !playAgain() {
			// Final statistics summary
			fmt.Println()
			fmt.Println(ColorCyan + "Final typing statistics:" + ColorReset)
			fmt.Println(strings.Repeat("=", 50))
			printStats()

			// Best performance details
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
			break
		}
	}
	printStats()
}
