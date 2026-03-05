package console

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
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
var textPools = map[int][]string{
	1: {
		"The cat sat on the mat.",
		"Dog runs in the park.",
		"Birds fly in the sky.",
		"The sun is bright today.",
		"I like to read books.",
	},
	2: {
		"The quick brown fox jumps over the lazy dog.",
		"She sells seashells by the seashore.",
		"A journey of a thousand miles begins with a single step.",
		"The early bird catches the worm.",
		"Time and tide wait for no man.",
	},
	3: {
		"Programming is the art of telling a computer what to do.",
		"The best way to predict the future is to create it.",
		"Success is not final, failure is not fatal.",
		"In the middle of difficulty lies opportunity.",
		"Knowledge is power but enthusiasm pulls the switch.",
	},
	4: {
		"Learning to code is like learning a new language. At first it seems impossible, then it becomes challenging, and finally it becomes natural.",
		"The only way to do great work is to love what you do. If you haven't found it yet, keep looking. Don't settle.",
		"Technology is best when it brings people together. It enables us to connect, share, and learn from each other.",
		"Every expert was once a beginner. Practice makes progress, and consistency is the key to mastery.",
	},
	5: {
		"Software development is not just about writing code; it's about solving problems and creating elegant solutions that make complex systems manageable.",
		"The debugger is twice as hard as writing the code in the first place. Therefore, if you write the code as cleverly as possible, you are by definition not smart enough to debug it.",
		"Any fool can write code that a computer can understand. Good programmers write code that humans can understand. Simplicity is the soul of efficiency.",
	},
	6: {
		"In the world of software, the best code is no code at all. Every new line of code you willingly bring into the world is code that has to be debugged, code that has to be read and understood, and code that has to be supported.",
		"The function of good software is to make the complex appear to be simple to the user. This requires deep understanding of both technology and human psychology.",
		"Programming isn't about what you know; it's about what you can figure out. The only way to go fast, is to go well. Quality is not an act, it is a habit.",
	},
	7: {
		"Design patterns are reusable solutions to commonly occurring problems in software design. They represent best practices evolved over time and provide a standard terminology that makes communication between developers more efficient.",
		"Clean code is happy code. It is writing code that is easy to understand, easy to modify, and easy to extend. The cost of cleaning up code is always less than the cost of maintaining messy code.",
		"Unit testing is not about finding bugs, it is about regression testing. It ensures that changes you make today don't break functionality that worked yesterday.",
	},
	8: {
		"Object-oriented programming was supposed to unify the perspectives of the programmer and the end user. However, modern OOP has become so complex that it often creates more problems than it solves.",
		"The premature optimization is the root of all evil. Yet we should not miss our opportunities to optimize critical sections of code that are executed millions of times.",
		"Dependency injection and inversion of control are powerful patterns that promote loose coupling and make systems more testable and maintainable over time.",
	},
	9: {
		"Functional programming concepts like immutability, higher-order functions, and pure functions can dramatically improve code quality by reducing side effects and making behavior more predictable.",
		"Microservices architecture enables teams to deploy independently, scale horizontally, and adopt different technologies. However, it introduces complexity in distributed systems management.",
		"Event-driven architectures allow systems to be more responsive and loosely coupled. By processing events asynchronously, applications can handle high loads while maintaining responsiveness.",
	},
	10: {
		"Concurrent programming in Go leverages goroutines and channels to create highly efficient and scalable systems. The select statement enables multiplexing between channel operations, while the sync package provides primitives for synchronization.",
		"Distributed systems must handle network partitions, partial failures, and eventual consistency. Understanding CAP theorem trade-offs and implementing proper retry mechanisms with exponential backoff is essential.",
		"Type-driven development and dependent types allow us to encode business rules at the type level, making illegal states unrepresentable and eliminating entire classes of runtime errors through compile-time verification.",
	},
}

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
	texts := textPools[difficulty]
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

func displayResults(wpm float64, errorRate float64, target string, typed string, duration time.Duration) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(ColorCyan + ColorBold + "                    RESULTS" + ColorReset)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Format WPM with color based on performance
	wpmColor := ColorReset
	if wpm >= 60 {
		wpmColor = ColorGreen
	} else if wpm >= 40 {
		wpmColor = ColorYellow
	} else {
		wpmColor = ColorRed
	}

	fmt.Printf("  %sTyping Speed:%s  %s%.2f WPM%s\n", ColorCyan, ColorReset, wpmColor, wpm, ColorReset)
	fmt.Printf("  %sError Rate:%s    %.2f%%\n", ColorCyan, ColorReset, errorRate)
	fmt.Printf("  %sTime Taken:%s    %.2f seconds\n", ColorCyan, ColorReset, duration.Seconds())
	fmt.Println()

	// Show typed text with visual feedback
	fmt.Println(ColorCyan + "Your typing:" + ColorReset)
	showTypedText(target, typed)
	fmt.Println()

	// Performance message
	if errorRate < 5 && wpm > 50 {
		fmt.Println(ColorGreen + ColorBold + "  ★ Excellent! Outstanding typing performance! ★" + ColorReset)
	} else if errorRate < 10 && wpm > 30 {
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

func RunGame() {
	printWelcome()

	for {
		// Print statistics at the beginning of each attempt
		printStats()

		difficulty := getDifficulty()
		text := getText(difficulty)
		displayText(text)

		// Start timer
		startTime := time.Now()

		// Get user input
		typed := getUserInput()

		// Stop timer
		duration := time.Since(startTime)

		// Calculate statistics
		wpm := calculateWPM(len(typed), duration)
		errorRate := calculateErrorRate(typed, text)
		errors := int(float64(len(text)) * errorRate / 100)

		// Update global statistics
		updateStats(wpm, text, errors)

		// Display results
		displayResults(wpm, errorRate, text, typed, duration)

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
