package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Ssnakerss/TypingHero/internal/console"
	"github.com/Ssnakerss/TypingHero/internal/database"
	"github.com/Ssnakerss/TypingHero/internal/web"
)

func main() {
	// Initialize database
	dbPath := "data/typing_sessions.sqlite"

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Warning: Could not initialize database: %v", err)
	}
	defer db.CloseDB()

	// Define command line flags
	consoleMode := flag.Bool("c", false, "Run in console mode")
	webMode := flag.Bool("w", false, "Run in web mode")

	// Parse command line arguments
	flag.Parse()

	// Check if both modes are specified
	if *consoleMode && *webMode {
		log.Fatal("Error: Cannot run in both console and web modes simultaneously. Please specify only one mode.")
	}

	// Run in console mode if -c flag is present
	if *consoleMode {
		fmt.Println("Starting Typing Hero in console mode...")
		console.RunGame(db)
	}

	// Run in web mode if -w flag is present
	if *webMode {
		fmt.Println("Starting Typing Hero in web mode...")
		web.StartWeb()
		return
	}

	if !(*consoleMode || *webMode) {
		// Default behavior: show help if no mode is specified
		fmt.Println("No mode specified. Please use:")
		fmt.Println("  -c : Console mode")
		fmt.Println("  -w : Web mode")
		fmt.Println("Example: go run main.go -c  (for console mode)")
		fmt.Println("         go run main.go -w  (for web mode)")
	}
}
